package profile_runner

import (
	"context"
	"strconv"
	"tsw_controller_app/config"
)

type SyncController_ControlState struct {
	Identifier   string
	CurrentValue float64
	TargetValue  float64
	/** [-1,0,1] -> decreasing, idle, increasing */
	Moving         int
	ControlProfile config.Config_Controller_Profile_Control_Assignment_SyncControl
}

type SyncController struct {
	SocketConnection            *SocketConnection
	ControlState                map[string]SyncController_ControlState
	ControlStateChangedChannels []chan SyncController_ControlState
}

func (c *SyncController) UpdateControlStateMoving(identifier string, moving int) {
	if state, has_state := c.ControlState[identifier]; has_state {
		state.Moving = moving
		c.ControlState[identifier] = state
		for _, channel := range c.ControlStateChangedChannels {
			channel <- state
		}
	}
}

func (c *SyncController) UpdateControlStateTargetValue(identifier string, targetValue float64, profile config.Config_Controller_Profile_Control_Assignment_SyncControl) {
	state, has_state := c.ControlState[identifier]
	if !has_state {
		state = SyncController_ControlState{
			Identifier:     identifier,
			CurrentValue:   targetValue,
			TargetValue:    targetValue,
			Moving:         0,
			ControlProfile: profile,
		}
	}
	state.TargetValue = targetValue
	c.ControlState[identifier] = state

	for _, channel := range c.ControlStateChangedChannels {
		channel <- state
	}
}

func (c *SyncController) Subscribe() (chan SyncController_ControlState, func()) {
	channel := make(chan SyncController_ControlState)
	c.ControlStateChangedChannels = append(c.ControlStateChangedChannels, channel)
	unsubscribe := func() {
		for index, ch := range c.ControlStateChangedChannels {
			if ch == channel {
				c.ControlStateChangedChannels = append(c.ControlStateChangedChannels[:index], c.ControlStateChangedChannels[index+1:]...)
				break
			}
		}
	}
	return channel, unsubscribe
}

func (c *SyncController) Run(ctx context.Context) func() {
	ctx_with_cancel, cancel := context.WithCancel(ctx)

	go func() {
		incoming_channel, unsubscribe := c.SocketConnection.Subscribe()
		defer unsubscribe()

		for {
			select {
			case <-ctx_with_cancel.Done():
				return
			case msg := <-incoming_channel:
				/* skip message if not sync_control message */
				if msg.EventName == "sync_control" {
					continue
				}

				control_state, has_control_state := c.ControlState[msg.Properties["name"]]
				if !has_control_state {
					control_state = SyncController_ControlState{
						Identifier:   msg.Properties["name"],
						CurrentValue: 0.0,
						TargetValue:  0.0,
						Moving:       0,
					}
				}
				current_value, err := strconv.ParseFloat(msg.Properties["value"], 64)
				if err != nil {
					control_state.CurrentValue = current_value
					for _, channel := range c.ControlStateChangedChannels {
						channel <- control_state
					}
				}
				c.ControlState[msg.Properties["name"]] = control_state
			}
		}
	}()

	return cancel
}

func NewSyncController(connection *SocketConnection) *SyncController {
	controller := SyncController{
		SocketConnection:            connection,
		ControlState:                map[string]SyncController_ControlState{},
		ControlStateChangedChannels: []chan SyncController_ControlState{},
	}
	return &controller
}
