package mock

import (
	"context"
	"toko-buku-api/internal/authors"

	"github.com/stretchr/testify/mock"
)

type AuthorUsecaseMock struct {
	mock.Mock
}

func (m *AuthorUsecaseMock) GetAuthors(ctx context.Context) (*[]authors.Authors, error) {
	m.Called(ctx)
	return &[]authors.Authors{
		{
			ID:     1,
			Author: "test 1",
		},
		{
			ID:     2,
			Author: "test 2",
		},
	}, nil
}

func (m *AuthorUsecaseMock) GetAuthorById(ctx context.Context, authorId uint16) (*authors.Authors, error) {
	m.Called(ctx, authorId)
	return &authors.Authors{
		ID:     1,
		Author: "test",
	}, nil
}

// Add stub implementations for all other methods required by authors.Usecase interface
func (m *AuthorUsecaseMock) CreateAuthor(ctx context.Context, author *authors.CreateAuthorRequest) (*authors.Authors, error) {
	m.Called(ctx, author)
	return &authors.Authors{
		ID:     1,
		Author: "new test",
	}, nil
}

func (m *AuthorUsecaseMock) UpdateAuthor(context.Context, *authors.UpdateAuthorRequest) (*authors.Authors, error) {
	return &authors.Authors{
		ID:     1,
		Author: "update test",
	}, nil
}

func (m *AuthorUsecaseMock) DeleteAuthor(ctx context.Context, authorId uint16) error {
	m.Called(ctx, authorId)
	return nil
}
