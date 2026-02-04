package user

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

var (
	ErrUserNotExist             = errors.New("user not found")
	ErrEmptyUsername            = errors.New("empty username")
	ErrUsernameHasWhitespace    = errors.New("username has leading or trailing whitespace")
	ErrDeleteLastAdminUser      = errors.New("final admin user cannot be deleted")
	ErrRemoveLastAdminRole      = errors.New("final admin role cannot be removed")
	ErrInternalError            = errors.New("internal error")
	ErrAccessDenied             = errors.New("access denied")
	ErrCurrentPasswordIncorrect = errors.New("current password incorrect")
	ErrUserAlreadyExists        = errors.New("user with that username already exists")
)

type UserSource interface {
	AllUsers(ctx context.Context) ([]*models.User, error)
	GetUser(ctx context.Context, username string) (*models.User, error)
	ValidateCredentials(ctx context.Context, username string, password string) bool

	CreateUser(ctx context.Context, u models.User, password string) error
	ReplaceUser(ctx context.Context, username string, updated models.User) error
	ChangeUserPassword(ctx context.Context, username string, newPassword string) error
	ChangeUserAPIKey(ctx context.Context, username string, newAPIKey string) error
	DeleteUser(ctx context.Context, username string) error
}

type Service struct {
	Store UserSource
}

func (s *Service) LoginRequired(ctx context.Context) bool {
	u, _ := s.Store.AllUsers(ctx)
	return len(u) > 0
}

func (s *Service) GetUser(ctx context.Context, username string) (*models.User, error) {
	return s.Store.GetUser(ctx, username)
}

func (s *Service) AllUsers(ctx context.Context) ([]*models.User, error) {
	return s.Store.AllUsers(ctx)
}

func userIsLocked(u *models.User) bool {
	return len(u.Roles) == 0
}

func (s *Service) ValidateCredentials(ctx context.Context, username string, password string) error {
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

	if userIsLocked(u) {
		logger.Infof("[login attempt] user %s is locked", username)
		return ErrAccessDenied
	}

	if !s.Store.ValidateCredentials(ctx, username, password) {
		logger.Infof("[login attempt] invalid credentials for user %s", username)
		return ErrAccessDenied
	}
	return nil
}

// AuthenticateUserByID authenticates a user by their username and returns the user object if successful.
// This is used for session-based authentication.
// It will return an error if the user does not exist or if the user is locked.
func (s *Service) AuthenticateUserByID(ctx context.Context, username string) (*models.User, error) {
	u, err := s.GetUser(ctx, username)
	if err != nil {
		logger.Errorf("error getting user for authentication: %v", err)
		return nil, ErrInternalError
	}

	if u == nil {
		logger.Infof("[authentication] user %s not found", username)
		return nil, ErrAccessDenied
	}

	if userIsLocked(u) {
		logger.Infof("[authentication] user %s is locked", username)
		return nil, ErrAccessDenied
	}

	return u, nil
}

func (s *Service) AuthenticateByAPIKey(ctx context.Context, apiKey string) (*models.User, error) {
	username, err := GetUserIDFromAPIKey(apiKey)
	if err != nil {
		logger.Errorf("error getting user ID from api key: %v", err)
		return nil, ErrInternalError
	}

	user, err := s.GetUser(ctx, username)
	if err != nil {
		logger.Errorf("error getting user by username: %v", err)
		return nil, ErrInternalError
	}

	if user == nil {
		logger.Infof("[apikey authentication] user %s not found", username)
		return nil, ErrAccessDenied
	}

	if userIsLocked(user) {
		logger.Infof("[apikey authentication] user %s is locked", username)
		return nil, ErrAccessDenied
	}

	// ensure apikey matches
	if user.ApiKey != apiKey {
		logger.Infof("[apikey authentication] invalid api key for user %s", username)
		return nil, ErrAccessDenied
	}

	return user, nil
}

func (s *Service) validateUsername(username string) error {
	if username == "" {
		return ErrEmptyUsername
	}

	// username must not have leading or trailing whitespace
	trimmed := strings.TrimSpace(username)

	if trimmed != username {
		return ErrUsernameHasWhitespace
	}

	return nil
}

func (s *Service) validatePassword(password string) error {
	if password == "" {
		return errors.New("password cannot be empty")
	}

	// add more password validation as needed

	return nil
}

func (s *Service) CreateUser(ctx context.Context, u models.User, password string) error {
	// validate input
	// ensure username is valid
	if err := s.validateUsername(u.Username); err != nil {
		return err
	}

	// check if user exists
	existingUser, err := s.GetUser(ctx, u.Username)
	if err != nil {
		return fmt.Errorf("error checking existing users: %w", err)
	}

	if existingUser != nil {
		return ErrUserAlreadyExists
	}

	// validate password
	if err := s.validatePassword(password); err != nil {
		return err
	}

	// if this is the first user, make them an admin
	users, err := s.AllUsers(ctx)
	if err != nil {
		return fmt.Errorf("error getting existing users: %w", err)
	}

	if len(users) == 0 && !u.Roles.HasRole(models.RoleEnumAdmin) {
		return errors.New("the first user must be an admin")
	}

	// create user in store
	if err := s.Store.CreateUser(ctx, u, password); err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	logger.Infof("[user] created %q", u.Username)

	return nil
}

