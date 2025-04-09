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

func (u *Usecase) GetAuthors(ctx context.Context) (*[]Author, error) {
	funcName := "usecase.GetAuthors"

	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, &funcName, "failed request body to get authors: repo db begin", err)
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	authors, err := u.Repo.GetAuthors(ctx, tx)
	if err != nil {
		u.Log.Warn(ctx, &funcName, "failed request body to get authors: not found", err)
		return nil, err
	}

	return authors, nil
}

func (u *Usecase) GetAuthorById(ctx context.Context, authorId uint16) (*Author, error) {
	funcName := "usecase.GetAuthorById"

	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, &funcName, "failed request body to get author by id: repo db begin", err)
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	author, err := u.Repo.GetAuthorById(ctx, tx, authorId)
	if err != nil {
		u.Log.Warn(ctx, &funcName, "failed request body to get author by id: not found by id", err)
		return nil, err
	}

	return author, nil
}

func (u *Usecase) CreateAuthor(ctx context.Context, request *CreateAuthorRequest) (*Author, error) {
	funcName := "usecase.CreateAuthor"

	err := u.Validate.Struct(request)
	if err != nil {
		u.Log.Warn(ctx, &funcName, "invalid request body to create author", err)
		return nil, err
	}

	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, &funcName, "failed request body to create author: repo db begin", err)
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	createdAuthor := &Author{
		Country_Id: request.Country_Id,
		Author:     request.Author,
		City:       request.City,
	}

	createdAuthor, err = u.Repo.CreateAuthor(ctx, tx, createdAuthor)
	if err != nil {
		u.Log.Warn(ctx, &funcName, "invalid request body to create author", err)
		return nil, err
	}

	return createdAuthor, nil
}

func (u *Usecase) UpdateAuthor(ctx context.Context, request *UpdateAuthorRequest) (*Author, error) {
	funcName := "usecase.UpdateAuthor"

	err := u.Validate.Struct(request)
	if err != nil {
		u.Log.Warn(ctx, &funcName, "invalid request body to update author", err)
	}

	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, &funcName, "failed request body to update: repo db begin", err)
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	updateAuthor, err := u.Repo.GetAuthorById(ctx, tx, uint16(request.ID))
	if err != nil {
		u.Log.Warn(ctx, &funcName, "failed request body to update author: repo GetAuthorById", err)
		return nil, err
	}

	updateAuthor.Author = request.Author

	// ??
	updateAuthor, err = u.Repo.UpdateAuthor(ctx, tx, updateAuthor)
	if err != nil {
		u.Log.Warn(ctx, &funcName, "failed request body to update author", err)
		return nil, err
	}

	return updateAuthor, nil
}

func (u *Usecase) DeleteAuthor(ctx context.Context, authorId uint16) error {
	funcName := "usecase.DeleteAuthor"

	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, &funcName, "failed request body to delete: repo db begin", err)
		return err
	}
	defer utils.CommitOrRollback(tx)

	author, err := u.Repo.GetAuthorById(ctx, tx, authorId)
	if err != nil {
		u.Log.Warn(ctx, &funcName, "failed request body to delete author", err)
		return err
	}

	return u.Repo.DeleteAuthor(ctx, tx, author)
}
