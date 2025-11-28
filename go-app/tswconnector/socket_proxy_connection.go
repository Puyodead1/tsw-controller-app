package tswconnector

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"
	"tsw_controller_app/chan_utils"
	"tsw_controller_app/logger"
	"tsw_controller_app/pubsub_utils"

	"github.com/gorilla/websocket"
)

const SOCKET_PROXY_CONNECTION_OUTGOING_QUEUE_BUFFER_SIZE = 32

var ErrCancelled = errors.New("cancelled")
var ErrUnknownMessage = errors.New("received unknown message type")
var ErrCloseMessage = errors.New("received close message")

type SocketProxyConnection_ConnectionResult struct {
	connection *websocket.Conn
	err        error
}

type SocketProxyConnection_MessageResult struct {
	message string
	err     error
}

type SocketProxyConnection struct {
	context         context.Context
	cancel          context.CancelFunc
	ServerAddr      string
	OutgoingChannel chan TSWConnector_Message
	Subscribers     *pubsub_utils.PubSubSlice[TSWConnector_Message]
}

var _ TSWConnector = (*SocketProxyConnection)(nil)

func (c *SocketProxyConnection) dial() chan SocketProxyConnection_ConnectionResult {
	done := make(chan SocketProxyConnection_ConnectionResult, 1)
	go func() {
		dialer := websocket.Dialer{
			HandshakeTimeout: 5 * time.Second,
		}
		u := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%d", c.ServerAddr, SOCKET_CONNECTION_PORT), Path: "/"}
		conn, _, err := dialer.Dial(u.String(), nil)
		if err != nil {
			logger.Logger.Error("[SocketProxyConnection::dial] could not connect to server", "error", err)
			done <- SocketProxyConnection_ConnectionResult{connection: nil, err: err}
		} else {
			done <- SocketProxyConnection_ConnectionResult{connection: conn, err: nil}
		}
	}()
	return done
}

func (c *SocketProxyConnection) waitForMessage(conn *websocket.Conn) chan SocketProxyConnection_MessageResult {
	received := make(chan SocketProxyConnection_MessageResult, 1)
	go func() {
		msg_type, msg, err := conn.ReadMessage()
		if err != nil {
			logger.Logger.Error("[SocketProxyConnection::waitForMessage] failed to read from connection", "error", err)
			received <- SocketProxyConnection_MessageResult{message: "", err: err}
		} else if msg_type == websocket.TextMessage {
			received <- SocketProxyConnection_MessageResult{message: string(msg), err: nil}
		} else if msg_type == websocket.CloseMessage {
			received <- SocketProxyConnection_MessageResult{message: "", err: ErrCloseMessage}
		} else {
			received <- SocketProxyConnection_MessageResult{message: "", err: ErrUnknownMessage}
		}
	}()
	return received
}

func (c *SocketProxyConnection) Subscribe() (chan TSWConnector_Message, func()) {
	return c.Subscribers.Subscribe()
}

func (c *SocketProxyConnection) Start() error {
	for {
		var connection *websocket.Conn
		select {
		case conn := <-c.dial():
			if conn.err != nil {
				<-time.After(5 * time.Second)
				continue
			}
			connection = conn.connection
		case <-c.context.Done():
			return nil
		}
		defer connection.Close()

		sender_ctx, cancel_sender := context.WithCancel(c.context)
		go func() {
			for {
				select {
				case <-sender_ctx.Done():
					return
				case message := <-c.OutgoingChannel:
					err := connection.WriteMessage(websocket.TextMessage, []byte(message.ToString()))
					if err != nil {
						cancel_sender()
						return
					}
				}
			}
		}()

	read_loop:
		for {
			select {
			case msg := <-c.waitForMessage(connection):
				if msg.err != nil {
					break read_loop
				}
				socket_message := TSWConnector_Message_FromString(msg.message)
				c.Subscribers.EmitTimeout(time.Second, socket_message)
			case <-c.context.Done():
				return nil
			}
		}
	}
}

func (c *SocketProxyConnection) Stop() error {
	c.cancel()
	return nil
}

func (c *SocketProxyConnection) Send(m TSWConnector_Message) error {
	return chan_utils.SendTimeout(c.OutgoingChannel, time.Second, m)
}

func NewSocketProxyConnection(ctx context.Context, addr string) *SocketProxyConnection {
	child_ctx, child_cancel := context.WithCancel(ctx)
	return &SocketProxyConnection{
		context:         child_ctx,
		cancel:          child_cancel,
		ServerAddr:      addr,
		OutgoingChannel: make(chan TSWConnector_Message, SOCKET_PROXY_CONNECTION_OUTGOING_QUEUE_BUFFER_SIZE),
		Subscribers:     pubsub_utils.NewPubSubSlice[TSWConnector_Message](),
	}
}
