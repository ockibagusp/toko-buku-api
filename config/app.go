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
	// typeRepository := types.NewRepository(appConfig.DB, appConfig.Log)
	// typeService := service.NewTypeService(typeRepository, db, validate)
	// typeController := controller.NewTypeController(typeService)

	// bookRepository := repository.NewTypeRepository()
	// bookService := service.NewTypeService(bookRepository, db, validate)
	// bookController := controller.NewTypeController(bookService)

	mux := http.NewServeMux()

	// handle author-related endpoints
	authorRepository := authors.NewRepository(appConfig.DB, appConfig.Log)
	authorUsecase := authors.NewUsecase(authorRepository, appConfig.Log, appConfig.Validate)
	authorHandler := v1.NewAuthorHandler(authorUsecase, appConfig.Validate)
	mux.HandleFunc("GET /authors", authorHandler.GetAuthors)
	mux.HandleFunc("GET /authors/{authorById}", authorHandler.GetAuthorById)
	mux.HandleFunc("POST /authors", authorHandler.CreateAuthor)
	mux.HandleFunc("PUT /authors/{authorById}", authorHandler.UpdateAuthor)
	mux.HandleFunc("DELETE /authors/{authorById}", authorHandler.DeleteAuthor)

	// handle country-related endpoints
	countryRepository := countries.NewRepository(appConfig.DB, appConfig.Log)
	countryUsecase := countries.NewUsecase(countryRepository, appConfig.Log, appConfig.Validate)
	countryHandler := v1.NewCountryHandler(countryUsecase, appConfig.Validate)
	mux.HandleFunc("GET /countries", countryHandler.GetCountries)
	mux.HandleFunc("GET /countries/{countryById}", countryHandler.GetCountryById)
	mux.HandleFunc("POST /countries", countryHandler.CreateCountry)
	mux.HandleFunc("PUT /countries/{countryById}", countryHandler.UpdateCountry)
	mux.HandleFunc("DELETE /countries/{countryById}", countryHandler.DeleteAuthor)

	return mux
}
