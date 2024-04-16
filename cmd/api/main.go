package main

import (
	"context"
	"devbubble-api/internal/config"
	customRepository "devbubble-api/internal/custom-repository"
	db "devbubble-api/internal/repository"
	"devbubble-api/internal/service"
	handler "devbubble-api/internal/transport/http/handler"
	"devbubble-api/internal/transport/http/middleware/jwt"
	mwLogger "devbubble-api/internal/transport/http/middleware/logger"
	socket "devbubble-api/internal/transport/socket"
	"devbubble-api/pkg/logger"
	"devbubble-api/pkg/logger/sl"
	customValidator "devbubble-api/pkg/validator"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(mwLogger.New(log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(cors.Handler(cors.Options{

		AllowedOrigins: []string{"https://*", "http://*"},

		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	var pl *pgxpool.Pool
	var ctx = context.Background()

	pool, err := pgxpool.New(ctx, fmt.Sprintf("postgres://%s:%s@%s:%s/%s?timezone=Europe/Moscow", cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.DBName))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	pl = pool
	defer pl.Close()
	log.Info(fmt.Sprintf("DB connected [PORT]:%s [HOST]:%s", cfg.DB.Port, cfg.DB.Host))

	queries := db.New(pl)
	customQueries := customRepository.New(pl)

	validator := customValidator.New(validator.New())
	updater := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {

			return true
		}}

	//Services
	messageService := service.NewMessageService(queries, log, cfg)
	privateChatService := service.NewPrivateChatService(queries, log, cfg, pool, customQueries)
	emailService := service.NewEmailService(cfg.SMTP.From, cfg.SMTP.Password, cfg.SMTP.Host, cfg.SMTP.Port, log)
	authService := service.NewAuthService(queries, emailService, log, cfg)
	userService := service.NewUserService(queries, authService, log)
	tagService := service.NewTagService(queries, log, cfg)

	//Middleware
	jwtMiddleware := jwt.New(authService)

	//Handlers
	messageHandler := handler.NewMessageHandler(messageService, validator, jwtMiddleware, log)
	privateChatHandler := handler.NewChatHandler(privateChatService, validator, jwtMiddleware, log, messageService)
	tagHandler := handler.NewTagHandler(tagService, validator, jwtMiddleware, log)
	socketHandler := socket.NewSocketHandler(validator, jwtMiddleware, updater, log, queries)
	userHandler := handler.NewUserHandler(userService, authService, validator, log, jwtMiddleware)
	authHandler := handler.NewAuthHandler(authService, validator, jwtMiddleware, log)

	//Routes
	r.Mount("/user", userHandler.InitRouter())
	r.Mount("/ws", socketHandler.InitRouter())
	r.Mount("/auth", authHandler.InitRouter())
	r.Mount("/tags", tagHandler.InitRouter())
	r.Mount("/chat", privateChatHandler.InitRouter())
	r.Mount("/message", messageHandler.InitRouter())

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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	log.Info("server stopped")
}
