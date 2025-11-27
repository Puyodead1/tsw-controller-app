package action_sequencer

import (
	"context"
	"strings"
	"time"
	"tsw_controller_app/chan_utils"
	"tsw_controller_app/logger"

	"github.com/go-vgo/robotgo"
)

const ACTIONS_QUEUE_BUFFER_SIZE = 32

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
		ActionsQueue: make(chan ActionSequencerAction, ACTIONS_QUEUE_BUFFER_SIZE),
	}
}

func (seq *ActionSequencer) Enqueue(action ActionSequencerAction) {
	chan_utils.SendTimeout(seq.ActionsQueue, time.Second, action)
}

func (seq *ActionSequencer) ToggleKeys(keys []string, modifiers []string, state string) {
	execution_groups := [][]string{}
	switch state {
	case "down":
		execution_groups = [][]string{modifiers, keys}
	case "up":
		execution_groups = [][]string{keys, modifiers}
	}
	for _, key_group := range execution_groups {
		for _, key := range key_group {
			robotgo.KeyToggle(key, state)
		}
		robotgo.MilliSleep(30)
	}
}

func (seq *ActionSequencer) Run(ctx context.Context) context.CancelFunc {
	modifier_keys_map := map[string]bool{
		"cmd":     true,
		"lcmd":    true,
		"rcmd":    true,
		"alt":     true,
		"lalt":    true,
		"ralt":    true,
		"ctrl":    true,
		"lctrl":   true,
		"rctrl":   true,
		"control": true,
		"shift":   true,
		"lshift":  true,
		"rshift":  true,
	}

	ctx_with_cancel, cancel := context.WithCancel(ctx)
	go func() {
		for {
			select {
			case <-ctx_with_cancel.Done():
				return
			case action := <-seq.ActionsQueue:
				logger.Logger.Debug("[ActionSequencer::Run] received action from queue", "action", action)
				keys_list := strings.Split(action.Keys, "+")
				modifier_keys := []string{}
				other_keys := []string{}
				for _, input := range keys_list {
					key := strings.ToLower(input)
					if is_modifier_key, has_is_modifier_key := modifier_keys_map[key]; has_is_modifier_key && is_modifier_key {
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
