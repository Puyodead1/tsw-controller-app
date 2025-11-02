package pubsub_utils

import (
	"sync"
	"time"
	"tsw_controller_app/chan_utils"
)

type PubSubSlice[T any] struct {
	lock        sync.RWMutex
	subscribers []chan T
}

func (pbs *PubSubSlice[T]) Subscribe() (chan T, func()) {
	pbs.lock.Lock()
	defer pbs.lock.Unlock()
	sub_channel := make(chan T)
	pbs.subscribers = append(pbs.subscribers, sub_channel)
	return sub_channel, func() {
		pbs.lock.Lock()
		defer pbs.lock.Unlock()
		defer close(sub_channel)
		for index, sub := range pbs.subscribers {
			if sub == sub_channel {
				pbs.subscribers = append(pbs.subscribers[:index], pbs.subscribers[index+1:]...)
				break
			}
		}
	}
}

func (pbs *PubSubSlice[T]) EmitTimeout(duration time.Duration, e T) {
	pbs.lock.RLock()
	defer pbs.lock.RUnlock()
	for _, c := range pbs.subscribers {
		chan_utils.SendTimeout(c, duration, e)
	}
}

func NewPubSubSlice[T any]() *PubSubSlice[T] {
	return &PubSubSlice[T]{
		subscribers: []chan T{},
	}
}
