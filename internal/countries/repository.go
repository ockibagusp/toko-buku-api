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
	var countries []Country
	query := `SELECT id, updated_at, iso3, country, nice_country, currency FROM country`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		r.Log.Debug(ctx, "repository.GetCountries", "get query context with error", err)

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		country, err := scanIntoGetCountries(rows)
		if err != nil {
			return nil, err
		}

		countries = append(countries, country)
	}

	rerr := rows.Close()
	if rerr != nil {
		r.Log.Debug(ctx, "repository.GetCountries", "get query context close with error", err)

		return nil, err
	}

	if err := rows.Err(); err != nil {
		r.Log.Debug(ctx, "repository.GetCountries", "get query context .Err() with error", err)
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
	query := `SELECT id, updated_at, iso3, country, nice_country, currency FROM country WHERE id = ?`

	row := tx.QueryRowContext(ctx, query, countryID)

	country, err = scanRowIntoGetCountryByID(row)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf(countryNotFoundError, country.ID)
		}
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
	query := "INSERT into country(id, updated_at, iso3, country, nice_country, currency) VALUES (?)"
	result, err := tx.ExecContext(ctx, query, country.ID, country.Updated_At, country.Iso3, country.Country, country.Nice_Country, country.Currency)
	if err != nil {
		return nil, err
	}

	countryID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	country.ID = uint8(countryID)
	return country, nil
}

func (r Repository) UpdateCountry(ctx context.Context, tx *sql.Tx, country *Country) (*Country, error) {
	// TODO: Why?
	query := "UPDATE country set iso3 = ?, country = ?, nice_country = ?, currency = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, country.Iso3, country.Country, country.Currency, country.ID)
	if err != nil {
		return nil, err
	}

	return country, nil
}

func (r Repository) DeleteCountry(ctx context.Context, tx *sql.Tx, country *Country) error {
	query := "DELETE FROM country WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, country.ID)
	if err != nil {
		return err
	}

	return nil
}
