package action_sequencer

import (
	"context"
	"strings"
	"tsw_controller_app/logger"

	"github.com/go-vgo/robotgo"
)

type ActionSequencerAction struct {
	Keys      string
	PressTime float64
	WaitTime  float64
	Release   bool
}

type ActionSequencer struct {
	ActionsQueue chan ActionSequencerAction
}

func New() *ActionSequencer {
	return &ActionSequencer{
		ActionsQueue: make(chan ActionSequencerAction),
	}
}

func (seq *ActionSequencer) Enqueue(action ActionSequencerAction) {
	seq.ActionsQueue <- action
}

func (seq *ActionSequencer) ToggleKeys(keys []string, modifiers []string, state string) {
	if state == "up" {
		for _, key := range keys {
			robotgo.KeyToggle(key, "up")
		}
		for _, key := range modifiers {
			robotgo.KeyToggle(key, "up")
		}
	} else if state == "down" {
		for _, key := range modifiers {
			robotgo.KeyToggle(key, "down")
		}
		for _, key := range keys {
			robotgo.KeyToggle(key, "down")
		}
	}
}

func (seq *ActionSequencer) Run(ctx context.Context) context.CancelFunc {
	ctx_with_cancel, cancel := context.WithCancel(ctx)
	go func() {
		for {
			select {
			case <-ctx_with_cancel.Done():
				return
			case action := <-seq.ActionsQueue:
				logger.Logger.Info("[ActionSequencer::Run] received action from queue", "action", action)
				keys_list := strings.Split(action.Keys, "+")
				modifier_keys := []string{}
				other_keys := []string{}
				for _, input := range keys_list {
					key := strings.ToLower(input)
					if key == "ctrl" || key == "control" || key == "alt" || key == "meta" || key == "cmd" || key == "command" {
						modifier_keys = append(modifier_keys, key)
					} else {
						other_keys = append(other_keys, key)
					}
				}

				if action.Release {
					seq.ToggleKeys(other_keys, modifier_keys, "up")
				} else {
					seq.ToggleKeys(other_keys, modifier_keys, "down")
					if action.PressTime != 0 {
						robotgo.MilliSleep(int(action.PressTime * 1000))
						seq.ToggleKeys(other_keys, modifier_keys, "up")
					}
					if action.WaitTime != 0 {
						robotgo.MilliSleep(int(action.WaitTime * 1000))
					}
				}
			}
		}
	}()
	return cancel
}
