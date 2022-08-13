package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"tschwaa.com/api/server"
	"tschwaa.com/api/storage"
)

// release is set through the linker at build time, generally from a git sha.
// Used for logging and error reporting
var release string

func main() {
	os.Exit(start())
}

func start() int {
	logEnv := getStringOrDefault("LOG_ENV", "development")
	log, err := createLogger(logEnv)
	if err != nil {
		fmt.Println("Error setting up the logger: ", err)
		return 1
	}

	log = log.With(zap.String("release", release))

	defer func() {
		// If we cannot sync, there's probably something wrong with outputting logs,
		// so we probably cannot write using fmt.Println either. So just ignore the error.
		_ = log.Sync()
	}()

	host := getStringOrDefault("HOST", "localhost")
	port := getIntOrDefault("PORT", 8080)

	s := server.New(server.Options{
		Database: createDatabase(log),
		Host:     host,
		Port:     port,
		Log:      log,
	})

	var eg errgroup.Group
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	eg.Go(func() error {
		<-ctx.Done()
		if err := s.Stop(); err != nil {
			log.Info("Error stopping server", zap.Error(err))
			return err
		}
		return nil
	})

	if err := s.Start(); err != nil {
		log.Info("Error starting server", zap.Error(err))
		return 1
	}

	if err := eg.Wait(); err != nil {
		return 1
	}

	return 0
}

func createLogger(env string) (*zap.Logger, error) {
	switch env {
	case "production":
		return zap.NewProduction()
	case "development":
		return zap.NewDevelopment()
	default:
		return zap.NewNop(), nil
	}
}

func createDatabase(log *zap.Logger) *storage.Database {
	return storage.NewDatabase(storage.NewDatabaseOptions{
		Host:                  getStringOrDefault("DB_HOST", "localhost"),
		Port:                  getIntOrDefault("DB_PORT", 5433),
		User:                  getStringOrDefault("DB_USER", "schwaa"),
		Password:              getStringOrDefault("DB_PASSWORD", "123"),
		Name:                  getStringOrDefault("DB_NAME", "schwaa"),
		MaxOpenConnections:    getIntOrDefault("DB_MAX_OPEN_CONNECTION", 10),
		MaxIdleConnections:    getIntOrDefault("DB_MAX_OPEN_CONNECTION", 10),
		ConnectionMaxLifetime: getDurationOrDefault("DB_CONNECTION_MAX_LIFETIME", time.Hour),
		Log:                   log,
	})
}

func getStringOrDefault(name, defaultV string) string {
	v, ok := os.LookupEnv(name)
	if !ok {
		return defaultV
	}
	return v
}

func getIntOrDefault(name string, defaultV int) int {
	v, ok := os.LookupEnv(name)
	if !ok {
		return defaultV
	}

	vAsInt, err := strconv.Atoi(v)
	if err != nil {
		return defaultV
	}
	return vAsInt
}

func getDurationOrDefault(name string, defaultV time.Duration) time.Duration {
	v, ok := os.LookupEnv(name)
	if !ok {
		return defaultV
	}
	vAsDuration, err := time.ParseDuration(v)
	if err != nil {
		return defaultV
	}
	return vAsDuration
}
