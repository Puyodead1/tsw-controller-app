package profile_runner

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"tsw_controller_app/logger"

	"github.com/gorilla/websocket"
)

type DirectController_Command struct {
	Controls   string
	InputValue float64
	Flags      []string
}

type DirectController struct {
	WsUpgrader     *websocket.Upgrader
	Server         *http.Server
	ControlChannel chan DirectController_Command
}

func (command *DirectController_Command) ToString() string {
	return fmt.Sprintf("%s,%f,%s", command.Controls, command.InputValue, strings.Join(command.Flags, "|"))
}

func (c *DirectController) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := c.WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Logger.Error("[DirectController::WebsocketHandler] websocket upgrade error", "error", err.Error())
		return
	}
	defer conn.Close()

	ctx_with_cancel, cancel_sender := context.WithCancel(r.Context())
	go func() {
		for {
			select {
			case <-ctx_with_cancel.Done():
				return
			case command := <-c.ControlChannel:
				command_str := command.ToString()
				err := conn.WriteMessage(websocket.TextMessage, []byte(command_str))
				if err != nil {
					cancel_sender()
					return
				}
			}
		}
	}()

	for {
		msg_type, _, err := conn.ReadMessage()
		if err != nil {
			logger.Logger.Error("[DirectController::WebsocketHandler] message read error", "error", err)
			return
		}

		if msg_type == websocket.CloseMessage {
			logger.Logger.Info("[DirectController::WebsocketHandler] received close message from client")
			break
		}
	}

	cancel_sender()
}

func (c *DirectController) Start() error {
	return c.Server.ListenAndServe()
}

func NewDirectController() *DirectController {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":63241",
		Handler: mux,
	}
	controller := DirectController{
		WsUpgrader: &websocket.Upgrader{},
		Server:     server,
	}
	mux.HandleFunc("/", controller.WebsocketHandler)
	return &controller
}
