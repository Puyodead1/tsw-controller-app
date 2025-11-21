package cabdebugger

import (
	"context"
	"errors"
	"net"
	"strconv"
	"sync"
	"time"
	"tsw_controller_app/map_utils"
	"tsw_controller_app/tswapi"
	"tsw_controller_app/tswconnector"
)

type PropertyName = string

type CabDebugger_ControlState_Control struct {
	Identifier             string
	PropertyName           PropertyName
	CurrentValue           float64
	CurrentNormalizedValue float64
}

type CabDebugger_ControlState struct {
	DrivableActorName string
	Controls          *map_utils.LockMap[PropertyName, CabDebugger_ControlState_Control]
}

type CabDebugger_Config struct {
	TSWAPISubscriptionIDStart int
}

type CabDebugger struct {
	updateControlStateFromAPIMutex sync.Mutex
	SocketConnection               *tswconnector.SocketConnection
	TSWAPI                         *tswapi.TSWAPI
	Config                         CabDebugger_Config
	State                          CabDebugger_ControlState
}

var ErrAlreadyLocked = errors.New("already locked error")

func (cd *CabDebugger) updateControlStateFromAPI() error {
	if cd.TSWAPI.Enabled() {
		/* try to acquire lock ; if already locked we skip */
		did_lock := cd.updateControlStateFromAPIMutex.TryLock()
		if !did_lock {
			return ErrAlreadyLocked
		}
		defer cd.updateControlStateFromAPIMutex.Unlock()

		result, err := cd.TSWAPI.GetCurrentDrivableActorSubscription(cd.Config.TSWAPISubscriptionIDStart)
		if (err != nil &&
			/* don't do anything further if the comm api key is missing */
			!errors.Is(err, tswapi.ErrMissingCommAPIKey) &&
			/* don't do anything for an OpError */
			!errors.As(err, new(*net.OpError))) ||
			result.ObjectClass != cd.State.DrivableActorName {
			cd.TSWAPI.DeleteSubscription(cd.Config.TSWAPISubscriptionIDStart)
			if err := cd.TSWAPI.CreateCurrentDrivableActorSubscription(cd.Config.TSWAPISubscriptionIDStart); err != nil {
				return err
			}
			if result, err = cd.TSWAPI.GetCurrentDrivableActorSubscription(cd.Config.TSWAPISubscriptionIDStart); err != nil {
				return err
			}
		}

		cd.State.Controls.Clear()
		cd.State.DrivableActorName = result.ObjectClass
		for property_name, control := range result.Controls {
			control_state, _ := cd.State.Controls.Get(property_name)
			control_state.Identifier = control.Identifier
			control_state.PropertyName = control.PropertyName
			control_state.CurrentValue = control.CurrentValue
			control_state.CurrentNormalizedValue = control.CurrentNormalizedValue
			cd.State.Controls.Set(property_name, control_state)
		}
	}

	return nil
}

func (cd *CabDebugger) UpdateConfig(config CabDebugger_Config) {
	cd.Config = config
}

func (cd *CabDebugger) Clear() {
	cd.State.Controls.Clear()
}

func (cd *CabDebugger) Start(ctx context.Context) {
	childctx := context.WithoutCancel(ctx)
	go func() {
		socket_channel, unsubscribe_socket_channel := cd.SocketConnection.Subscribe()
		ticker := time.NewTicker(333 * time.Millisecond)
		for {
			select {
			case msg := <-socket_channel:
				if msg.EventName == "sync_control_value" {
					control_state, has_control_state := cd.State.Controls.Get(msg.Properties["property"])
					if cd.TSWAPI.Enabled() && !has_control_state {
						/* if the API is enabled - it should drive the existance of the controls */
						continue
					}

					control_state.Identifier = msg.Properties["name"]
					control_state.PropertyName = msg.Properties["property"]
					current_value, _ := strconv.ParseFloat(msg.Properties["value"], 64)
					current_normalized_value, _ := strconv.ParseFloat(msg.Properties["normalized_value"], 64)
					control_state.CurrentValue = current_value
					control_state.CurrentNormalizedValue = current_normalized_value
					cd.State.Controls.Set(msg.Properties["property"], control_state)
				}
			case <-ticker.C:
				go cd.updateControlStateFromAPI()
			case <-childctx.Done():
				ticker.Stop()
				unsubscribe_socket_channel()
				return
			}
		}
	}()

}

func NewCabDebugger(tswapi *tswapi.TSWAPI, socket_conn *tswconnector.SocketConnection, config CabDebugger_Config) *CabDebugger {
	return &CabDebugger{
		updateControlStateFromAPIMutex: sync.Mutex{},
		SocketConnection:               socket_conn,
		TSWAPI:                         tswapi,
		Config:                         config,
		State: CabDebugger_ControlState{
			DrivableActorName: "",
			Controls:          map_utils.NewLockMap[PropertyName, CabDebugger_ControlState_Control](),
		},
	}
}
