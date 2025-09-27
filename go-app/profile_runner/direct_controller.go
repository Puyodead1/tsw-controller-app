package profile_runner

import (
	"context"
	"fmt"
	"strings"
)

type DirectController_Command struct {
	Controls   string
	InputValue float64
	Flags      []string
}

type DirectController struct {
	SocketConnection *SocketConnection
	ControlChannel   chan DirectController_Command
}

func (command *DirectController_Command) ToString() string {
	return fmt.Sprintf("%s,%f,%s", command.Controls, command.InputValue, strings.Join(command.Flags, "|"))
}

func (controller *DirectController) Run(ctx context.Context) func() {
	ctx_with_cancel, cancel := context.WithCancel(ctx)

	go func() {
		for {
			select {
			case <-ctx_with_cancel.Done():
				return
			case command := <-controller.ControlChannel:
				controller.SocketConnection.OutgoingChannel <- command.ToString()
			}
		}
	}()

	return cancel
}

func NewDirectController(connection *SocketConnection) *DirectController {
	controller := DirectController{
		SocketConnection: connection,
		ControlChannel:   make(chan DirectController_Command),
	}
	return &controller
}
