package user

import (
	"context"
)

type Store interface {
	GetUser(ctx context.Context, username string) (*User, error)
	HasCredentials(ctx context.Context) (bool, error)
	GetPasswordHash(ctx context.Context, username string) (string, error)
}

type Service struct {
	Store Store
}

func (s *Service) GetUser(ctx context.Context, username string) (*User, error) {
	return s.Store.GetUser(ctx, username)
}
