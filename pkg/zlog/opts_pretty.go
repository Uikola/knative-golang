package zlog

import (
	"time"

	"github.com/rs/zerolog"
)

func WithPretty(enable bool) OptionFn {
	return func(opts *LoggerOptions) {
		if !enable {
			return
		}

		opts.writer = zerolog.ConsoleWriter{ //nolint: exhaustruct // don't need
			Out:        opts.writer,
			TimeFormat: time.RFC3339,
		}
	}
}
