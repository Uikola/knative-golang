package zlog_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Uikola/knative-golang/pkg/zlog"
)

type logRecord struct {
	Level   string    `json:"level"`
	Version string    `json:"version"`
	Time    time.Time `json:"time"`
	Caller  string    `json:"caller"`
	Message string    `json:"message"`
}

func TestOutput(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := zlog.New(
		zlog.WithWriter(&buf),
		zlog.WithPretty(false),
		zlog.WithCaller(0),
		zlog.WithVersion("dev"),
	)
	log.Err(fmt.Errorf("test")).Msg("hello")

	rec := &logRecord{}
	require.NoError(t, json.NewDecoder(&buf).Decode(rec))
	require.NotEmpty(t, rec.Time)
	require.Equal(t, "error", rec.Level)
	require.Equal(t, "dev", rec.Version)
	require.Equal(t, "zlog/logger_test.go:33", rec.Caller)
	require.Equal(t, "hello", rec.Message)
}
