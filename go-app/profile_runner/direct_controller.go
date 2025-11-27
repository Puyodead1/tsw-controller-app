package profile_runner

import (
	"context"
	"fmt"
	"strings"
	"tsw_controller_app/tswconnector"
)

const DIRECT_CONTROLLER_QUEUE_BUFFER_SIZE = 32

type DirectController_Command struct {
	Controls   string
	InputValue float64
	Flags      []string
}

type DirectController struct {
	SocketConnection *tswconnector.SocketConnection
	ControlChannel   chan DirectController_Command
}

func (command *DirectController_Command) ToSocketMessage() tswconnector.TSWConnector_Message {
	return tswconnector.TSWConnector_Message{
		EventName: "direct_control",
		Properties: map[string]string{
			"controls": command.Controls,
			"value":    fmt.Sprintf("%f", command.InputValue),
			"flags":    strings.Join(command.Flags, "|"),
		},
	}
}

func (controller *DirectController) Run(ctx context.Context) func() {
	ctx_with_cancel, cancel := context.WithCancel(ctx)

	go func() {
		for {
			select {
			case <-ctx_with_cancel.Done():
				return
			case command := <-controller.ControlChannel:
				controller.SocketConnection.Send(command.ToSocketMessage())
			}
		}
	}()

	return cancel
}

func NewDirectController(connection *tswconnector.SocketConnection) *DirectController {
	controller := DirectController{
		SocketConnection: connection,
		ControlChannel:   make(chan DirectController_Command, DIRECT_CONTROLLER_QUEUE_BUFFER_SIZE),
	}
	return &controller
}
