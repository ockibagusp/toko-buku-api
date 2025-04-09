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

type Repository struct {
	DB  *sql.DB
	Log *logger.Logger
}

const (
	authorBaseError     = "author  %d: %v"
	authorNotFoundError = "author %d: not found"
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

func (r Repository) GetAuthors(ctx context.Context, tx *sql.Tx) (*[]Author, error) {
	var authors []Author
	var funcName = "repository.GetAuthors"
	query := `SELECT a.*, c.* FROM author a LEFT JOIN country c ON a.country_id = c.id`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		r.Log.Debug(ctx, &funcName, "get query context with error", err)

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		author, err := scanIntoGetAuthors(rows)
		if err != nil {
			r.Log.Debug(ctx, &funcName, "get scan into get authors with error", err)
			return nil, err
		}

		authors = append(authors, author)
	}

	if err := rows.Close(); err != nil {
		r.Log.Debug(ctx, &funcName, "get rows close with error", err)
		return nil, err
	}

	if err := rows.Err(); err != nil {
		r.Log.Debug(ctx, &funcName, "get rows err with error", err)
		return nil, err
	}

	return &authors, nil
}

func scanIntoGetAuthors(rows *sql.Rows) (selectedAuthor Author, err error) {
	var country countries.Country
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

func (r Repository) GetAuthorById(ctx context.Context, tx *sql.Tx, authorId uint16) (*Author, error) {
	funcName := "repository.GetAuthorById"
	query := `SELECT a.*, c.* FROM author a LEFT JOIN country c ON a.country_id = c.id WHERE a.id = ?`

	row := tx.QueryRowContext(ctx, query, authorId)

	author, err := scanRowIntoGetAuthorById(row, authorId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.Log.Debug(ctx, &funcName, "get scan row into get author by id with errors.Is", err)
			return nil, fmt.Errorf(authorNotFoundError, author.ID)
		}
		r.Log.Debug(ctx, &funcName, "get scan row into get author by id with error", err)
		return nil, fmt.Errorf(authorBaseError, author.ID, err)
	}

	return author, nil
}

func scanRowIntoGetAuthorById(row *sql.Row, authorId uint16) (*Author, error) {
	var author = Author{}
	var country = countries.Country{}

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
			return nil, fmt.Errorf("albumsById %d: no such album", authorId)
		}
		return nil, fmt.Errorf("albumsById %d: %v", authorId, err)
	}

	author.Country = &country
	return &author, nil
}

///////

func (r Repository) CreateAuthor(ctx context.Context, tx *sql.Tx, author *Author) (auther *Author, err error) {
	funcName := "repository.CreateAuthor"

	query := "INSERT into author(country_id, author, city) VALUES (?, ?, ?)"
	result, err := tx.ExecContext(ctx, query, author.Country_Id, author.Author, author.City)
	if err != nil {
		r.Log.Debug(ctx, &funcName, "get exec context with create author error", err)
		return nil, err
	}

	authorId, err := result.LastInsertId()
	if err != nil {
		r.Log.Debug(ctx, &funcName, "get result last insert id with create author error", err)
		return nil, err
	}
	author.ID = uint16(authorId)
	return author, nil
}

func (r Repository) UpdateAuthor(ctx context.Context, tx *sql.Tx, author *Author) (*Author, error) {
	funcName := "repository.UpdateAuthor"

	query := "UPDATE author set author = ?, city = ? WHERE author_id = ?"
	_, err := tx.ExecContext(ctx, query, author.Author, author.City, author.ID)
	if err != nil {
		r.Log.Debug(ctx, &funcName, "get exec context with update author error", err)
		return nil, err
	}

	return author, nil
}

func (r Repository) DeleteAuthor(ctx context.Context, tx *sql.Tx, author *Author) error {
	funcName := "repository.DeleteAuthor"

	query := "DELETE FROM author WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, author.ID)
	if err != nil {
		r.Log.Debug(ctx, &funcName, "get exec context with delete author error", err)
		return err
	}

	return nil
}
