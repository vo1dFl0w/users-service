package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	http_adaptor "github.com/vo1dFl0w/users-service/internal/app/adapters/http"
	"github.com/vo1dFl0w/users-service/internal/app/adapters/http/jwt"
	"github.com/vo1dFl0w/users-service/internal/app/adapters/storage/postgres"
	"github.com/vo1dFl0w/users-service/internal/app/config"
	"github.com/vo1dFl0w/users-service/internal/app/logger"
	"github.com/vo1dFl0w/users-service/internal/app/usecase/auth_usecase"
	"github.com/vo1dFl0w/users-service/internal/app/usecase/user_usecase"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := run(ctx); err != nil {
		log.Println(ctx, "startup", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	log := logger.LoadLogger(cfg)

	// TODO create databseDSN
	// TODO LaunchDB
	databaseDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.DBname,
		cfg.DB.Sslmode,
	)

	db, err := sql.Open("postgres", databaseDSN)
	if err != nil {
		return fmt.Errorf("failed to launch database: %w", err)
	}
	defer db.Close()

	store := postgres.New(db)

	tokenService := jwt.New([]byte(cfg.Secret))
	authService := auth_usecase.NewService(store, tokenService)
	userService := user_usecase.NewService(store)

	server := &http.Server{
		Addr:    cfg.HTTPaddr,
		Handler: http_adaptor.NewHandler(log, tokenService, authService, userService),
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	serverErr := make(chan error, 1)
	go func() {
		log.Info("server started", "host", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		} else {
			serverErr <- nil
		}
	}()

	select {
	case err := <-serverErr:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Info("shutdown server received", "signal", sig.String())
		shutdownCtx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Error("gracefull shutdown failed", "err", err)
			return err
		}
		log.Info("server gracefully stopped")
		return nil
	}
}
