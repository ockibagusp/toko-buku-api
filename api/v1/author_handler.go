package v1

import (
	"encoding/json"
	"fmt"
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
	h.Log.Info(ctx, "receive get authors request", "func_name", funcName)

	authors, err := h.Usecase.GetAuthors(ctx)
	if err != nil {
		h.Log.Warn(ctx, "receive get authors with error request", "error", err, "func_name", funcName)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}

	h.Log.Info(ctx, "receive response to get authors response", "response", "waiting", "func_name", funcName)
	response := utils.StatusOK(authors)
	h.Log.Info(ctx, "receive response to get authors response", "response", "ok", "func_name", funcName)
	utils.RespondWithJSON(writer, http.StatusOK, response)
}

func (h AuthorHandler) GetAuthorById(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	funcName := "handler.GetAuthorById"

	authorById := request.PathValue("authorById")
	h.Log.Info(ctx, fmt.Sprintf("receive get author by id: %+v", authorById), "func_name", funcName)

	id, err := strconv.Atoi(authorById)
	if err != nil {
		h.Log.Error(ctx, fmt.Sprintf("receive get author by id: %+v with error", authorById), "error", err, "func_name", funcName)

		responseErr := utils.StatusInternalServerError(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusInternalServerError, responseErr)
		return
	}

	h.Log.Info(ctx, fmt.Sprintf("receive get author by id: %+v", id), "func_name", funcName)
	author, err := h.Usecase.GetAuthorById(ctx, uint16(id))
	if err != nil {
		h.Log.Warn(ctx, "failed to parse request body", "error", err, "func_name", funcName)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}

	response := utils.StatusOK(author)
	h.Log.Info(ctx, fmt.Sprintf("receive response to get author by id: %+v", id), "response", "ok", "func_name", funcName)
	utils.RespondWithJSON(writer, http.StatusOK, response)
}

// ///
func (h AuthorHandler) CreateAuthor(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	funcName := "handler.CreateAuthor"
	h.Log.Info(ctx, "receive request to create author", "func_name", funcName)

	createRequestAuthor := new(authors.CreateAuthorRequest)
	err := json.NewDecoder(request.Body).Decode(&createRequestAuthor)
	if err != nil {
		h.Log.Error(ctx, "failed to parse create author with error request", "error", err, "func_name", funcName)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}

	createAuthor, err := h.Usecase.CreateAuthor(ctx, createRequestAuthor)
	if err != nil {
		h.Log.Warn(ctx, "invalid to parse create author with error request", "error", err, "func_name", funcName)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}

	response := utils.StatusOK(createAuthor)
	h.Log.Info(ctx, "receive response to create author", "response", "ok", "func_name", funcName)
	utils.RespondWithJSON(writer, http.StatusOK, response)
}

func (h AuthorHandler) UpdateAuthor(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	funcName := "handler.UpdateAuthor"

	authorById := request.PathValue("authorById")
	h.Log.Debug(ctx, fmt.Sprintf("receive request to update author: %+v", authorById), "func_name", funcName)

	id, err := strconv.Atoi(authorById)
	if err != nil {
		h.Log.Warn(ctx, fmt.Sprintf("receive update author by id: %+v with error", authorById), "error", err, "func_name", funcName)

		responseErr := utils.StatusInternalServerError(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusInternalServerError, responseErr)
		return
	}

	updateRequestAuthor := new(authors.UpdateAuthorRequest)
	err = json.NewDecoder(request.Body).Decode(&updateRequestAuthor)
	if err != nil {
		h.Log.Debug(ctx, "failed to parse update author with error request", "error", err, "func_name", funcName)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}
	updateRequestAuthor.ID = uint16(id)

	authorResponse, err := h.Usecase.UpdateAuthor(ctx, updateRequestAuthor)
	if err != nil {
		h.Log.Warn(ctx, "invalid to parse update author with error request", "error", err, "func_name", funcName)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}

	response := utils.StatusOK(authorResponse)
	h.Log.Info(ctx, fmt.Sprintf("receive response to update author: %+v", id), "response", "ok", "func_name", funcName)
	utils.RespondWithJSON(writer, http.StatusOK, response)
}

func (h AuthorHandler) DeleteAuthor(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	funcName := "handler.DeleteAuthor"
	h.Log.Debug(ctx, "receive request to delete author", "func_name", funcName)
	authorById := request.PathValue("authorById")
	id, err := strconv.Atoi(authorById)
	if err != nil {
		h.Log.Error(ctx, fmt.Sprintf("receive get delete author by id: %+v with error", authorById), "error", err, "func_name", funcName)

		responseErr := utils.StatusInternalServerError(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusInternalServerError, responseErr)
		return
	}

	err = h.Usecase.DeleteAuthor(ctx, uint16(id))
	if err != nil {
		h.Log.Warn(ctx, "failed to parse request body", "error", err, "func_name", funcName)

		responseErr := utils.StatusBadRequest(err.Error())
		utils.RespondErrorWithJSON(writer, http.StatusBadRequest, responseErr)
		return
	}

	response := utils.StatusOK(struct{}{})
	h.Log.Info(ctx, fmt.Sprintf("receive response to delete author: %+v", id), "response", "ok", "func_name", funcName)
	utils.RespondWithJSON(writer, http.StatusOK, response)
}
