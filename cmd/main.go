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

	server := http.Server{
		Addr:         fmt.Sprintf("%s:%s", viper.GetString("server.host"), viper.GetString("server.port")),
		Handler:      routing,
		ReadTimeout:  time.Duration(viper.GetInt("server.readTimeout")) * time.Second,
		WriteTimeout: time.Duration(viper.GetInt("server.writeTimeout")) * time.Second,
		IdleTimeout:  time.Duration(viper.GetInt("server.idleTimeout")) * time.Second,
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
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", signalChan)

		ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			server.Close()
			log.Error(ctx, "startup", "err", fmt.Sprintf("could not stop server gracefully: %v", err))
		}
	}
}

func test(db *sql.DB) {
	var author authors.Author
	row := db.QueryRow("SELECT a.id, a.updated_at, a.author, a.city, c.id, c.updated_at, c.iso3, c.country, c.nice_country, c.currency FROM author a LEFT JOIN country c ON a.country_id = c.id WHERE a.id = ?", 1)

	var country countries.Country
	if err := row.Scan(
		&author.ID,
		&author.Updated_At,
		&author.Author,
		&author.City,
		&country.ID,
		&country.Updated_At,
		&country.Iso3,
		&country.Country,
		&country.Nice_Country,
		&country.Currency,
	); err != nil {
		if err == sql.ErrNoRows {
			fmt.Errorf("albumsById %d: no such album", 1)
		}
		fmt.Errorf("albumsById %d: %v", 1, err)
	}

	author.Country = &country

	fmt.Println(author, country)
}
