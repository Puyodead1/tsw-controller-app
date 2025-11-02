package func_utils

import (
	"sync"
	"time"
)

func Throttle(fn func(), interval time.Duration) func() {
	var mu sync.Mutex
	var last time.Time

	return func() {
		mu.Lock()
		defer mu.Unlock()

		now := time.Now()
		if now.Sub(last) >= interval {
			last = now
			fn()
		}
	}
}
