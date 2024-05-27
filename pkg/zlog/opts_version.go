package zlog

func WithVersion(version string) OptionFn {
	return func(opts *LoggerOptions) {
		opts.l = opts.l.With().Str("version", version).Logger()
	}
}
