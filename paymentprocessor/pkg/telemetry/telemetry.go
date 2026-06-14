package telemetry

import (
	"context"
	"log/slog"
	"sync"
)

type logBagKey struct{}

type LogBag struct {
	mu     sync.Mutex
	fields []any
}

func (l *LogBag) AddField(ctx context.Context, fields ...any) {
	l.mu.Lock()
	l.fields = append(l.fields, fields...)
	l.mu.Unlock()
}

func ExecuteWithLog(ctx context.Context, f func(ctx context.Context) error) error {
	bag := LogBag{
		mu:     sync.Mutex{},
		fields: make([]any, 0),
	}
	exCtx := context.WithValue(ctx, logBagKey{}, &bag)
	err := f(exCtx)

	msg := "process finished"
	if err != nil {
		slog.Error(msg, bag.fields...)
	} else {
		slog.Info(msg, bag.fields...)
	}
	return err
}
