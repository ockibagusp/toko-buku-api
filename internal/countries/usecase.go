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
	funcName := "usecase.GetCountries"

	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, &funcName, "failed request body to get countries: repo db begin", err)
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	countrys, err := u.Repo.GetCountries(ctx, tx)
	if err != nil {
		u.Log.Warn(ctx, &funcName, "failed request body to get countries: not found", err)
		return nil, err
	}

	return countrys, nil
}

func (u *Usecase) GetCountryByID(ctx context.Context, countryID uint16) (*Country, error) {
	funcName := "usecase.GetCountryByID"

	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, &funcName, "failed request body to get country by id", err)
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	country, err := u.Repo.GetCountryByID(ctx, tx, countryID)
	if err != nil {
		u.Log.Warn(ctx, &funcName, "failed request body to get country by id: not found by id", err)
		return nil, err
	}

	return country, nil
}

func (u *Usecase) CreateCountry(ctx context.Context, request *CreateCountryRequest) (*Country, error) {
	funcName := "usecase.CreateCountry"

	err := u.Validate.Struct(request)
	if err != nil {
		u.Log.Warn(ctx, &funcName, "invalid request body to create country", err)
		return nil, err
	}

	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, &funcName, "failed request body to create country: repo db begin", err)
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
		u.Log.Warn(ctx, &funcName, "invalid request body to create country", err)
		return nil, err
	}

	return createdCountry, nil
}

func (u *Usecase) UpdateCountry(ctx context.Context, request *UpdateCountryRequest) (*Country, error) {
	funcName := "usecase.UpdateCountry"

	err := u.Validate.Struct(request)
	if err != nil {
		u.Log.Warn(ctx, &funcName, "invalid request body to update country", err)
	}

	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, &funcName, "failed request body to update country: repo db begin", err)
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	updatedCountry, err := u.Repo.GetCountryByID(ctx, tx, uint16(request.ID))
	if err != nil {
		u.Log.Warn(ctx, &funcName, "failed request body to update country", err)
		return nil, err
	}

	updatedCountry.Country = request.Country

	updatedCountry, err = u.Repo.UpdateCountry(ctx, tx, updatedCountry)
	if err != nil {
		u.Log.Warn(ctx, &funcName, "failed request body to update country: repo UpdateCountry", err)
		return nil, err
	}

	return updatedCountry, nil
}

func (u *Usecase) DeleteCountry(ctx context.Context, countryID uint16) error {
	funcName := "usecase.DeleteCountry"

	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, &funcName, "failed request body to delete country: repo db begin", err)
		return err
	}
	defer utils.CommitOrRollback(tx)

	country, err := u.Repo.GetCountryByID(ctx, tx, countryID)
	if err != nil {
		u.Log.Warn(ctx, &funcName, "failed request body to delete country", err)
		return err
	}

	return u.Repo.DeleteCountry(ctx, tx, country)
}
