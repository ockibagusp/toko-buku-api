package countries

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"toko-buku-api/pkg/logger"
)

// Database access methods for country data

type Repository struct {
	DB  *sql.DB
	Log *logger.Logger
}

const (
	countryBaseError     = "country  %d: %v"
	countryNotFoundError = "country %d: not found"
)

func NewRepository(db *sql.DB, logger *logger.Logger) Repository {
	return Repository{
		DB:  db,
		Log: logger,
	}
}

// FindAllWithJoin(..)
// FindAllComplete(..)
// FindCompletePersonByID(..)

func (r Repository) GetCountries(ctx context.Context, tx *sql.Tx) ([]Country, error) {
	var funcName = "repository.GetCountries"
	var countries []Country
	query := `SELECT id, updated_at, iso3, country, nice_country, currency FROM country`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		r.Log.Debug(ctx, &funcName, "get query context with error", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		country, err := scanIntoGetCountries(rows)
		if err != nil {
			r.Log.Debug(ctx, &funcName, "get scan into get countries with error", err)
			return nil, err
		}

		countries = append(countries, country)
	}

	rerr := rows.Close()
	if rerr != nil {
		r.Log.Debug(ctx, &funcName, "get rows close with error", err)

		return nil, err
	}

	if err := rows.Err(); err != nil {
		r.Log.Debug(ctx, &funcName, "get rows err with error", err)
		return nil, err
	}

	return countries, nil
}

func scanIntoGetCountries(rows *sql.Rows) (country Country, err error) {
	err = rows.Scan(
		&country.ID,
		&country.Updated_At,
		&country.Iso3,
		&country.Country,
		&country.Nice_Country,
		&country.Currency,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return country, fmt.Errorf(countryNotFoundError, country.ID)
		}
		return country, fmt.Errorf(countryBaseError, country.ID, err)
	}

	return
}

func (r Repository) GetCountryByID(ctx context.Context, tx *sql.Tx, countryID uint16) (country *Country, err error) {
	funcName := "repository.GetCountryByID"
	query := `SELECT id, updated_at, iso3, country, nice_country, currency FROM country WHERE id = ?`

	row := tx.QueryRowContext(ctx, query, countryID)

	country, err = scanRowIntoGetCountryByID(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.Log.Debug(ctx, &funcName, "get scan row into get coutry by id with errors.Is", err)
			return nil, fmt.Errorf(countryNotFoundError, country.ID)
		}
		r.Log.Debug(ctx, &funcName, "get scan row into get coutry by id with error", err)
		return nil, fmt.Errorf(countryBaseError, country.ID, err)
	}

	return country, nil
}

func scanRowIntoGetCountryByID(row *sql.Row) (*Country, error) {
	country := Country{}

	err := row.Scan(
		&country.ID,
		&country.Updated_At,
		&country.Iso3,
		&country.Country,
		&country.Nice_Country,
		&country.Currency,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf(countryNotFoundError, country.ID)
		}
		return nil, fmt.Errorf(countryBaseError, country.ID, err)
	}

	return &country, nil
}

func (r Repository) CreateCountry(ctx context.Context, tx *sql.Tx, country *Country) (auther *Country, err error) {
	funcName := "repository.CreateCountry"

	query := "INSERT into country(id, updated_at, iso3, country, nice_country, currency) VALUES (?)"
	result, err := tx.ExecContext(ctx, query, country.ID, country.Updated_At, country.Iso3, country.Country, country.Nice_Country, country.Currency)
	if err != nil {
		r.Log.Debug(ctx, &funcName, "get exec context with create country error", err)
		return nil, err
	}

	countryID, err := result.LastInsertId()
	if err != nil {
		r.Log.Debug(ctx, &funcName, "get result last insert id with create country error:", err)
		return nil, err
	}
	country.ID = uint8(countryID)
	return country, nil
}

func (r Repository) UpdateCountry(ctx context.Context, tx *sql.Tx, country *Country) (*Country, error) {
	funcName := "repository.UpdateCountry"

	// TODO: Why?
	query := "UPDATE country set iso3 = ?, country = ?, nice_country = ?, currency = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, country.Iso3, country.Country, country.Currency, country.ID)
	if err != nil {
		r.Log.Debug(ctx, &funcName, "get exec context with update country error", err)
		return nil, err
	}

	return country, nil
}

func (r Repository) DeleteCountry(ctx context.Context, tx *sql.Tx, country *Country) error {
	funcName := "repository.DeleteCountry"

	query := "DELETE FROM country WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, country.ID)
	if err != nil {
		r.Log.Debug(ctx, &funcName, "get exec context with delete country error", err)
		return err
	}

	return nil
}
