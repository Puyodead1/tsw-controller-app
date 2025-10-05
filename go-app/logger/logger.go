package logger

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
)

type GlobalLogger struct {
	mutex     sync.RWMutex
	slogger   *slog.Logger
	listeners []chan string
}

func (g *GlobalLogger) Listen() (chan string, func()) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	channel := make(chan string)
	g.listeners = append(g.listeners, channel)
	unsubscribe := func() {
		for index, c := range g.listeners {
			if c == channel {
				g.listeners = append(g.listeners[:index], g.listeners[index+1:]...)
				break
			}
		}
		close(channel)
	}
	return channel, unsubscribe
}

func (g *GlobalLogger) PropertiesFromArgs(args ...any) map[string]string {
	properties := map[string]string{}
	for index, arg := range args {
		if index%2 == 1 {
			/* uneven indexes are the values */
			properties[fmt.Sprintf("%v", args[index-1])] = fmt.Sprintf("%#v", arg)
		}
	}
	return properties
}

func (g *GlobalLogger) Debug(msg string, args ...any) {
	g.slogger.Debug(msg, args...)
}

func (g *GlobalLogger) Info(msg string, args ...any) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.slogger.Info(msg, args...)

	if len(g.listeners) > 0 {
		properties := g.PropertiesFromArgs(args...)
		for _, c := range g.listeners {
			c <- fmt.Sprintf("%s | %v", msg, properties)
		}
	}
}

func (g *GlobalLogger) Error(msg string, args ...any) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.slogger.Error(msg, args...)

	if len(g.listeners) > 0 {
		properties := g.PropertiesFromArgs(args...)
		for _, c := range g.listeners {
			c <- fmt.Sprintf("%s | %v", msg, properties)
		}
	}
}

var Logger = GlobalLogger{
	slogger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	})),
	listeners: []chan string{},
}
