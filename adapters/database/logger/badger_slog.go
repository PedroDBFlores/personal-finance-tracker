package logger

import (
	"fmt"
	"log/slog"

	"github.com/dgraph-io/badger/v4"
)

type BadgerSlog struct {
	*slog.Logger
}

func NewBadgerSlog(logger *slog.Logger) badger.Logger {
	return &BadgerSlog{logger}
}

func (b *BadgerSlog) Errorf(format string, args ...any) {
	b.Error(fmt.Sprintf(format, args...))
}

func (b *BadgerSlog) Warningf(format string, args ...any) {
	b.Warn(fmt.Sprintf(format, args...))
}

func (b *BadgerSlog) Infof(format string, args ...any) {
	b.Info(fmt.Sprintf(format, args...))
}

func (b *BadgerSlog) Debugf(format string, args ...any) {
	b.Debug(fmt.Sprintf(format, args...))
}
