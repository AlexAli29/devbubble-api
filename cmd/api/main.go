package main

import (
	"context"
	"devbubble-api/internal/config"
	handler "devbubble-api/internal/transport/http/handler"
	mwLogger "devbubble-api/internal/transport/http/middleware/logger"
	"devbubble-api/pkg/logger"
	"devbubble-api/pkg/logger/sl"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(mwLogger.New(log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	userRouter := handler.NewUserHandler()
	authRouter := handler.NewAuthHandler()

	r.Mount("/user", userRouter)
	r.Mount("/auth", authRouter)

	apiRouter := chi.NewRouter()
	apiRouter.Mount("/api", r)

	log.Info("starting devbubble-api", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	log.Info("starting server", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      apiRouter,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	// TODO: move timeout to config
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	// TODO: close storage

	log.Info("server stopped")
}
