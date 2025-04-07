package config

import (
	"database/sql"
	"net/http"
	v1 "toko-buku-api/api/v1"
	"toko-buku-api/internal/authors"
	"toko-buku-api/internal/countries"
	"toko-buku-api/pkg/logger"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// Configurations files and setup

type AppConfig struct {
	Viper    *viper.Viper
	DB       *sql.DB
	Log      *logger.Logger
	Validate *validator.Validate
}

func NewApp(appConfig *AppConfig) *http.ServeMux {
	// typeRepository := types.NewRepository(appDB, appLog)
	// typeService := service.NewTypeService(typeRepository, db, validate)
	// typeController := controller.NewTypeController(typeService)

	// bookRepository := repository.NewTypeRepository()
	// bookService := service.NewTypeService(bookRepository, db, validate)
	// bookController := controller.NewTypeController(bookService)

	mux := http.NewServeMux()

	// handle author-related endpoints
	authorLog := logger.NewService("AUTHOR")

	authorRepository := authors.NewRepository(appConfig.DB, authorLog)
	authorUsecase := authors.NewUsecase(authorRepository, authorLog, appConfig.Validate)
	authorHandler := v1.NewAuthorHandler(authorUsecase, authorLog, appConfig.Validate)
	mux.HandleFunc("GET /authors", authorHandler.GetAuthors)
	mux.HandleFunc("GET /authors/{authorById}", authorHandler.GetAuthorById)
	mux.HandleFunc("POST /authors", authorHandler.CreateAuthor)
	mux.HandleFunc("PUT /authors/{authorById}", authorHandler.UpdateAuthor)
	mux.HandleFunc("DELETE /authors/{authorById}", authorHandler.DeleteAuthor)

	// handle country-related endpoints
	countryLog := logger.NewService("COUNTRY")

	countryRepository := countries.NewRepository(appConfig.DB, countryLog)
	countryUsecase := countries.NewUsecase(countryRepository, countryLog, appConfig.Validate)
	countryHandler := v1.NewCountryHandler(countryUsecase, countryLog, appConfig.Validate)
	mux.HandleFunc("GET /countries", countryHandler.GetCountries)
	mux.HandleFunc("GET /countries/{countryById}", countryHandler.GetCountryById)
	mux.HandleFunc("POST /countries", countryHandler.CreateCountry)
	mux.HandleFunc("PUT /countries/{countryById}", countryHandler.UpdateCountry)
	mux.HandleFunc("DELETE /countries/{countryById}", countryHandler.DeleteAuthor)

	return mux
}
