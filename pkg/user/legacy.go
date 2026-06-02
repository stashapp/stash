package user

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

// for schema versions before the introduction of the user table, we load the user information from the config file, and create a default user if in single user mode

var (
	ErrUnsupported       = errors.New("unsupported operation in legacy mode")
	ErrInvalidLegacyHash = errors.New("invalid legacy hash format")
)

// EncodeLegacyHash encodes a bcrypt hash in the format for database storage.
func EncodeLegacyHash(hash string) (string, error) {
	encodedHash := base64.RawStdEncoding.EncodeToString([]byte(hash))
	hash = fmt.Sprintf("$bcrypt$%s", encodedHash)
	return hash, nil
}

// DecodeHash expects a hash created from this package, and parses it to return the params used to
// create it, as well as the salt and key (password hash).
// There may be $ characters in the key.
func decodeLegacyHash(input string) (hash string, err error) {
	vals := strings.Split(input, "$")
	if len(vals) != 3 {
		return "", ErrInvalidLegacyHash
	}

	if vals[1] != "bcrypt" {
		return "", ErrInvalidLegacyHash
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(vals[2])
	if err != nil {
		return "", fmt.Errorf("error decoding legacy hash: %w", err)
	}

	return string(decodedHash), nil
}

type LegacyService struct {
	Username     string
	PasswordHash string
}

func (s *LegacyService) Init(ctx context.Context) error {
	return nil
}

func (s *LegacyService) IsSingleUserMode() bool {
	return s.Username == "" || s.PasswordHash == ""
}

func (s *LegacyService) LoginRequired(ctx context.Context) (bool, error) {
	return s.Username != "" && s.PasswordHash != "", nil
}

func (s *LegacyService) SetGuestUserEnabled(enabled bool) error {
	return ErrUnsupported
}

// GetGuestUser returns the guest user if it exists, or nil if it does not exist.
// The guest user is a special user that is used for unauthenticated access.
func (s *LegacyService) GetGuestUser(ctx context.Context) *models.User {
	return nil
}

// GetSingleUser returns the single user if it exists.
// It will return nil if there is no single user or if there are multiple users
// (since single user can only be used if it is the only user).
func (s *LegacyService) GetSingleUser(ctx context.Context) (*models.User, error) {
	return nil, nil
}

func (s *LegacyService) GetUser(ctx context.Context, username string) (*models.User, error) {
	if s.Username == "" || s.PasswordHash == "" {
		return nil, nil
	}

	if username != s.Username {
		return nil, nil
	}

	return &models.User{
		Username: s.Username,
		Roles:    []models.RoleEnum{models.RoleEnumAdmin},
	}, nil
}

func (s *LegacyService) AllUsers(ctx context.Context) ([]*models.User, error) {
	if s.Username == "" || s.PasswordHash == "" {
		return nil, nil
	}

	return []*models.User{
		{
			Username: s.Username,
			Roles:    []models.RoleEnum{models.RoleEnumAdmin},
		},
	}, nil
}

func (s *LegacyService) LockUser(ctx context.Context, username string) error {
	return ErrUnsupported
}

func (s *LegacyService) UnlockUser(ctx context.Context, username string) error {
	return ErrUnsupported
}

func checkLegacyHash(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil, err
}

func (s *LegacyService) ValidateCredentials(ctx context.Context, username string, password string) error {
	// ensure user is not locked
	u, err := s.GetUser(ctx, username)
	if err != nil {
		logger.Errorf("error getting user for credential validation: %v", err)
		return ErrInternalError
	}

	if u == nil {
		logger.Infof("[login attempt] user %s not found during credential validation", username)
		return ErrAccessDenied
	}

	match, _ := checkLegacyHash(password, s.PasswordHash)
	if !match {
		logger.Infof("[login attempt] invalid credentials for user %s", username)
		return ErrAccessDenied
	}
	return nil
}

// AuthenticateSession authenticates a user by their username and login time and returns the user object if successful.
// This is used for session-based authentication.
// It will return an error if the user does not exist, if the user is locked or if the user has been updated since the login time.
func (s *LegacyService) AuthenticateSession(ctx context.Context, username string, loginTime time.Time) (*models.User, error) {
	u, err := s.GetUser(ctx, username)
	if err != nil {
		logger.Errorf("error getting user for authentication: %v", err)
		return nil, ErrInternalError
	}

	if u == nil {
		logger.Infof("[authentication] user %s not found", username)
		return nil, ErrAccessDenied
	}

	return u, nil
}

func (s *LegacyService) AuthenticateByAPIKey(ctx context.Context, apiKey string) (*models.User, error) {
	return nil, ErrUnsupported
}

func (s *LegacyService) CreateUser(ctx context.Context, u models.User, password string) error {
	return ErrUnsupported
}

func (s *LegacyService) UpdateUser(ctx context.Context, username string, updated models.User) error {
	return ErrUnsupported
}

func (s *LegacyService) ChangePassword(ctx context.Context, username, currentPassword, newPassword string) error {
	return ErrUnsupported
}

func (s *LegacyService) ChangeUserPassword(ctx context.Context, username, newPassword string) error {
	return ErrUnsupported
}

// ResetUserPassword resets the password for the specified user and returns the new password.
// Used for emergency password resets from the command line when the user has been locked out.
func (s *LegacyService) ResetUserPassword(ctx context.Context, username string) (string, error) {
	return "", ErrUnsupported
}

func (s *LegacyService) GenerateAPIKey(ctx context.Context, username string) (string, error) {
	return "", ErrUnsupported
}

func (s *LegacyService) ClearAPIKey(ctx context.Context, username string) error {
	return ErrUnsupported
}

func (s *LegacyService) DeleteUser(ctx context.Context, username string) error {
	return ErrUnsupported
}
