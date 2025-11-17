package tickerutils

import (
	"context"
	"time"
)

type PausableTicker struct {
	ctx      context.Context
	dur      time.Duration
	ticker   *time.Ticker
	pauseC   chan time.Time
	pausedCs []chan time.Time
	C        chan time.Time
}

func (pt *PausableTicker) Paused() chan time.Time {
	c := make(chan time.Time)
	pt.pausedCs = append(pt.pausedCs, c)
	return c
}

func (pt *PausableTicker) Pause() {
	pause_time := time.Now()
	pt.pauseC <- pause_time
}

func (pt *PausableTicker) Start() {
	pt.ticker = time.NewTicker(pt.dur)
	go func() {
		select {
		case t := <-pt.ticker.C:
			pt.C <- t
		case pause_time := <-pt.pauseC:
			pt.ticker.Stop()
			for _, c := range pt.pausedCs {
				c <- pause_time
			}
			return
		case <-pt.ctx.Done():
			return
		}
	}()
}

func NewPausableTicker(ctx context.Context, d time.Duration) *PausableTicker {
	return &PausableTicker{
		ctx:      context.WithoutCancel(ctx),
		dur:      d,
		ticker:   nil,
		pauseC:   make(chan time.Time),
		pausedCs: []chan time.Time{},
		C:        make(chan time.Time),
	}
}
