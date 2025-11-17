package profile_runner

import (
	"context"
	"fmt"
	"tsw_controller_app/tswapi"
)

type ApiController_Command struct {
	Controls   string
	InputValue float64
}

type ApiController struct {
	API            *tswapi.TSWAPI
	ControlChannel chan ApiController_Command
}

func (c *ApiController_Command) ToString() string {
	return fmt.Sprintf("api_control_command:%s:%f", c.Controls, c.InputValue)
}

func (controller *ApiController) Run(ctx context.Context) func() {
	ctx_with_cancel, cancel := context.WithCancel(ctx)

	go func() {
		for {
			select {
			case <-ctx_with_cancel.Done():
				return
			case command := <-controller.ControlChannel:
				go controller.API.SetInputValue(command.Controls, command.InputValue)
			}
		}
	}()

	return cancel
}

func NewAPIController(twapi *tswapi.TSWAPI) *ApiController {
	controller := ApiController{
		API:            twapi,
		ControlChannel: make(chan ApiController_Command),
	}
	return &controller
}
