package main

import (
	"context"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type Connection struct {
	ctx         context.Context
	cancel      context.CancelFunc
	socket      *websocket.Conn
	outgoing    chan string
	subscribers []chan string
}

func (c *Connection) TryConnect() (*websocket.Conn, error) {
	url := url.URL{Scheme: "ws", Host: "127.0.0.1:63241"}
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		return nil, err
		// log.("[Connection::NewConnection] could not connect to server %e", err)
	}
	return conn, nil
}

func (c *Connection) Listen() {
	for {
		sock, err := c.TryConnect()
		if err != nil {
			log.Fatalf("[Connection::Listen] failed to connect to websocket %e", err)
			time.Sleep(5 * time.Second)
			continue
		}

		c.socket = sock
		sock_closed_ctx, done := context.WithCancel(context.Background())

		/* reader */
		go func() {
			defer done()
			for {
				msg_type, msg, err := sock.ReadMessage()
				if err != nil {
					log.Fatalf("[Connection::Listen] socket read error %e", err)
					return
				}
				if msg_type == websocket.CloseMessage {
					log.Printf("[Connection::Listen] closing current connection; received close message")
					return
				}
				if msg_type == websocket.TextMessage {
					for _, sub := range c.subscribers {
						sub <- string(msg)
					}
				} else {
					log.Printf("[Connection::Listen] received unsupported message type %d", msg_type)
				}
			}
		}()

		go func() {
			for msg := range c.outgoing {
				err := sock.WriteMessage(websocket.TextMessage, []byte(msg))
				if err != nil {
					log.Fatalf("[Connection::Listen] write error %e", err)
					return
				}
			}
		}()

		select {
		case <-c.ctx.Done():
			/* parent context closed -> shutdown */
			sock.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			time.Sleep(time.Second) // allow close frame to flush
		case <-sock_closed_ctx.Done():
			/* connection broke - restart connection loop */
			c.socket = nil
		}
	}
}

func NewConnection(ct context.Context) *Connection {
	ctx, cancel := context.WithCancel(ct)
	conn := &Connection{
		ctx:      ctx,
		cancel:   cancel,
		outgoing: make(chan string, 50),
	}
	return conn
}

func (c *Connection) Cancel() {
	c.cancel()
}

func (c *Connection) Send(msg string) {
	if c.socket != nil {
		c.outgoing <- msg
	}
}

func (c *Connection) Subscribe() (chan string, func()) {
	channel := make(chan string)
	c.subscribers = append(c.subscribers, channel)
	return channel, func() {
		for index, sub := range c.subscribers {
			if sub == channel {
				c.subscribers = append(c.subscribers[:index], c.subscribers[index+1:]...)
				break
			}
		}
	}
}
