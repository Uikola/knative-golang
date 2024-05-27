package zlog

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func New(options ...OptionFn) zerolog.Logger {
	l := log.With().Logger()

	opts := &LoggerOptions{ //nolint: exhaustruct // don't need
		writer: os.Stdout,
		l:      l,
	}

	for _, fn := range options {
		fn(opts)
	}

	l = opts.l.Output(opts.writer)

	for _, h := range opts.hooks {
		l = l.Hook(h)
	}

	return l
}

func Default(color bool, version string, level zerolog.Level) zerolog.Logger {
	return New(defaultOpts(color, version, level)...)
}

func defaultOpts(color bool, version string, level zerolog.Level) []OptionFn {
	return []OptionFn{
		WithVersion(version),
		WithLevel(level),
		WithPretty(color),
	}
}
