package profile_runner

import (
	"context"
	"net/http"
	"tsw_controller_app/logger"

	"github.com/gorilla/websocket"
)

type SocketConnection struct {
	WsUpgrader      *websocket.Upgrader
	Server          *http.Server
	OutgoingChannel chan string
	Subscribers     []chan string
}

func (c *SocketConnection) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := c.WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Logger.Error("[SocketConnection::WebsocketHandler] websocket upgrade error", "error", err.Error())
		return
	}
	defer conn.Close()

	ctx_with_cancel, cancel_sender := context.WithCancel(r.Context())
	go func() {
		for {
			select {
			case <-ctx_with_cancel.Done():
				return
			case message := <-c.OutgoingChannel:
				err := conn.WriteMessage(websocket.TextMessage, []byte(message))
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
			for _, sub := range c.Subscribers {
				sub <- string(msg)
			}
		} else {
			logger.Logger.Info("[ProfileRunner::WebsocketHandler] received unsupported message %d", msg_type)
		}
	}

	cancel_sender()
}

func (c *SocketConnection) Subscribe() (chan string, func()) {
	channel := make(chan string)
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

func NewSocketConnection() *SocketConnection {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":63241",
		Handler: mux,
	}
	controller := SocketConnection{
		WsUpgrader:      &websocket.Upgrader{},
		Server:          server,
		OutgoingChannel: make(chan string),
	}
	mux.HandleFunc("/", controller.WebsocketHandler)
	return &controller
}