func (s *Service) UpdateUser(ctx context.Context, username string, updated models.User) error {
	// validate input
	// check if user exists
	existingUser, err := s.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("error getting existing user: %w", err)
	}

	if existingUser == nil {
		return ErrUserNotExist
	}

	existingRoles := existingUser.Roles

	// ensure username is valid
	if username != updated.Username {
		if err := s.validateUsername(updated.Username); err != nil {
			return err
		}

		// ensure new username doesn't already exist
		otherUser, err := s.GetUser(ctx, updated.Username)
		if err != nil {
			return fmt.Errorf("error checking existing user: %w", err)
		}

		if otherUser != nil {
			return ErrUserAlreadyExists
		}
	}

	// validate roles
	// don't allow removing admin from last admin user
	if existingRoles.HasRole(models.RoleEnumAdmin) && !updated.Roles.HasRole(models.RoleEnumAdmin) {
		users, err := s.AllUsers(ctx)
		if err != nil {
			return fmt.Errorf("error getting all users: %w", err)
		}

		hasAdmin := false
		for _, u := range users {
			if u.Username != existingUser.Username && u.Roles.HasRole(models.RoleEnumAdmin) {
				hasAdmin = true
				break
			}
		}

		if !hasAdmin {
			return ErrRemoveLastAdminRole
		}
	}

	// update user in store
	if err := s.Store.ReplaceUser(ctx, username, updated); err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	if username != updated.Username {
		logger.Infof("[user] updated name %q -> %q", username, updated.Username)
	}

	if !slices.Equal(existingRoles, updated.Roles) {
		logger.Infof("[user] updated roles for user %q", updated.Username)
	}

	return nil
}

func (s *Service) ChangePassword(ctx context.Context, username, currentPassword, newPassword string) error {
	// validate current credentials
	if err := s.ValidateCredentials(ctx, username, currentPassword); err != nil {
		logger.Infof("[user] failed password change attempt for %q: incorrect current password", username)
		return ErrCurrentPasswordIncorrect
	}

	return s.ChangeUserPassword(ctx, username, newPassword)
}

func (s *Service) ChangeUserPassword(ctx context.Context, username, newPassword string) error {
	// check if user exists
	existingUser, err := s.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("error getting existing user: %w", err)
	}

	if existingUser == nil {
		return ErrUserNotExist
	}

	// validate new password
	if err := s.validatePassword(newPassword); err != nil {
		return err
	}

	// change password in store
	if err := s.Store.ChangeUserPassword(ctx, username, newPassword); err != nil {
		return fmt.Errorf("error changing user password: %w", err)
	}

	logger.Infof("[user] changed password for %q", username)

	return nil
}

func (s *Service) GenerateAPIKey(ctx context.Context, username string) (string, error) {
	// check if user exists
	existingUser, err := s.GetUser(ctx, username)
	if err != nil {
		return "", fmt.Errorf("error getting existing user: %w", err)
	}

	if existingUser == nil {
		return "", ErrUserNotExist
	}

	// generate new api key
	newAPIKey, err := generateAPIKey(username)
	if err != nil {
		return "", fmt.Errorf("error generating api key: %w", err)
	}

	if err := s.Store.ChangeUserAPIKey(ctx, username, newAPIKey); err != nil {
		return "", fmt.Errorf("error updating user with new api key: %w", err)
	}

	logger.Infof("[user] generated new API key for %q", username)

	return newAPIKey, nil
}

func (s *Service) ClearAPIKey(ctx context.Context, username string) error {
	// check if user exists
	existingUser, err := s.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("error getting existing user: %w", err)
	}

	if existingUser == nil {
		return ErrUserNotExist
	}

	// clear api key
	if err := s.Store.ChangeUserAPIKey(ctx, username, ""); err != nil {
		return fmt.Errorf("error clearing user api key: %w", err)
	}

	logger.Infof("[user] cleared API key for %q", username)

	return nil
}

func (s *Service) DeleteUser(ctx context.Context, username string) error {
	// check if user exists
	existingUser, err := s.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("error getting existing user: %w", err)
	}

	if existingUser == nil {
		return ErrUserNotExist
	}

	// don't allow deleting last admin user unless it is the last user
	if existingUser.Roles.HasRole(models.RoleEnumAdmin) {
		users, err := s.AllUsers(ctx)
		if err != nil {
			return fmt.Errorf("error getting all users: %w", err)
		}

		hasAdmin := false
		for _, u := range users {
			if u.Username != username && u.Roles.HasRole(models.RoleEnumAdmin) {
				hasAdmin = true
				break
			}
		}

		// allow deleting last admin if it is the only user
		if !hasAdmin && len(users) > 1 {
			return ErrDeleteLastAdminUser
		}
	}

	// delete user from store
	if err := s.Store.DeleteUser(ctx, username); err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	logger.Infof("[user] deleted %q", username)

	return nil
}
