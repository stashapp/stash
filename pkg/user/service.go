package user

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/session"
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
	ErrLockSelf                 = errors.New("cannot lock/unlock yourself")
	ErrCannotModifyGuestUser    = errors.New("guest user cannot be modified")
	ErrUsersExist               = errors.New("users already exist in the system")
)

const (
	Argon2Time    = 5
	Argon2Memory  = 7 * 1024
	Argon2Threads = 5
	Argon2KeyLen  = 32

	SaltLength = 16
)

var Argon2Params = &argon2id.Params{
	Memory:      Argon2Memory,
	Iterations:  Argon2Time,
	Parallelism: Argon2Threads,
	SaltLength:  SaltLength,
	KeyLength:   Argon2KeyLen,
}

const AdminUsername = "admin"
const GuestUsername = "guest"

type UserSource interface {
	All(ctx context.Context) ([]*models.User, error)
	Count(ctx context.Context) (int, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	FindAdminUsers(ctx context.Context) ([]*models.User, error)

	GetPasswordHash(ctx context.Context, id int) (string, error)

	Create(ctx context.Context, u *models.User, password string) error
	Update(ctx context.Context, updated *models.User) error
	SetUserPassword(ctx context.Context, id int, newPassword string) error
	SetUserAPIKey(ctx context.Context, id int, newAPIKey string) error
	Destroy(ctx context.Context, id int) error
	SetLock(ctx context.Context, id int, locked bool) error
}

type UserServiceConfig interface {
	IsNewSystem() bool
	GetPublicAccess() bool

	GetGuestUserEnabled() bool
	SetGuestUserEnabled(bool)
}

type Service struct {
	Store  UserSource
	Config UserServiceConfig

	startedAt time.Time

	singleUserModeMutex sync.RWMutex
	singleUserMode      bool

	guestUserEnabled bool
}

func (s *Service) IsSingleUserMode() bool {
	s.singleUserModeMutex.RLock()
	defer s.singleUserModeMutex.RUnlock()

	return s.singleUserMode
}

func (s *Service) LoginRequired(ctx context.Context) (bool, error) {
	s.singleUserModeMutex.RLock()
	defer s.singleUserModeMutex.RUnlock()

	return !s.singleUserMode, nil
}

func (s *Service) Init(ctx context.Context) error {
	s.startedAt = time.Now()

	// determine if we are in single user mode
	// don't do this for publically accessible instances
	if !s.Config.GetPublicAccess() {
		s.singleUserModeMutex.Lock()
		defer s.singleUserModeMutex.Unlock()

		count, err := s.Store.Count(ctx)
		if err != nil {
			return fmt.Errorf("error counting users during initialization: %w", err)
		}

		if count == 1 {
			// determine if we are in single user mode based on whether there is exactly one user and that user has no password set
			u, err := s.Store.All(ctx)
			if err != nil {
				return fmt.Errorf("error getting users during initialization: %w", err)
			}

			if len(u) != 1 {
				return fmt.Errorf("expected exactly one user, got %d", len(u))
			}

			pwHash, err := s.Store.GetPasswordHash(ctx, u[0].ID)
			if err != nil {
				return fmt.Errorf("error getting password hash during initialization: %w", err)
			}

			s.singleUserMode = pwHash == ""
			if s.singleUserMode {
				logger.Infof("Single user mode enabled since there is exactly one user with no password set")
			}
		}
	}

	s.guestUserEnabled = s.Config.GetGuestUserEnabled()
	if s.guestUserEnabled {
		if s.singleUserMode {
			logger.Warnf("Guest user cannot be enabled in single user mode, ignoring guest user enabled setting")
			s.guestUserEnabled = false
		} else {
			logger.Info("Guest user enabled")
		}

	}

	return nil
}

func (s *Service) SetGuestUserEnabled(enabled bool) error {
	if enabled && s.singleUserMode {
		return fmt.Errorf("cannot enable guest user in single user mode")
	}
	s.guestUserEnabled = enabled
	s.Config.SetGuestUserEnabled(enabled)
	return nil
}

// GetGuestUser returns the guest user if it exists, or nil if it does not exist.
// The guest user is a special user that is used for unauthenticated access.
func (s *Service) GetGuestUser(ctx context.Context) *models.User {
	if !s.guestUserEnabled {
		return nil
	}

	return &models.User{
		Username: GuestUsername,
		Roles:    []models.RoleEnum{models.RoleEnumRead},
	}
}

// GetSingleUser returns the single user if it exists.
// It will return nil if there is no single user or if there are multiple users
// (since single user can only be used if it is the only user).
func (s *Service) GetSingleUser(ctx context.Context) (*models.User, error) {
	s.singleUserModeMutex.RLock()
	defer s.singleUserModeMutex.RUnlock()

	if !s.singleUserMode {
		return nil, nil
	}

	count, err := s.Store.Count(ctx)
	if err != nil {
		logger.Errorf("error counting users: %v", err)
		return nil, ErrInternalError
	}

	if count != 1 {
		return nil, ErrUsersExist
	}

	all, err := s.AllUsers(ctx)
	if err != nil {
		logger.Errorf("error getting all users: %v", err)
		return nil, ErrInternalError
	}
	// sanity check
	if len(all) != 1 {
		panic(fmt.Sprintf("expected exactly one user, got %d", len(all)))
	}
	return all[0], nil
}

func (s *Service) GetUser(ctx context.Context, username string) (*models.User, error) {
	return s.Store.FindByUsername(ctx, username)
}

func (s *Service) AllUsers(ctx context.Context) ([]*models.User, error) {
	return s.Store.All(ctx)
}

func userIsLocked(u *models.User) bool {
	return u.Locked || len(u.Roles) == 0
}

func (s *Service) LockUser(ctx context.Context, username string) error {
	// ensure caller not locking themselves
	cur := session.GetCurrentUser(ctx)
	if cur != nil && cur.Username == username {
		return ErrLockSelf
	}

	existingUser, err := s.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("error getting existing user: %w", err)
	}

	if existingUser == nil {
		return ErrUserNotExist
	}

	if err := s.Store.SetLock(ctx, existingUser.ID, true); err != nil {
		return fmt.Errorf("error locking user: %w", err)
	}

	logger.Infof("[user] locked %q by %q", username, cur.Username)

	return nil
}

