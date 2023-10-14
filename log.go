package try

// A common handler: log the error

import (
	"context"
	"log/slog"
	"runtime"
	"time"
)

type logErrorOptions struct {
	Logger *slog.Logger
	Level  slog.Level
	Msg    string
}

func defaultLogErrorOptions() logErrorOptions {
	return logErrorOptions{
		Logger: slog.Default(),
		Level:  slog.LevelError,
		Msg:    "try got an error",
	}
}

type LogOption func(*logErrorOptions)

func WithLogger(logger *slog.Logger) LogOption {
	return func(opts *logErrorOptions) {
		opts.Logger = logger
	}
}

func WithLevel(level slog.Level) LogOption {
	return func(opts *logErrorOptions) {
		opts.Level = level
	}
}

func WithMsg(msg string) LogOption {
	return func(opts *logErrorOptions) {
		opts.Msg = msg
	}
}

// Log returns a ErrorHandleFunc that logs the error
// with the given logger, and returns the error unchanged.
//
// If options is nil, the default options will be used.
// Or else the first element of options will be used.
func Log(options ...LogOption) HandlerFunc {
	opts := defaultLogErrorOptions()
	for _, option := range options {
		option(&opts)
	}

	return func(err error) error {
		// To skip callers:
		// https://github.com/golang/go/issues/59145#issuecomment-1481920720

		if !opts.Logger.Enabled(context.Background(), opts.Level) {
			return err
		}
		var pcs [1]uintptr
		runtime.Callers(4, pcs[:]) // skip [Callers, Log.func1, Check]
		r := slog.NewRecord(time.Now(), opts.Level, opts.Msg, pcs[0])
		_ = opts.Logger.With(slog.String("err", err.Error())).Handler().Handle(context.Background(), r)

		return err
	}
}
