package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"toko-buku-api/internal/countries"
	"toko-buku-api/pkg/logger"
	"toko-buku-api/utils"

	"github.com/go-playground/validator/v10"
)

// Handler for country-related endpoints

type CountryHandler struct {
	Usecase countries.Usecase
	Log     *logger.Logger
}

func NewCountryHandler(usercase countries.Usecase, logger *logger.Logger, validate *validator.Validate) *CountryHandler {
	return &CountryHandler{
		Usecase: usercase,
		Log:     logger,
	}
}

func (h CountryHandler) GetCountries(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	funcName := "handler.GetCountries"
	h.Log.Info(ctx, "receive get countries request", "func_name", funcName)

	countries, err := h.Usecase.GetCountries(ctx)
	if err != nil {
		h.Log.Warn(ctx, "receive get countries with error request", "error", err, "func_name", funcName)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}

	h.Log.Info(ctx, "receive response to get countries response", "response", "waiting", "func_name", funcName)
	response := utils.StatusOK(countries)
	h.Log.Info(ctx, "receive response to get authors response", "response", "ok", "func_name", funcName)
	utils.RespondWithJSON(writer, http.StatusOK, response)
}

func (h CountryHandler) GetCountryById(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	funcName := "handler.GetCountryById"

	countryById := request.PathValue("countryById")
	h.Log.Info(ctx, fmt.Sprintf("receive request to get country by id: %+v", countryById), "func_name", funcName)

	id, err := strconv.Atoi(countryById)
	if err != nil {
		h.Log.Error(ctx, fmt.Sprintf("receive get country by id: %+v with error", countryById), "error", err, "func_name", funcName)

		responseErr := utils.StatusInternalServerError(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusInternalServerError, responseErr)
		return
	}

	h.Log.Info(ctx, fmt.Sprintf("receive get country by id: %+v", id), "func_name", funcName)
	country, err := h.Usecase.GetCountryByID(ctx, uint16(id))
	if err != nil {
		h.Log.Warn(ctx, "failed to parse request body", "error", err, "func_name", funcName)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}

	response := utils.StatusOK(country)
	h.Log.Info(ctx, fmt.Sprintf("receive response to get country by id: %+v", id), "response", "ok", "func_name", funcName)
	utils.RespondWithJSON(writer, http.StatusOK, response)
}

func (h CountryHandler) CreateCountry(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	funcName := "handler.CreateCountry"
	h.Log.Info(ctx, "receive request to create country", "func_name", funcName)

	createRequestAuthor := new(countries.CreateCountryRequest)
	err := json.NewDecoder(request.Body).Decode(&createRequestAuthor)
	if err != nil {
		h.Log.Warn(ctx, "failed to parse create country with error request", "error", err, "func_name", funcName)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}

	createAuthor, err := h.Usecase.CreateCountry(ctx, createRequestAuthor)
	if err != nil {
		h.Log.Warn(ctx, "invalid to parse create country with error request", err, "func_name", funcName)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}
	response := utils.StatusOK(createAuthor)
	h.Log.Info(ctx, "receive response to create country", "response", "ok", "func_name", funcName)
	utils.RespondWithJSON(writer, http.StatusOK, response)
}

func (h CountryHandler) UpdateCountry(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	funcName := "handler.UpdateCountry"

	countryById := request.PathValue("countryById")
	h.Log.Info(ctx, fmt.Sprintf("receive request to update country by id: %+v", countryById), "func_name", funcName)

	id, err := strconv.Atoi(countryById)
	if err != nil {
		h.Log.Warn(ctx, fmt.Sprintf("receive update country by id: %+v with error", countryById), "error", err, "func_name", funcName)

		responseErr := utils.StatusInternalServerError(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusInternalServerError, responseErr)
		return
	}

	updateRequestAuthor := new(countries.UpdateCountryRequest)
	err = json.NewDecoder(request.Body).Decode(&updateRequestAuthor)
	if err != nil {
		h.Log.Debug(ctx, "failed to parse update country with error request", "error", err, "func_name", funcName)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}
	updateRequestAuthor.ID = uint8(id)

	countryResponse, err := h.Usecase.UpdateCountry(ctx, updateRequestAuthor)
	if err != nil {
		h.Log.Warn(ctx, "invalid to parse update country with error request", "error", err, "func_name", funcName)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}
	response := utils.StatusOK(countryResponse)
	h.Log.Info(ctx, "receive response to update country", "response", "ok", "func_name", funcName)
	utils.RespondWithJSON(writer, http.StatusOK, response)
}

func (h CountryHandler) DeleteAuthor(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	funcName := "handler.DeleteCountry"

	countryById := request.PathValue("countryById")
	h.Log.Info(ctx, fmt.Sprintf("receive request to delete country by id: %+v", countryById), "func_name", funcName)

	id, err := strconv.Atoi(countryById)
	if err != nil {
		h.Log.Error(ctx, fmt.Sprintf("receive to delete country by id: %+v with error", countryById), "error", err, "func_name", funcName)

		responseErr := utils.StatusInternalServerError(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusInternalServerError, responseErr)
		return
	}

	err = h.Usecase.DeleteCountry(ctx, uint16(id))
	if err != nil {
		h.Log.Warn(ctx, "failed to parse request body", "error", err, "func_name", funcName)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}

	response := utils.StatusOK(struct{}{})
	h.Log.Info(ctx, fmt.Sprintf("receive response to delete country: %+v", id), "response", "ok", "func_name", funcName)
	utils.RespondWithJSON(writer, http.StatusOK, response)
}