func (s *Service) UnlockUser(ctx context.Context, username string) error {
	// ensure caller not unlocking themselves (same rule)
	cur := session.GetCurrentUser(ctx)
	if cur != nil && cur.Username == username {
		return ErrLockSelf
	}

	existingUser, err := s.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("error getting existing user: %w", err)
	}

	if existingUser == nil {
		return ErrUserNotExist
	}

	if err := s.Store.SetLock(ctx, existingUser.ID, false); err != nil {
		return fmt.Errorf("error unlocking user: %w", err)
	}

	logger.Infof("[user] unlocked %q by %q", username, cur.Username)

	return nil
}

func checkHash(password, hash string) (bool, error) {
	ret, _, err := argon2id.CheckHash(password, hash)
	return ret, err
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

	passwordHash, err := s.Store.GetPasswordHash(ctx, u.ID)
	if err != nil {
		logger.Errorf("error getting password hash for user %s: %v", username, err)
		return ErrInternalError
	}

	match, err := checkHash(password, passwordHash)
	if err != nil {
		logger.Errorf("error checking password hash for user %s: %v", username, err)
		return ErrInternalError
	}

	if !match {
		logger.Infof("[login attempt] invalid credentials for user %s", username)
		return ErrAccessDenied
	}
	return nil
}

// AuthenticateSession authenticates a user by their username and login time and returns the user object if successful.
// This is used for session-based authentication.
// It will return an error if the user does not exist, if the user is locked or if the user has been updated since the login time.
func (s *Service) AuthenticateSession(ctx context.Context, username string, loginTime time.Time) (*models.User, error) {
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

	// check if the user is logging in after a restart
	if s.startedAt.After(loginTime) {
		logger.Infof("[authentication] user %s logged in before service start time", username)
		return nil, ErrAccessDenied
	}

	// check if the user has been updated since the login time
	if u.UpdatedAt.After(loginTime) {
		logger.Infof("[authentication] user %s has been updated since login", username)
		return nil, ErrAccessDenied
	}

	return u, nil
}

func checkApiKey(apiKey, hash string) (bool, error) {
	hasher := sha256Hasher{}
	return hasher.CompareHash(apiKey, hash)
}

