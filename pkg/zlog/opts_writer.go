package zlog

import (
	"io"
)

func WithWriter(writer io.Writer) OptionFn {
	return func(opts *LoggerOptions) {
		opts.writer = writer
	}
}
