package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/oustrix/homeset/internal/config"
	"github.com/oustrix/homeset/internal/domain/users"
	"github.com/oustrix/homeset/internal/handlers"
	"github.com/oustrix/homeset/pkg/httpserver"
	"github.com/oustrix/homeset/pkg/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//
	// Flags
	//
	flags := parseFlags()

	//
	// Config
	//
	cfg, err := config.New(flags.configPath)
	if err != nil {
		panic(fmt.Errorf("failed to create config: %w", err))
	}

	//
	// Logger
	//
	logger.Configure(logger.Config{
		Writer: os.Stdout,
		Level:  cfg.Logger.Level,
	})

	slog.WarnContext(
		ctx,
		"logger configured",
		"level", cfg.Logger.Level,
	)

	//
	// Storage
	//
	storage, err := resolveStorage(ctx, storageConfig{
		DBMSName: cfg.DBMS,
		sqliteConfig: sqliteConfig{
			DSN: cfg.SQLite.DSN,
		},
	})
	defer storage.Close(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to create storage: %w", err))
	}

	slog.WarnContext(
		ctx,
		"storage created",
		"dbms", cfg.DBMS,
	)

	//
	// Use Cases
	//

	// Users
	createUser := users.NewCreateUser(users.CreateUserConfig{
		Storage: storage,
	})
	getUser := users.NewGetUser(users.GetUserConfig{
		Storage: storage,
	})

	slog.InfoContext(ctx, "use cases created")

	//
	// HTTP
	//
	router, err := handlers.NewRouter(handlers.RouterConfig{
		CreateUser: createUser,
		GetUser:    getUser,
	})
	if err != nil {
		panic(fmt.Errorf("failed to create router: %w", err))
	}

	httpServer := httpserver.New(
		router,
		httpserver.Port(cfg.HTTP.Port),
		httpserver.ReadTimeout(cfg.HTTP.ReadTimeout),
		httpserver.WriteTimeout(cfg.HTTP.WriteTimeout),
		httpserver.ShutdownTimeout(cfg.HTTP.ShutdownTimeout),
	)

	slog.WarnContext(ctx, "http server started", "port", cfg.HTTP.Port)

	//
	// Graceful shutdown
	//
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-interrupt:
		slog.ErrorContext(ctx, "got interrupt signal", "signal", sig)
	case err = <-httpServer.Notify():
		slog.ErrorContext(ctx, "got http server error", "error", err)
	case <-ctx.Done():
		slog.ErrorContext(ctx, "main context done", "error", ctx.Err())
	}

	err = httpServer.Shutdown()
	if err != nil {
		slog.ErrorContext(ctx, "failed to shutdown http server", "error", err)
	}

	slog.WarnContext(ctx, "app stopped")
}
