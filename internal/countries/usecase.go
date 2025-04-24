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

func (u *Usecase) GetCountries(ctx context.Context) ([]Countries, error) {
	funcName := "usecase.GetCountries"

	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, "failed request body to get countries: repo db begin", "error", err, "func_name", funcName)
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	countrys, err := u.Repo.GetCountries(ctx, tx)
	if err != nil {
		u.Log.Warn(ctx, "failed request body to get countries: not found", "error", err, "func_name", funcName)
		return nil, err
	}

	return countrys, nil
}

func (u *Usecase) GetCountryByID(ctx context.Context, countryID uint16) (*Countries, error) {
	funcName := "usecase.GetCountryByID"

	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, "failed request body to get country by id", "error", err, "func_name", funcName)
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	country, err := u.Repo.GetCountryByID(ctx, tx, countryID)
	if err != nil {
		u.Log.Warn(ctx, "failed request body to get country by id: not found by id", "error", err, "func_name", funcName)
		return nil, err
	}

	return country, nil
}

func (u *Usecase) CreateCountry(ctx context.Context, request *CreateCountryRequest) (*Countries, error) {
	funcName := "usecase.CreateCountry"

	err := u.Validate.Struct(request)
	if err != nil {
		u.Log.Warn(ctx, "invalid request body to create country", "error", err, "func_name", funcName)
		return nil, err
	}

	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, "failed request body to create country: repo db begin", "error", err, "func_name", funcName)
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	createdCountry := &Countries{
		Iso3:         request.Iso3,
		Country:      request.Country,
		Nice_Country: request.Nice_Country,
		Currency:     request.Currency,
	}

	createdCountry, err = u.Repo.CreateCountry(ctx, tx, createdCountry)
	if err != nil {
		u.Log.Warn(ctx, "invalid request body to create country", "error", err, "func_name", funcName)
		return nil, err
	}

	return createdCountry, nil
}

func (u *Usecase) UpdateCountry(ctx context.Context, request *UpdateCountryRequest) (*Countries, error) {
	funcName := "usecase.UpdateCountry"

	err := u.Validate.Struct(request)
	if err != nil {
		u.Log.Warn(ctx, "invalid request body to update country", "error", err, "func_name", funcName)
		return nil, err
	}

	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, "failed request body to update country: repo db begin", "error", err, "func_name", funcName)
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	oldCountry, err := u.Repo.GetCountryByID(ctx, tx, uint16(request.ID))
	if err != nil {
		u.Log.Warn(ctx, "failed request body to update country", "error", err, "func_name", funcName)
		return nil, err
	}

	if request.Country != "" {
		oldCountry.Country = request.Country
	}

	if request.Currency != "" {
		oldCountry.Currency = request.Currency
	}
	if request.Iso3 != "" {
		oldCountry.Iso3 = request.Iso3
	}

	if request.Nice_Country != "" {
		oldCountry.Nice_Country = request.Nice_Country
	}

	updatedCountry, err := u.Repo.UpdateCountry(ctx, tx, oldCountry)
	if err != nil {
		u.Log.Warn(ctx, "failed request body to update country: repo UpdateCountry", "error", err, "func_name", funcName)
		return nil, err
	}

	return updatedCountry, nil
}

func (u *Usecase) DeleteCountry(ctx context.Context, countryID uint16) error {
	funcName := "usecase.DeleteCountry"

	tx, err := u.Repo.DB.Begin()
	if err != nil {
		u.Log.Warn(ctx, "failed request body to delete country: repo db begin", "error", err, "func_name", funcName)
		return err
	}
	defer utils.CommitOrRollback(tx)

	country, err := u.Repo.GetCountryByID(ctx, tx, countryID)
	if err != nil {
		u.Log.Warn(ctx, "failed request body to delete country", "error", err, "func_name", funcName)
		return err
	}

	return u.Repo.DeleteCountry(ctx, tx, country)
}
