package chan_utils

import (
	"errors"
	"time"
	"tsw_controller_app/logger"
)

var ErrTimeout = errors.New("timeout error")

func SendTimeout[T any](channel chan T, timeout_duration time.Duration, payload T) error {
	timeout := time.After(timeout_duration)
	select {
	case channel <- payload:
		// value sent to channel ok
		return nil
	case <-timeout:
		// timed out
		logger.Logger.Error("[SendTimeout] message timed out", "payload", payload, "timeout", timeout_duration)
		return ErrTimeout
	}
}
