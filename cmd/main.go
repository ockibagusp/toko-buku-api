package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"toko-buku-api/config"
	"toko-buku-api/pkg/logger"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
)

// Entrypoints for different executables (e.g., main.go)

// Entrypoint for the main application

func main() {
	viper := config.NewViper()
	log := logger.New(os.Stdout, logger.LevelDebug, "MAIN", nil)
	db := config.NewDatabase(viper, log)
	validate := validator.New()

	routing := config.NewApp(&config.AppConfig{
		Viper:    viper,
		DB:       db,
		Log:      log,
		Validate: validate,
	})

	// ---
	// Start API Service
	ctx := context.Background()

	log.Info(ctx, "startup", "status", "initializing V1 API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	host := viper.GetString("server.host")
	readTimeout := viper.GetInt("server.readTimeout")
	writeTimeout := viper.GetInt("server.writeTimeout")
	idleTimeout := viper.GetInt("server.idleTimeout")

	server := http.Server{
		Addr:         fmt.Sprintf("%s:%s", host, viper.GetString("server.port")),
		Handler:      routing,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		IdleTimeout:  time.Duration(idleTimeout) * time.Second,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrs := make(chan error, 1)
	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", server.Addr)

		serverErrs <- server.ListenAndServe()
	}()

	// test(db)

	// -------------------------------------------------------------------------
	// Shutdown

	var shutdownTimeout = 20 * time.Millisecond

	select {
	case err := <-serverErrs:
		log.Error(ctx, "startup", "err", fmt.Sprintf("server error: %v", err))

	case signalChan := <-shutdown:
		fmt.Println("\n------")
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", signalChan)
		defer func() {
			log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", signalChan)
			os.Exit(0)
		}()

		log.Info(ctx, "shutdown", "status", "database close started", "status", "waiting...")
		defer func() {
			err := db.Close()
			if err != nil {
				log.Fatal(context.Background(), "got error when closing the DB connection", err)
			}
			log.Info(ctx, "shutdown", "status", "database close complete", "status", "ok")
		}()

		ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			server.Close()
			log.Fatal(ctx, "startup", "err", fmt.Sprintf("could not stop server gracefully: %v", err))
		}
	}
}

func test(db *sql.DB) {

}
