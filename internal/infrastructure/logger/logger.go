package logger

import (
	"context"
	"fmt"
	"os"
)

type Logger struct {
	logs chan string
}

func NewLogger() *Logger {
	return &Logger{
		logs: make(chan string, 256),
	}
}

func (l *Logger) Log(msg string) {
	select {
	case l.logs <- msg:
	default:
	}
}

func (l *Logger) Logf(format string, args ...any) {
	l.Log(fmt.Sprintf(format, args...))
}

func (l *Logger) Run(ctx context.Context) {
	for {
		select {
		case msg := <-l.logs:
			_, _ = fmt.Fprintln(os.Stdout, msg)
		case <-ctx.Done():
			return
		}
	}
}
