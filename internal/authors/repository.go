package authors

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"toko-buku-api/internal/countries"
	"toko-buku-api/pkg/logger"
)

// Database access methods for author data
type Repository interface {
	GetDB() *sql.DB
	GetLog() *logger.Logger
	GetAuthors(ctx context.Context, tx *sql.Tx) (*[]Authors, error)
	GetAuthorById(ctx context.Context, tx *sql.Tx, authorId uint16) (*Authors, error)
	CreateAuthor(ctx context.Context, tx *sql.Tx, author *Authors) (*Authors, error)
	UpdateAuthor(ctx context.Context, tx *sql.Tx, author *Authors) (*Authors, error)
	DeleteAuthor(ctx context.Context, tx *sql.Tx, author *Authors) error
}

type repository struct {
	DB  *sql.DB
	Log *logger.Logger
}

const (
	authorBaseError     = "author %d: %v"
	authorNotFoundError = "author %d: not found"
)

func NewRepository(db *sql.DB, logger *logger.Logger) Repository {
	return repository{
		DB:  db,
		Log: logger,
	}
}

// FindAllWithJoin(..)
// FindAllComplete(..)
// FindCompletePersonByID(..)

func (r repository) GetDB() *sql.DB {
	return r.DB
}

func (r repository) GetLog() *logger.Logger {
	return r.Log
}

func (r repository) GetAuthors(ctx context.Context, tx *sql.Tx) (*[]Authors, error) {
	var authors []Authors
	var funcName = "repository.GetAuthors"
	query := `SELECT a.*, c.* FROM authors a LEFT JOIN countries c ON a.country_id = c.id`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		r.Log.Error(ctx, "get query context with error", "error", err, "func_name", funcName)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		author, err := scanIntoGetAuthors(rows)
		if err != nil {
			r.Log.Error(ctx, "get scan into get author with error", "error", err, "func_name", funcName)
			return nil, err
		}

		authors = append(authors, author)
	}

	if err := rows.Close(); err != nil {
		r.Log.Error(ctx, "get rows close with error", "error", err, "func_name", funcName)
		return nil, err
	}

	if err := rows.Err(); err != nil {
		r.Log.Error(ctx, "get rows err with error", "error", err, "func_name", funcName)
		return nil, err
	}

	return &authors, nil
}

func scanIntoGetAuthors(rows *sql.Rows) (selectedAuthor Authors, err error) {
	var country countries.Countries
	err = rows.Scan(
		&selectedAuthor.ID,
		&selectedAuthor.Updated_At,
		&selectedAuthor.Country_Id,
		&selectedAuthor.Author,
		&selectedAuthor.City,
		&country.ID,
		&country.Updated_At,
		&country.Iso3,
		&country.Country,
		&country.Nice_Country,
		&country.Currency,
	)

	// TODO
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return selectedAuthor, fmt.Errorf(authorNotFoundError, selectedAuthor.ID)
		}
		return selectedAuthor, fmt.Errorf(authorBaseError, selectedAuthor.ID, err)
	}

	selectedAuthor.Country = &country
	return
}

func (r repository) GetAuthorById(ctx context.Context, tx *sql.Tx, authorId uint16) (*Authors, error) {
	funcName := "repository.GetAuthorById"
	query := `SELECT a.*, c.* FROM authors a LEFT JOIN countries c ON a.country_id = c.id WHERE a.id = ?`

	row := tx.QueryRowContext(ctx, query, authorId)

	author, err := scanRowIntoGetAuthorById(row, authorId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.Log.Error(ctx, "get scan row into get author by id with errors.Is", "error", err, "func_name", funcName)
			return nil, errors.New("not found")
		}
		r.Log.Error(ctx, "get scan row into get author by id with error", "error", err, "func_name", funcName)
		return nil, errors.New("not found")
	}

	return author, nil
}

func scanRowIntoGetAuthorById(row *sql.Row, authorId uint16) (*Authors, error) {
	var author = Authors{}
	var country = countries.Countries{}

	if err := row.Scan(
		&author.ID,
		&author.Updated_At,
		&author.Country_Id,
		&author.Author,
		&author.City,
		&country.ID,
		&country.Updated_At,
		&country.Iso3,
		&country.Country,
		&country.Nice_Country,
		&country.Currency,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf(authorNotFoundError, authorId)
		}
		return nil, fmt.Errorf(authorBaseError, authorId, err)
	}

	author.Country = &country
	return &author, nil
}

func (r repository) CreateAuthor(ctx context.Context, tx *sql.Tx, author *Authors) (auther *Authors, err error) {
	funcName := "repository.CreateAuthor"

	query := "INSERT INTO authors(country_id, author, city) VALUES (?, ?, ?)"
	result, err := tx.ExecContext(ctx, query, author.Country_Id, author.Author, author.City)
	if err != nil {
		r.Log.Error(ctx, "get exec context with create author error", "error", err, "func_name", funcName)
		return nil, err
	}

	authorId, err := result.LastInsertId()
	if err != nil {
		r.Log.Error(ctx, "get result last insert id with create author error", "error", err, "func_name", funcName)
		return nil, err
	}
	author.ID = uint16(authorId)
	return author, nil
}

func (r repository) UpdateAuthor(ctx context.Context, tx *sql.Tx, author *Authors) (*Authors, error) {
	query := "UPDATE authors SET author = ?, city = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, author.Author, author.City, author.ID)
	if err != nil {
		r.Log.Error(ctx, "get exec context with update author error", "error", err, "func_name", "repository.UpdateAuthor")
		return nil, err
	}

	return author, nil
}

func (r repository) DeleteAuthor(ctx context.Context, tx *sql.Tx, author *Authors) error {
	query := "DELETE FROM authors WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, author.ID)
	if err != nil {
		r.Log.Error(ctx, "get exec context with delete author error", "error", err, "func_name", "repository.DeleteAuthor")
		return err
	}

	return nil
}
