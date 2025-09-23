package profile_runner

import (
	"net/http"
	"strconv"
	"strings"
	"tsw_controller_app/config"
	"tsw_controller_app/logger"

	"github.com/gorilla/websocket"
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
	WsUpgrader                  *websocket.Upgrader
	Server                      *http.Server
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

func (c *SyncController) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := c.WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Logger.LogF("[SyncController::WebsocketHandler] websocket upgrade error (%e)", err)
		return
	}
	defer conn.Close()

	for {
		msg_type, msg, err := conn.ReadMessage()
		if err != nil {
			logger.Logger.LogF("[SyncController::WebsocketHandler] read error (%e)", err)
			return
		}

		if msg_type == websocket.CloseMessage {
			logger.Logger.LogF("[SyncController::WebsocketHandler] received close message")
			break
		}

		if msg_type == websocket.TextMessage {
			/* message should follow format sync_control,{identifier},{value} */
			message_parts := strings.Split(string(msg), ",")
			/* skip message if not sync_control message */
			if message_parts[0] != "sync_control" || len(message_parts) != 3 {
				continue
			}

			control_state, has_control_state := c.ControlState[message_parts[1]]
			if !has_control_state {
				control_state = SyncController_ControlState{
					Identifier:   message_parts[1],
					CurrentValue: 0.0,
					TargetValue:  0.0,
					Moving:       0,
				}
			}
			current_value, err := strconv.ParseFloat(message_parts[2], 64)
			if err != nil {
				control_state.CurrentValue = current_value
				for _, channel := range c.ControlStateChangedChannels {
					channel <- control_state
				}
			}
			c.ControlState[message_parts[1]] = control_state
		}
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

func (c *SyncController) Start() error {
	return c.Server.ListenAndServe()
}

func NewSyncController() *SyncController {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":63242",
		Handler: mux,
	}
	controller := SyncController{
		WsUpgrader:   &websocket.Upgrader{},
		Server:       server,
		ControlState: map[string]SyncController_ControlState{},
	}
	mux.HandleFunc("/", controller.WebsocketHandler)
	return &controller
}
