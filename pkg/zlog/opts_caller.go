package zlog

import (
	"fmt"

	"github.com/rs/zerolog"
)

func WithCaller(skipFrame int) OptionFn {
	return func(opts *LoggerOptions) {
		zerolog.CallerMarshalFunc = CallerMarshalFunc
		if skipFrame == 0 {
			opts.l = opts.l.With().Caller().Logger()
		} else {
			opts.l = opts.l.With().CallerWithSkipFrameCount(skipFrame).Logger()
		}
	}
}

func CallerMarshalFunc(pc uintptr, file string, line int) string {
	short := file
	z := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			z++
			if z > 1 {
				break
			}
		}
	}
	file = short

	return fmt.Sprintf("%s:%d", file, line)
}