func (s *Service) AuthenticateByAPIKey(ctx context.Context, apiKey string) (*models.User, error) {
	username, err := GetUserIDFromAPIKey(apiKey)
	if errors.Is(err, ErrInvalidToken) {
		logger.Infof("[apikey authentication] invalid api key: %s", apiKey)
		return nil, ErrAccessDenied
	}

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

	match, err := checkApiKey(apiKey, user.ApiKeyHash)
	if err != nil {
		logger.Errorf("error checking api key hash for user %s: %v", username, err)
		return nil, ErrInternalError
	}

	// ensure apikey matches
	if !match {
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

func hashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, Argon2Params)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func (s *Service) CreateUser(ctx context.Context, u models.User, password string) error {
	s.singleUserModeMutex.Lock()
	defer s.singleUserModeMutex.Unlock()

	if s.singleUserMode {
		return errors.New("password must be set on initial user before other users can be created")
	}

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

	count, err := s.Store.Count(ctx)
	if err != nil {
		return fmt.Errorf("error getting existing users: %w", err)
	}

	hashedPassword := ""
	setSingleUserMode := false

	// special case for no password set
	if password == "" {
		// must be the first user created
		if count != 0 {
			return errors.New("password cannot be empty for non-initial users")
		}

		logger.Warnf("Creating initial user %q with no password set", u.Username)
		setSingleUserMode = true
	} else {
		// validate password
		if err := s.validatePassword(password); err != nil {
			return err
		}

		// hash the password and store it
		hashedPassword, err = hashPassword(password)
		if err != nil {
			return fmt.Errorf("error hashing password: %w", err)
		}
	}

	// if this is the first user, make them an admin
	if count == 0 && !u.Roles.HasRole(models.RoleEnumAdmin) {
		return errors.New("the first user must be an admin")
	}

	u.CreatedAt = time.Now()
	u.UpdatedAt = u.CreatedAt

	// create user in store
	if err := s.Store.Create(ctx, &u, hashedPassword); err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	cur := session.GetCurrentUser(ctx)
	logger.Infof("[user] created %q by %q", u.Username, cur.Username)
	if setSingleUserMode {
		s.singleUserMode = true
	}

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
		users, err := s.Store.FindAdminUsers(ctx)
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

	updated.ID = existingUser.ID

	// update user in store
	if err := s.Store.Update(ctx, &updated); err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	if username != updated.Username {
		logger.Infof("[user] updated name %q -> %q", username, updated.Username)
	}

	if !slices.Equal(existingRoles, updated.Roles) {
		logger.Infof("[user] updated roles for user %q", updated.Username)
	}

	cur := session.GetCurrentUser(ctx)
	logger.Infof("[user] updated %q by %q", updated.Username, cur.Username)

	return nil
}

func (s *Service) ChangePassword(ctx context.Context, username, currentPassword, newPassword string) error {
	if username == GuestUsername {
		return ErrCannotModifyGuestUser
	}

	// if we're in single user mode, allow changing password without validating current password since there is no other user
	if s.singleUserMode {
		return s.ChangeUserPassword(ctx, username, newPassword)
	}

	// validate current credentials
	if err := s.ValidateCredentials(ctx, username, currentPassword); err != nil {
		logger.Infof("[user] failed password change attempt for %q: incorrect current password", username)
		return ErrCurrentPasswordIncorrect
	}

	return s.ChangeUserPassword(ctx, username, newPassword)
}

func (s *Service) ChangeUserPassword(ctx context.Context, username, newPassword string) error {
	if username == GuestUsername {
		return ErrCannotModifyGuestUser
	}

	// check if user exists
	existingUser, err := s.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("error getting existing user: %w", err)
	}

	if existingUser == nil {
		return ErrUserNotExist
	}

	// if new password is empty, we are unsetting the password and enabling single user mode
	if newPassword == "" {
		s.singleUserModeMutex.Lock()
		defer s.singleUserModeMutex.Unlock()

		if s.Config.GetPublicAccess() {
			return errors.New("cannot unset password in public access mode")
		}

		count, err := s.Store.Count(ctx)
		if err != nil {
			return fmt.Errorf("error counting users: %w", err)
		}

		if count != 1 {
			return errors.New("can only unset password if there is exactly one user")
		}

		if err := s.Store.SetUserPassword(ctx, existingUser.ID, ""); err != nil {
			return fmt.Errorf("error unsetting user password: %w", err)
		}

		logger.Warnf("Unsetting password for user %q, enabling single user mode", username)
		s.singleUserMode = true
		return nil
	}

	// validate new password
	if err := s.validatePassword(newPassword); err != nil {
		return err
	}

	// hash the password and store it
	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	// change password in store
	if err := s.Store.SetUserPassword(ctx, existingUser.ID, hashedPassword); err != nil {
		return fmt.Errorf("error changing user password: %w", err)
	}

	cur := session.GetCurrentUser(ctx)
	logger.Infof("[user] changed password for %q by %q", username, cur.Username)

	s.singleUserModeMutex.Lock()
	defer s.singleUserModeMutex.Unlock()

	if s.singleUserMode {
		// turn off single user mode since password is now set and default user can no longer be used
		s.singleUserMode = false

		logger.Infof("Single user mode disabled since password has been set for user %q", username)
	}

	return nil
}

