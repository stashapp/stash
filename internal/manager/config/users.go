package config

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

const (
	Username = "username"
	Password = "password"
	Users    = "users"
	Roles    = "roles"
)

type StoredUser struct {
	Username     string            `json:"username" koanf:"username"`
	PasswordHash string            `json:"passwordhash" koanf:"passwordhash"`
	Roles        []models.RoleEnum `json:"roles" koanf:"roles"`
}

type UserStore struct {
	*Config

	cachedUsers map[string]StoredUser
}

func (s *Config) GetUsername() string {
	return s.getString(Username)
}

func (s *Config) GetPasswordHash() string {
	return s.getString(Password)
}

func (s *UserStore) legacyUser() *StoredUser {
	un := s.getString(Username)
	pwHash := s.getString(Password)

	if un != "" && pwHash != "" {
		return &StoredUser{
			Username:     un,
			PasswordHash: pwHash,
			Roles:        []models.RoleEnum{models.RoleEnumAdmin},
		}
	}

	return nil
}

func (s *UserStore) loadUsers() error {
	// done outside lock to avoid deadlock
	legacyUser := s.legacyUser()

	s.RLock()
	defer s.RUnlock()

	var ret []*StoredUser
	err := s.unmarshalKey(Users, &ret)
	if err != nil {
		return err
	}

	// add legacy username
	if legacyUser != nil {
		ret = append(ret, legacyUser)
	}

	s.cachedUsers = make(map[string]StoredUser)
	for _, u := range ret {
		s.cachedUsers[u.Username] = *u
	}

	return nil
}

func (s *UserStore) convertUser(su StoredUser) *models.User {
	return &models.User{
		Username: su.Username,
		Roles:    su.Roles,
	}
}

func (s *UserStore) getUser(username string) *StoredUser {
	u, ok := s.cachedUsers[username]
	if !ok {
		return nil
	}

	return &u
}

func (s *UserStore) GetUser(ctx context.Context, username string) (*models.User, error) {
	s.RLock()
	defer s.RUnlock()

	u := s.getUser(username)
	if u == nil {
		return nil, nil
	}

	return s.convertUser(*u), nil
}

func (s *UserStore) AllUsers(ctx context.Context) ([]*models.User, error) {
	var users []*models.User

	s.RLock()
	defer s.RUnlock()

	for _, su := range s.cachedUsers {
		users = append(users, s.convertUser(su))
	}

	return users, nil
}

func (s *UserStore) LoginRequired(ctx context.Context) bool {
	return len(s.cachedUsers) > 0
}

func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	return string(hash)
}

func (s *UserStore) ValidateCredentials(ctx context.Context, username string, password string) bool {
	u := s.getUser(username)
	if u == nil {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))

	return err == nil
}

func (s *UserStore) saveUsers() error {
	// convert to list
	users := make([]StoredUser, 0, len(s.cachedUsers))
	for _, u := range s.cachedUsers {
		users = append(users, u)
	}

	s.setInterfaceNoLock(Users, users)
	return s.writeNoLock()
}

func (s *UserStore) ChangeUserPassword(ctx context.Context, username string, newPassword string) error {
	s.Lock()
	defer s.Unlock()

	u := s.getUser(username)
	if u == nil {
		return fmt.Errorf("user not found")
	}

	newHash := hashPassword(newPassword)

	updatedUser := *u
	updatedUser.PasswordHash = newHash
	s.cachedUsers[username] = updatedUser

	return s.saveUsers()
}

func (s *UserStore) CreateUser(ctx context.Context, u models.User, password string) error {
	s.Lock()
	defer s.Unlock()

	existingUser := s.getUser(u.Username)
	if existingUser != nil {
		return fmt.Errorf("user already exists")
	}

	newUser := StoredUser{
		Username:     u.Username,
		PasswordHash: hashPassword(password),
		Roles:        u.Roles,
	}

	s.cachedUsers[u.Username] = newUser

	return s.saveUsers()
}

func (s *UserStore) ReplaceUser(ctx context.Context, username string, updated models.User) error {
	s.Lock()
	defer s.Unlock()

	existingUser := s.getUser(username)
	if existingUser == nil {
		return fmt.Errorf("user not found")
	}

	updatedUser := StoredUser{
		Username:     updated.Username,
		PasswordHash: existingUser.PasswordHash,
		Roles:        updated.Roles,
	}

	// if username changed, remove old entry
	if username != updated.Username {
		delete(s.cachedUsers, username)
	}

	s.cachedUsers[updated.Username] = updatedUser

	return s.saveUsers()
}

func (s *UserStore) DeleteUser(ctx context.Context, username string) error {
	s.Lock()
	defer s.Unlock()

	existingUser := s.getUser(username)
	if existingUser == nil {
		return fmt.Errorf("user not found")
	}

	delete(s.cachedUsers, username)

	return s.saveUsers()
}
