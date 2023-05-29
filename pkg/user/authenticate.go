package user

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordMismatch = fmt.Errorf("password mismatch")
	ErrUserNotFound     = fmt.Errorf("user not found")
)

func (s *Service) ValidateCredentials(ctx context.Context, username string, password string) error {
	hc, err := s.Store.HasCredentials(ctx)
	if err != nil {
		return fmt.Errorf("error checking if credentials exist: %w", err)
	}

	if !hc {
		// don't need to authenticate if no credentials saved
		return nil
	}

	authPWHash, err := s.Store.GetPasswordHash(ctx, username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return err
		}
		return fmt.Errorf("error getting password hash: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(authPWHash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrPasswordMismatch
		}

		return fmt.Errorf("error comparing password hash: %w", err)
	}

	return nil
}