// GenerateRandomPassword generates a random password of a specified length.
func GenerateRandomPassword(length uint32) (string, error) {
	return generateRandomString(length)
}

// ResetUserPassword resets the password for the specified user and returns the new password.
// Used for emergency password resets from the command line when the user has been locked out.
func (s *Service) ResetUserPassword(ctx context.Context, username string) (string, error) {
	if s.singleUserMode {
		return "", errors.New("cannot reset user password in single user mode")
	}

	if username == GuestUsername {
		return "", ErrCannotModifyGuestUser
	}

	// check if user exists
	existingUser, err := s.GetUser(ctx, username)
	if err != nil {
		return "", fmt.Errorf("error getting existing user: %w", err)
	}

	if existingUser == nil {
		return "", ErrUserNotExist
	}

	const passwordLen = 16
	newPassword, err := generateRandomString(passwordLen)
	if err != nil {
		return "", fmt.Errorf("error generating new password: %w", err)
	}

	// hash the password and store it
	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}

	if err := s.Store.SetUserPassword(ctx, existingUser.ID, hashedPassword); err != nil {
		return "", fmt.Errorf("error changing user password: %w", err)
	}

	logger.Warnf("Password for user %q has been reset", username)
	return newPassword, nil
}

// func (s *Service) CreateGuestUser(ctx context.Context) error {
// 	if s.singleUserMode {
// 		return errors.New("cannot create guest user in single user mode")
// 	}

// 	u := models.User{
// 		Username:  "guest",
// 		Roles:     []models.RoleEnum{models.RoleEnumRead},
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 	}

// 	// create user in store with empty password
// 	if err := s.Store.Create(ctx, &u, ""); err != nil {
// 		return fmt.Errorf("error creating guest user: %w", err)
// 	}

// 	logger.Infof("[user] created guest user")

// 	return nil
// }

func hashAPIKey(apiKey string) (string, error) {
	// use faster SHA-256 hashing for API keys
	// https://cybersierra.co/blog/bcrypt-performance-issues-api/

	hasher := sha256Hasher{}
	return hasher.GenerateHash(apiKey)
}

func (s *Service) GenerateAPIKey(ctx context.Context, username string) (string, error) {
	if username == GuestUsername {
		return "", ErrCannotModifyGuestUser
	}

	if s.singleUserMode {
		// don't allow generating API key in single user mode since default user is meant to be used without password or API key
		return "", errors.New("cannot generate API key in single user mode")
	}

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

	// hash the api key and store it
	hashedAPIKey, err := hashAPIKey(newAPIKey)
	if err != nil {
		return "", fmt.Errorf("error hashing api key: %w", err)
	}

	if err := s.Store.SetUserAPIKey(ctx, existingUser.ID, hashedAPIKey); err != nil {
		return "", fmt.Errorf("error updating user with new api key: %w", err)
	}

	cur := session.GetCurrentUser(ctx)
	logger.Infof("[user] generated new API key for %q by %q", username, cur.Username)

	return newAPIKey, nil
}

func (s *Service) ClearAPIKey(ctx context.Context, username string) error {
	if username == GuestUsername {
		return ErrCannotModifyGuestUser
	}

	if s.singleUserMode {
		// don't allow clearing API key in single user mode since default user is meant to be used without password or API key
		return errors.New("cannot clear API key in single user mode")
	}

	// check if user exists
	existingUser, err := s.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("error getting existing user: %w", err)
	}

	if existingUser == nil {
		return ErrUserNotExist
	}

	// clear api key
	if err := s.Store.SetUserAPIKey(ctx, existingUser.ID, ""); err != nil {
		return fmt.Errorf("error clearing user api key: %w", err)
	}

	cur := session.GetCurrentUser(ctx)
	logger.Infof("[user] cleared API key for %q by %q", username, cur.Username)

	return nil
}

func (s *Service) DeleteUser(ctx context.Context, username string) error {
	if s.singleUserMode {
		return errors.New("cannot delete user in single user mode")
	}

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
		users, err := s.Store.FindAdminUsers(ctx)
		if err != nil {
			return fmt.Errorf("error getting admin users: %w", err)
		}

		hasAdmin := false
		for _, u := range users {
			if u.Username != username && u.Roles.HasRole(models.RoleEnumAdmin) {
				hasAdmin = true
				break
			}
		}

		if !hasAdmin {
			return ErrDeleteLastAdminUser
		}
	}

	// delete user from store
	if err := s.Store.Destroy(ctx, existingUser.ID); err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	cur := session.GetCurrentUser(ctx)
	logger.Infof("[user] deleted %q by %q", username, cur.Username)

	return nil
}
