package v1

import (
	"encoding/json"
	"net/http"
	"strconv"
	"toko-buku-api/internal/authors"
	"toko-buku-api/pkg/logger"
	"toko-buku-api/utils"

	"github.com/go-playground/validator/v10"
)

// Handler for author-related endpoints

type AuthorHandler struct {
	Usecase authors.Usecase
	Log     *logger.Logger
}

func NewAuthorHandler(usercase authors.Usecase, logger *logger.Logger, validate *validator.Validate) *AuthorHandler {
	// file, err := os.OpenFile("./tmp/app-info.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	// if err != nil {
	//     log.Fatal(err)
	// }
	// log.SetOutput(file)
	// log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	// var log = logger.NewWithFiles(os.Stdout, logger.LevelDebug, "AUTHOR", nil)

	return &AuthorHandler{
		Usecase: usercase,
		Log:     logger,
	}
}

func (h AuthorHandler) GetAuthors(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	funcName := "handler.GetAuthors"
	h.Log.Debug(ctx, &funcName, "receive get all authors request")

	authors, err := h.Usecase.GetAuthors(ctx)
	if err != nil {
		h.Log.Debug(ctx, &funcName, "receive get authors with error request", err)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}

	response := utils.StatusOK(authors)
	utils.RespondWithJSON(writer, http.StatusOK, response)
}

func (h AuthorHandler) GetAuthorById(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	funcName := "handler.GetAuthorById"

	authorById := request.PathValue("authorById")
	h.Log.Debug(ctx, &funcName, "receive get author by id request", authorById)

	id, err := strconv.Atoi(authorById)
	if err != nil {
		h.Log.Warn(ctx, &funcName, "receive get author by id: %+v with error", authorById)

		responseErr := utils.StatusInternalServerError(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusInternalServerError, responseErr)
		return
	}

	h.Log.Warn(ctx, &funcName, "receive get author by id: %+v request", id)

	author, err := h.Usecase.GetAuthorById(ctx, uint16(id))
	if err != nil {
		h.Log.Warn(ctx, &funcName, "failed to parse request body", err)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}

	response := utils.StatusOK(author)
	utils.RespondWithJSON(writer, http.StatusOK, response)
}

// ///
func (h AuthorHandler) CreateAuthor(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	funcName := "handler.CreateAuthor"
	h.Log.Debug(ctx, &funcName, "receive create author request")

	createRequestAuthor := new(authors.CreateAuthorRequest)
	err := json.NewDecoder(request.Body).Decode(&createRequestAuthor)
	if err != nil {
		h.Log.Debug(ctx, &funcName, "failed to parse create author with error request", err)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}

	createAuthor, err := h.Usecase.CreateAuthor(ctx, createRequestAuthor)
	if err != nil {
		h.Log.Warn(ctx, &funcName, "invalid to parse create author with error request", err)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}

	response := utils.StatusOK(createAuthor)
	utils.RespondWithJSON(writer, http.StatusOK, response)
}

func (h AuthorHandler) UpdateAuthor(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	funcName := "handler.UpdateAuthor"
	h.Log.Debug(ctx, &funcName, "receive update author request")

	authorById := request.URL.Query().Get("authorById")
	id, err := strconv.Atoi(authorById)
	if err != nil {
		h.Log.Warn(ctx, &funcName, "receive get author by id: %+v with error", authorById)

		responseErr := utils.StatusInternalServerError(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusInternalServerError, responseErr)
		return
	}

	updateRequestAuthor := new(authors.UpdateAuthorRequest)
	err = json.NewDecoder(request.Body).Decode(&updateRequestAuthor)
	if err != nil {
		h.Log.Debug(ctx, &funcName, "failed to parse create author with error request", err)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}
	updateRequestAuthor.ID = uint16(id)

	authorResponse, err := h.Usecase.UpdateAuthor(ctx, updateRequestAuthor)
	if err != nil {
		h.Log.Warn(ctx, &funcName, "invalid to parse create author with error request", err)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}
	response := utils.StatusOK(authorResponse)
	utils.RespondWithJSON(writer, http.StatusOK, response)
}

func (h AuthorHandler) DeleteAuthor(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	funcName := "handler.DeleteAuthor"
	h.Log.Debug(ctx, &funcName, "receive delete author request")

	authorById := request.URL.Query().Get("authorById")
	id, err := strconv.Atoi(authorById)
	if err != nil {
		h.Log.Warn(ctx, &funcName, "receive get author by id: %+v with error", authorById)

		responseErr := utils.StatusInternalServerError(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusInternalServerError, responseErr)
		return
	}

	err = h.Usecase.DeleteAuthor(ctx, uint16(id))
	if err != nil {
		h.Log.Warn(ctx, &funcName, "failed to parse request body", err)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}

	response := utils.StatusOK(struct{}{})
	utils.RespondWithJSON(writer, http.StatusOK, response)
}
