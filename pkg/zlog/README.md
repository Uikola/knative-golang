# Usage

usage of native logger

```go
log := zlog.Default(
    // pretty/color print
    true,
    // app version
    "app_version",
    // logger level
    zerolog.InfoLevel,
)
```

opts example

```go
log := zlog.New(
    zlog.WithPretty(false),
    zlog.WithCaller(0),
    zlog.WithVersion("dev"),
)
logger.Info().Msg("hello")
```
