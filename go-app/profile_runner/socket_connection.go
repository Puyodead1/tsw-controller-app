package profile_runner

import (
	"context"
	"net/http"
	"sort"
	"strings"
	"tsw_controller_app/logger"
	"tsw_controller_app/map_utils"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type SocketConnection_Message struct {
	EventName  string
	Properties map[string]string
}

type SocketConnection struct {
	WsUpgrader       *websocket.Upgrader
	Server           *http.Server
	OutgoingChannels *map_utils.LockMap[uuid.UUID, chan SocketConnection_Message]
	Subscribers      []chan SocketConnection_Message
}

func SocketConnectionMessage_FromString(msg string) SocketConnection_Message {
	parts := strings.Split(msg, ",")
	result := SocketConnection_Message{
		EventName:  "",
		Properties: make(map[string]string),
	}

	if len(parts) == 0 {
		return result
	}

	// first part is the event name
	result.EventName = parts[0]

	// the rest are key=value pairs
	for _, p := range parts[1:] {
		if kv := strings.SplitN(p, "=", 2); len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			val := strings.TrimSpace(kv[1])
			result.Properties[key] = val
		}
	}

	return result
}

func (msg SocketConnection_Message) ToString() string {
	var sb strings.Builder

	sb.WriteString(msg.EventName)

	keys := make([]string, 0, len(msg.Properties))
	for k := range msg.Properties {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		sb.WriteString(",")
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(msg.Properties[k])
	}

	return sb.String()
}

func (c *SocketConnection) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := c.WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Logger.Error("[SocketConnection::WebsocketHandler] websocket upgrade error", "error", err.Error())
		return
	}
	defer conn.Close()

	conn_id := uuid.New()
	outgoing_channel := make(chan SocketConnection_Message)
	c.OutgoingChannels.Set(conn_id, outgoing_channel)
	defer close(outgoing_channel)
	defer c.OutgoingChannels.Delete(conn_id)

	ctx_with_cancel, cancel_sender := context.WithCancel(r.Context())
	go func() {
		for {
			select {
			case <-ctx_with_cancel.Done():
				return
			case message := <-outgoing_channel:
				err := conn.WriteMessage(websocket.TextMessage, []byte(message.ToString()))
				if err != nil {
					cancel_sender()
					return
				}
			}
		}
	}()

	for {
		msg_type, msg, err := conn.ReadMessage()
		if err != nil {
			logger.Logger.Error("[ProfileRunner::WebsocketHandler] message read error", "error", err)
			return
		}

		if msg_type == websocket.CloseMessage {
			logger.Logger.Info("[ProfileRunner::WebsocketHandler] received close message from client")
			break
		}

		if msg_type == websocket.TextMessage {
			socket_message := SocketConnectionMessage_FromString(string(msg))
			logger.Logger.Info("[ProfileRunner::WebsocketHandler] received message from client", "message", socket_message)
			for _, sub := range c.Subscribers {
				sub <- socket_message
			}
		} else {
			logger.Logger.Info("[ProfileRunner::WebsocketHandler] received unsupported message %d", "message_type", msg_type)
		}
	}

	cancel_sender()
}

func (c *SocketConnection) Subscribe() (chan SocketConnection_Message, func()) {
	channel := make(chan SocketConnection_Message)
	c.Subscribers = append(c.Subscribers, channel)
	return channel, func() {
		for index, sub := range c.Subscribers {
			if sub == channel {
				c.Subscribers = append(c.Subscribers[:index], c.Subscribers[index+1:]...)
				break
			}
		}
	}
}

func (c *SocketConnection) Start() error {
	return c.Server.ListenAndServe()
}

func (c *SocketConnection) Send(m SocketConnection_Message) {
	c.OutgoingChannels.ForEach(func(channel chan SocketConnection_Message, key uuid.UUID) bool {
		channel <- m
		return true
	})
}

func NewSocketConnection() *SocketConnection {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":63241",
		Handler: mux,
	}
	controller := SocketConnection{
		WsUpgrader:       &websocket.Upgrader{},
		Server:           server,
		OutgoingChannels: map_utils.NewLockMap[uuid.UUID, chan SocketConnection_Message](),
	}
	mux.HandleFunc("/", controller.WebsocketHandler)
	return &controller
}
