package chan_utils

import (
	"time"
	"tsw_controller_app/logger"
)

func SendTimeout[T any](channel chan T, timeout_duration time.Duration, payload T) {
	timeout := time.After(timeout_duration)
	select {
	case channel <- payload:
		// value sent to channel ok
	case <-timeout:
		// timed out
		logger.Logger.Error("[SendTimeout] message timed out", "payload", payload, "timeout", timeout_duration)
	}
}
