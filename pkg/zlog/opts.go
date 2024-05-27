package zlog

import (
	"io"

	"github.com/rs/zerolog"
)

type OptionFn func(opts *LoggerOptions)

type LoggerOptions struct {
	l      zerolog.Logger
	writer io.Writer
	hooks  []zerolog.Hook
}
