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
	args := m.Called(ctx)
	return args.Get(0).(*[]authors.Authors), args.Error(1)
}

func (m *AuthorUsecaseMock) GetAuthorById(ctx context.Context, authorId uint16) (*authors.Authors, error) {
	args := m.Called(ctx, authorId)
	return args.Get(0).(*authors.Authors), args.Error(1)
}

// Add stub implementations for all other methods required by authors.Usecase interface
func (m *AuthorUsecaseMock) CreateAuthor(ctx context.Context, author *authors.CreateAuthorRequest) (*authors.Authors, error) {
	args := m.Called(ctx, author)
	return args.Get(0).(*authors.Authors), args.Error(1)
}

func (m *AuthorUsecaseMock) UpdateAuthor(ctx context.Context, author *authors.UpdateAuthorRequest) (*authors.Authors, error) {
	args := m.Called(ctx, author)
	return args.Get(0).(*authors.Authors), args.Error(1)
}

func (m *AuthorUsecaseMock) DeleteAuthor(ctx context.Context, authorId uint16) error {
	return m.Called(ctx).Error(1)
}
