package logger

import "fmt"

const LOGGER_BUFSIZE = 1000

type logger struct {
	Logs    []string
	Channel chan string
}

func (l *logger) LogF(format string, a ...any) {
	log := fmt.Sprintf(format, a...)
	fmt.Print(log)

	l.Logs = append(l.Logs, log)
	l.Channel <- log
}

var Logger = logger{
	Logs:    []string{},
	Channel: make(chan string, LOGGER_BUFSIZE),
}
