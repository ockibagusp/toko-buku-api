package authors

import (
	"context"
	"toko-buku-api/pkg/logger"
	"toko-buku-api/utils"

	"github.com/go-playground/validator/v10"
)

// Core business logic for author operations

type Usecase struct {
	Repo     Repository
	Log      *logger.Logger
	Validate *validator.Validate
}

func NewUsecase(repo Repository, logger *logger.Logger, validate *validator.Validate) Usecase {
	return Usecase{
		Repo:     repo,
		Log:      logger,
		Validate: validate,
	}
}

func (u *Usecase) GetAuthors(ctx context.Context) (*[]Authors, error) {
	funcName := "usecase.GetAuthors"

	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, "failed request body to get authors: repo db begin", "error", err, "func_name", funcName)
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	authors, err := u.Repo.GetAuthors(ctx, tx)
	if err != nil {
		u.Log.Warn(ctx, "failed request body to get authors: not found", "error", err, "func_name", funcName)
		return nil, err
	}

	return authors, nil
}

func (u *Usecase) GetAuthorById(ctx context.Context, authorId uint16) (*Authors, error) {
	funcName := "usecase.GetAuthorById"

	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, "failed request body to get author by id: repo db begin", "error", err, "func_name", funcName)
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	author, err := u.Repo.GetAuthorById(ctx, tx, authorId)
	if err != nil {
		u.Log.Warn(ctx, "failed request body to get author by id: not found by id", "error", err, "func_name", funcName)
		return nil, err
	}

	return author, nil
}

func (u *Usecase) CreateAuthor(ctx context.Context, request *CreateAuthorRequest) (*Authors, error) {
	funcName := "usecase.CreateAuthor"

	err := u.Validate.Struct(request)
	if err != nil {
		u.Log.Warn(ctx, "invalid request body to create author", "error", err, "func_name", funcName)
		return nil, err
	}

	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, "failed request body to create author: repo db begin", "error", err, "func_name", funcName)
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	createdAuthor := &Authors{
		Country_Id: request.Country_Id,
		Author:     request.Author,
		City:       request.City,
	}

	createdAuthor, err = u.Repo.CreateAuthor(ctx, tx, createdAuthor)
	if err != nil {
		u.Log.Warn(ctx, "invalid request body to create author", "error", err, "func_name", funcName)
		return nil, err
	}

	return createdAuthor, nil
}

func (u *Usecase) UpdateAuthor(ctx context.Context, request *UpdateAuthorRequest) (*Authors, error) {
	funcName := "usecase.UpdateAuthor"

	err := u.Validate.Struct(request)
	if err != nil {
		u.Log.Warn(ctx, "invalid request body to update author", "error", err, "func_name", funcName)
		return nil, err
	}

	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, "failed request body to update: repo db begin", "error", err, "func_name", funcName)
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	oldAuthor, err := u.Repo.GetAuthorById(ctx, tx, uint16(request.ID))
	if err != nil {
		u.Log.Warn(ctx, "failed request body to update author: repo GetAuthorById", "error", err, "func_name", funcName)
		return nil, err
	}

	if request.Author != "" {
		oldAuthor.Author = request.Author
	}

	if request.City != "" {
		oldAuthor.City = request.City
	}

	if request.Country_Id > 0 {
		oldAuthor.Country_Id = request.Country_Id
	}

	updatedAuthor, err := u.Repo.UpdateAuthor(ctx, tx, oldAuthor)
	if err != nil {
		u.Log.Warn(ctx, "failed request body to update author", "error", err, "func_name", funcName)
		return nil, err
	}

	return updatedAuthor, nil
}

func (u *Usecase) DeleteAuthor(ctx context.Context, authorId uint16) error {
	funcName := "usecase.DeleteAuthor"

	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, "failed request body to delete author: repo db begin", "error", err, "func_name", funcName)
		return err
	}
	defer utils.CommitOrRollback(tx)

	author, err := u.Repo.GetAuthorById(ctx, tx, authorId)
	if err != nil {
		u.Log.Warn(ctx, "failed request body to delete author", "error", err, "func_name", funcName)
		return err
	}

	return u.Repo.DeleteAuthor(ctx, tx, author)
}
