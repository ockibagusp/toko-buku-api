package countries

import (
	"context"
	"toko-buku-api/pkg/logger"
	"toko-buku-api/utils"

	"github.com/go-playground/validator/v10"
)

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

func (u *Usecase) GetCountries(ctx context.Context) ([]Country, error) {
	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, "usecase.GetCountries", "failed request body to get authors", err)
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	authors, err := u.Repo.GetCountries(ctx, tx)
	if err != nil {
		u.Log.Warn(ctx, "usecase.GetCountries", "not found", err)
		return nil, err
	}

	return authors, nil
}

func (u *Usecase) GetCountryByID(ctx context.Context, authorID uint16) (*Country, error) {
	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, "usecase.GetCountryByID", "failed request body to get author by id", err)
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	author, err := u.Repo.GetCountryByID(ctx, tx, authorID)
	if err != nil {
		u.Log.Warn(ctx, "usecase.GetCountryByID", "not found by id", err)
		return nil, err
	}

	return author, nil
}

func (u *Usecase) CreateCountry(ctx context.Context, request *CreateCountryRequest) (*Country, error) {
	err := u.Validate.Struct(request)
	if err != nil {
		u.Log.Warn(ctx, "usecase.CreateCountry", "invalid request body to create author", err)
		return nil, err
	}

	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, "usecase.CreateCountry", "failed request body to create author", err)
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	createdCountry := &Country{
		Iso3:         request.Iso3,
		Country:      request.Country,
		Nice_Country: request.Nice_Country,
		Currency:     request.Currency,
	}

	createdCountry, err = u.Repo.CreateCountry(ctx, tx, createdCountry)
	if err != nil {
		u.Log.Warn(ctx, "usecase.CreateCountry", "invalid request body to create author", err)
		return nil, err
	}

	return createdCountry, nil
}

func (u *Usecase) UpdateCountry(ctx context.Context, request *UpdateCountryRequest) (*Country, error) {
	err := u.Validate.Struct(request)
	if err != nil {
		u.Log.Warn(ctx, "usecase.UpdateCountry", "invalid request body to update author", err)
	}

	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, "usecase.UpdateCountry", "failed request body to update author", err)
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	updateCountry, err := u.Repo.GetCountryByID(ctx, tx, uint16(request.ID))
	if err != nil {
		u.Log.Warn(ctx, "usecase.UpdateCountry", "failed request body to update author", err)
		return nil, err
	}

	updateCountry.Country = request.Country

	updateCountry, err = u.Repo.UpdateCountry(ctx, tx, updateCountry)
	if err != nil {
		return nil, err
	}

	return updateCountry, nil
}

func (u *Usecase) DeleteCountry(ctx context.Context, authorID uint16) error {
	tx, err := u.Repo.DB.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

	author, err := u.Repo.GetCountryByID(ctx, tx, authorID)
	if err != nil {
		u.Log.Warn(ctx, "usecase.DeleteCountry", "failed request body to update author", err)
		return err
	}

	return u.Repo.DeleteCountry(ctx, tx, author)
}
