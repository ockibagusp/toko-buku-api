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
	"toko-buku-api/internal/authors"
	"toko-buku-api/internal/countries"
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
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", signalChan)
		defer func() {
			log.Error(ctx, "shutdown", "status", "shutdown complete", "signal", signalChan)
			os.Exit(0)
		}()

		ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			server.Close()
			log.Error(ctx, "startup", "err", fmt.Sprintf("could not stop server gracefully: %v", err))
		}
	}
}

func test(db *sql.DB) {
	// var author authors.Author
	author1 := authors.Author{
		ID:         1,
		Updated_At: time.Now(),
		Country_Id: 1,
		Author:     "Author 1",
		City:       "City 1",
	}

	fmt.Printf("%+v\n", author1)

	author2 := authors.Author{
		ID:         2,
		Updated_At: time.Now(),
		Country_Id: 2,
		Country: &countries.Country{
			ID:   1,
			Iso3: "IDN",
		},
		Author: "Author 1",
		City:   "City 1",
	}

	fmt.Printf("%+v\n", author2)
}
