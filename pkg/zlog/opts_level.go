package zlog

import (
	"github.com/rs/zerolog"
)

func WithLevel(lvl zerolog.Level) OptionFn {
	return func(opts *LoggerOptions) {
		opts.l = opts.l.Level(lvl)
	}
}
