package main

import (
	"context"
	"errors"
	"github.com/Uikola/knative-golang/pkg/zlog"
	"net/http"
	"os"
	"sync"

	server "github.com/Uikola/knative-golang/internal/server/http"
	"github.com/rs/zerolog/log"
)

const (
	DEBUGLEVEL = 0
)

func main() {
	log.Logger = zlog.Default(true, "dev", DEBUGLEVEL)

	srv := server.NewServer()

	httpServer := &http.Server{
		Addr:    ":8000",
		Handler: srv,
	}

	go func() {
		log.Info().Msg("Starting server...")
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			log.Error().Msg(err.Error())
			os.Exit(1)
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Error().Msg(err.Error())
			os.Exit(1)
		}
	}()
	wg.Wait()
}
