package authors

import (
	"context"
	"fmt"
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
	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, "usecase.GetAuthors", "failed request body to get authors", err)
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	authors, err := u.Repo.GetAuthors(ctx, tx)
	if err != nil {
		u.Log.Warn(ctx, "usecase.GetAuthors", "not found", err)
		return nil, err
	}

	fmt.Println(authors)

	return authors, nil
}

func (u *Usecase) GetAuthorById(ctx context.Context, authorId uint16) (*Author, error) {
	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, "usecase.GetAuthorById", "failed request body to get author by id", err)
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	author, err := u.Repo.GetAuthorById(ctx, tx, authorId)
	if err != nil {
		u.Log.Warn(ctx, "usecase.GetAuthorById", "not found by id", err)
		return nil, err
	}

	return author, nil
}

func (u *Usecase) CreateAuthor(ctx context.Context, request *CreateAuthorRequest) (*Author, error) {
	err := u.Validate.Struct(request)
	if err != nil {
		u.Log.Warn(ctx, "usecase.CreateAuthor", "invalid request body to create author", err)
		return nil, err
	}

	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, "usecase.CreateAuthor", "failed request body to create author", err)
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
		u.Log.Warn(ctx, "usecase.CreateAuthor", "invalid request body to create author", err)
		return nil, err
	}

	return createdAuthor, nil
}

func (u *Usecase) UpdateAuthor(ctx context.Context, request *UpdateAuthorRequest) (*Author, error) {
	err := u.Validate.Struct(request)
	if err != nil {
		u.Log.Warn(ctx, "usecase.UpdateAuthor", "invalid request body to update author", err)
	}

	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, "usecase.UpdateAuthor", "failed request body to update author", err)
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	updateAuthor, err := u.Repo.GetAuthorById(ctx, tx, uint16(request.ID))
	if err != nil {
		u.Log.Warn(ctx, "usecase.UpdateAuthor", "failed request body to update author", err)
		return nil, err
	}

	updateAuthor.Author = request.Author

	// ??
	updateAuthor, err = u.Repo.UpdateAuthor(ctx, tx, updateAuthor)
	if err != nil {
		return nil, err
	}

	return updateAuthor, nil
}

func (u *Usecase) DeleteAuthor(ctx context.Context, authorId uint16) error {
	tx, err := u.Repo.DB.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

	author, err := u.Repo.GetAuthorById(ctx, tx, authorId)
	if err != nil {
		u.Log.Warn(ctx, "usecase.DeleteAuthor", "failed request body to update author", err)
		return err
	}

	return u.Repo.DeleteAuthor(ctx, tx, author)
}
