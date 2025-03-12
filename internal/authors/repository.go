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
	query := `SELECT a.*, c.* FROM author a LEFT JOIN country c ON a.country_id = c.id`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		r.Log.Debug(ctx, "repository.getauthors", "get query context with error", err)

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		author, err := scanIntoGetAuthors(rows)
		if err != nil {
			return nil, err
		}

		authors = append(authors, author)
	}

	if err := rows.Close(); err != nil {
		r.Log.Debug(ctx, "repository.getauthors", "get query context close with error", err)

		return nil, err
	}

	if err := rows.Err(); err != nil {
		r.Log.Debug(ctx, "repository.getauthors", "get query context .Err() with error", err)
		return nil, err
	}

	return &authors, nil
}

func scanIntoGetAuthors(rows *sql.Rows) (author Author, err error) {
	var country countries.Country

	err = rows.Scan(
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
	)

	// TODO
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return author, fmt.Errorf(authorNotFoundError, author.ID)
		}
		return author, fmt.Errorf(authorBaseError, author.ID, err)
	}

	author.Country = &country
	return
}

func (r Repository) GetAuthorById(ctx context.Context, tx *sql.Tx, authorId uint16) (*Author, error) {
	query := `SELECT a.*, c.* FROM author a LEFT JOIN country c ON a.country_id = c.id WHERE a.id = ?`

	row := tx.QueryRowContext(ctx, query, authorId)

	author, err := scanRowIntoGetAuthorById(row, authorId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf(authorNotFoundError, author.ID)
		}
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
	query := "INSERT into author(author, city) VALUES (?)"
	result, err := tx.ExecContext(ctx, query, author.Author, author.City)
	if err != nil {
		return nil, err
	}

	authorId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	author.ID = uint16(authorId)
	return author, nil
}

func (r Repository) UpdateAuthor(ctx context.Context, tx *sql.Tx, author *Author) (*Author, error) {
	query := "UPDATE author set author = ?, city = ? WHERE author_id = ?"
	_, err := tx.ExecContext(ctx, query, author.Author, author.City, author.ID)
	if err != nil {
		return nil, err
	}

	return author, nil
}

func (r Repository) DeleteAuthor(ctx context.Context, tx *sql.Tx, author *Author) error {
	query := "DELETE FROM author WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, author.ID)
	if err != nil {
		return err
	}

	return nil
}
