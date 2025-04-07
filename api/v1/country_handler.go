package v1

import (
	"encoding/json"
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
	h.Log.Debug(ctx, "handler.GetCountries", "receive get all countries request")

	countries, err := h.Usecase.GetCountries(ctx)
	if err != nil {
		h.Log.Debug(ctx, "handler.GetCountries", "receive get countries with error request", err)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}

	response := utils.StatusOK(countries)
	utils.RespondWithJSON(writer, http.StatusOK, response)
}

func (h CountryHandler) GetCountryById(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	countryById := request.URL.Query().Get("countryById")
	id, err := strconv.Atoi(countryById)
	if err != nil {
		h.Log.Warn(ctx, "handler.getcountrybyid", "receive get country by id: %+v with error", countryById)

		responseErr := utils.StatusInternalServerError(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusInternalServerError, responseErr)
		return
	}

	h.Log.Warn(ctx, "handler.getcountrybyid", "receive get country by id: %+v request", id)

	country, err := h.Usecase.GetCountryByID(ctx, uint16(id))
	if err != nil {
		h.Log.Warn(ctx, "handler.getcountrybyid", "failed to parse request body", err)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}

	response := utils.StatusOK(country)
	utils.RespondWithJSON(writer, http.StatusOK, response)
}

func (h CountryHandler) CreateCountry(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	h.Log.Debug(ctx, "handler.createcountry", "receive create country request")

	createRequestAuthor := new(countries.CreateCountryRequest)
	err := json.NewDecoder(request.Body).Decode(&createRequestAuthor)
	if err != nil {
		h.Log.Debug(ctx, "handler.createcountry", "failed to parse create country with error request", err)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}

	createAuthor, err := h.Usecase.CreateCountry(ctx, createRequestAuthor)
	if err != nil {
		h.Log.Warn(ctx, "handler.createcountry", "invalid to parse create country with error request", err)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}
	response := utils.StatusOK(createAuthor)
	utils.RespondWithJSON(writer, http.StatusOK, response)
}

func (h CountryHandler) UpdateCountry(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	h.Log.Debug(ctx, "handler.updatecountry", "receive update country request")

	countryById := request.URL.Query().Get("countryById")
	id, err := strconv.Atoi(countryById)
	if err != nil {
		h.Log.Warn(ctx, "handler.updatecountry", "receive get country by id: %+v with error", countryById)

		responseErr := utils.StatusInternalServerError(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusInternalServerError, responseErr)
		return
	}

	updateRequestAuthor := new(countries.UpdateCountryRequest)
	err = json.NewDecoder(request.Body).Decode(&updateRequestAuthor)
	if err != nil {
		h.Log.Debug(ctx, "handler.updatecountry", "failed to parse create country with error request", err)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}
	updateRequestAuthor.ID = uint8(id)

	countryResponse, err := h.Usecase.UpdateCountry(ctx, updateRequestAuthor)
	if err != nil {
		h.Log.Warn(ctx, "handler.createcountry", "invalid to parse create country with error request", err)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}
	response := utils.StatusOK(countryResponse)
	utils.RespondWithJSON(writer, http.StatusOK, response)
}

func (h CountryHandler) DeleteAuthor(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	h.Log.Debug(ctx, "handler.deletecountry", "receive delete country request")

	countryById := request.URL.Query().Get("countryById")
	id, err := strconv.Atoi(countryById)
	if err != nil {
		h.Log.Warn(ctx, "handler.deletecountry", "receive get country by id: %+v with error", countryById)

		responseErr := utils.StatusInternalServerError(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusInternalServerError, responseErr)
		return
	}

	err = h.Usecase.DeleteCountry(ctx, uint16(id))
	if err != nil {
		h.Log.Warn(ctx, "handler.deletecountry", "failed to parse request body", err)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}

	response := utils.StatusOK(struct{}{})
	utils.RespondWithJSON(writer, http.StatusOK, response)
}
