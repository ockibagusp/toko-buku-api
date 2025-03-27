package v1

import (
	"encoding/json"
	"net/http"
	"os"
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

func NewAuthorHandler(usercase authors.Usecase, validate *validator.Validate) *AuthorHandler {
	// file, err := os.OpenFile("./tmp/app-info.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	// if err != nil {
	//     log.Fatal(err)
	// }
	// log.SetOutput(file)
	// log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	// var log = logger.NewWithFiles(os.Stdout, logger.LevelDebug, "AUTHOR", nil)
	log := logger.New(os.Stdout, logger.LevelDebug, "AUTHOR", nil)

	return &AuthorHandler{
		Usecase: usercase,
		Log:     log,
	}
}

func (h AuthorHandler) GetAuthors(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	h.Log.Debug(ctx, "handler.GetAuthors", "receive get all authors request")

	authors, err := h.Usecase.GetAuthors(ctx)
	if err != nil {
		h.Log.Debug(ctx, "handler.GetAuthors", "receive get authors with error request", err)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}

	response := utils.StatusOK(authors)
	utils.RespondWithJSON(writer, http.StatusOK, response)
}

func (h AuthorHandler) GetAuthorById(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	authorById := request.PathValue("authorById")
	id, err := strconv.Atoi(authorById)
	if err != nil {
		h.Log.Warn(ctx, "handler.GetAuthorById", "receive get author by id: %+v with error", authorById)

		responseErr := utils.StatusInternalServerError(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusInternalServerError, responseErr)
		return
	}

	h.Log.Warn(ctx, "handler.GetAuthorById", "receive get author by id: %+v request", id)

	author, err := h.Usecase.GetAuthorById(ctx, uint16(id))
	if err != nil {
		h.Log.Warn(ctx, "handler.GetAuthorById", "failed to parse request body", err)

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
	h.Log.Debug(ctx, "handler.CreateAuthor", "receive create author request")

	createRequestAuthor := new(authors.CreateAuthorRequest)
	err := json.NewDecoder(request.Body).Decode(&createRequestAuthor)
	if err != nil {
		h.Log.Debug(ctx, "handler.CreateAuthor", "failed to parse create author with error request", err)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}

	createAuthor, err := h.Usecase.CreateAuthor(ctx, createRequestAuthor)
	if err != nil {
		h.Log.Warn(ctx, "handler.CreateAuthor", "invalid to parse create author with error request", err)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}

	response := utils.StatusOK(createAuthor)
	utils.RespondWithJSON(writer, http.StatusOK, response)
}

func (h AuthorHandler) UpdateAuthor(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	h.Log.Debug(ctx, "handler.UpdateAuthor", "receive update author request")

	authorById := request.URL.Query().Get("authorById")
	id, err := strconv.Atoi(authorById)
	if err != nil {
		h.Log.Warn(ctx, "handler.UpdateAuthor", "receive get author by id: %+v with error", authorById)

		responseErr := utils.StatusInternalServerError(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusInternalServerError, responseErr)
		return
	}

	updateRequestAuthor := new(authors.UpdateAuthorRequest)
	err = json.NewDecoder(request.Body).Decode(&updateRequestAuthor)
	if err != nil {
		h.Log.Debug(ctx, "handler.UpdateAuthor", "failed to parse create author with error request", err)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}
	updateRequestAuthor.ID = uint16(id)

	authorResponse, err := h.Usecase.UpdateAuthor(ctx, updateRequestAuthor)
	if err != nil {
		h.Log.Warn(ctx, "handler.UpdateAuthor", "invalid to parse create author with error request", err)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}
	response := utils.StatusOK(authorResponse)
	utils.RespondWithJSON(writer, http.StatusOK, response)
}

func (h AuthorHandler) DeleteAuthor(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	h.Log.Debug(ctx, "handler.DeleteAuthor", "receive delete author request")

	authorById := request.URL.Query().Get("authorById")
	id, err := strconv.Atoi(authorById)
	if err != nil {
		h.Log.Warn(ctx, "handler.DeleteAuthor", "receive get author by id: %+v with error", authorById)

		responseErr := utils.StatusInternalServerError(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusInternalServerError, responseErr)
		return
	}

	err = h.Usecase.DeleteAuthor(ctx, uint16(id))
	if err != nil {
		h.Log.Warn(ctx, "handler.DeleteAuthor", "failed to parse request body", err)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}

	response := utils.StatusOK(struct{}{})
	utils.RespondWithJSON(writer, http.StatusOK, response)
}
