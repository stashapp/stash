package user

import (
	"context"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/session"
)

type mockConfig struct {
	public bool
	guest  bool
}

func (m *mockConfig) GetPublicAccess() bool     { return m.public }
func (m *mockConfig) GetGuestUserEnabled() bool { return m.guest }
func (m *mockConfig) IsNewSystem() bool         { return false }
func (m *mockConfig) SetGuestUserEnabled(enabled bool) {
	m.guest = enabled
}

type mockStore struct {
	users     map[string]*models.User
	passwords map[int]string
	nextID    int
}

func newMockStore() *mockStore {
	return &mockStore{users: make(map[string]*models.User), passwords: make(map[int]string), nextID: 1}
}

func (m *mockStore) All(ctx context.Context) ([]*models.User, error) {
	ret := make([]*models.User, 0, len(m.users))
	for _, u := range m.users {
		ret = append(ret, u)
	}
	return ret, nil
}

func (m *mockStore) Count(ctx context.Context) (int, error) { return len(m.users), nil }

func (m *mockStore) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	if u, ok := m.users[username]; ok {
		return u, nil
	}
	return nil, nil
}

func (m *mockStore) FindAdminUsers(ctx context.Context) ([]*models.User, error) {
	var ret []*models.User
	for _, u := range m.users {
		if u.Roles.HasRole(models.RoleEnumAdmin) {
			ret = append(ret, u)
		}
	}
	return ret, nil
}

func (m *mockStore) GetPasswordHash(ctx context.Context, id int) (string, error) {
	return m.passwords[id], nil
}

func (m *mockStore) Create(ctx context.Context, u *models.User, password string) error {
	u.ID = m.nextID
	m.nextID++
	m.users[u.Username] = u
	if password != "" {
		m.passwords[u.ID] = password
	}
	return nil
}

func (m *mockStore) Update(ctx context.Context, updated *models.User) error {
	if _, ok := m.users[updated.Username]; !ok {
		// try to find by id and replace key if renamed
		for key, u := range m.users {
			if u.ID == updated.ID {
				delete(m.users, key)
				m.users[updated.Username] = updated
				return nil
			}
		}
		return nil
	}
	m.users[updated.Username] = updated
	return nil
}

func (m *mockStore) SetUserPassword(ctx context.Context, id int, newPassword string) error {
	m.passwords[id] = newPassword
	return nil
}

func (m *mockStore) SetUserAPIKey(ctx context.Context, id int, newAPIKey string) error { return nil }
func (m *mockStore) Destroy(ctx context.Context, id int) error {
	for k, u := range m.users {
		if u.ID == id {
			delete(m.users, k)
			delete(m.passwords, id)
			return nil
		}
	}
	return nil
}

func (m *mockStore) SetLock(ctx context.Context, id int, locked bool) error {
	for _, u := range m.users {
		if u.ID == id {
			u.Locked = locked
			return nil
		}
	}
	return nil
}

func TestDeleteUser(t *testing.T) {
	store := newMockStore()
	cfg := &mockConfig{public: false}
	svc := &Service{Store: store, Config: cfg}

	ctx := session.SetCurrentUser(context.Background(), models.User{Username: "admin"})

	// create admin user
	adminUser := &models.User{
		Username: "admin",
		Roles:    models.Roles{models.RoleEnumAdmin},
	}
	if err := svc.CreateUser(ctx, *adminUser, "password"); err != nil {
		t.Fatalf("failed to create admin user: %v", err)
	}

	// delete only admin user - should fail since it is the last admin user
	if err := svc.DeleteUser(ctx, "admin"); err == nil {
		t.Fatal("expected error when deleting last admin user, got nil")
	}

	// create non-admin user
	nonAdminUser := &models.User{
		Username: "user1",
		Roles:    models.Roles{models.RoleEnumRead},
	}
	if err := svc.CreateUser(ctx, *nonAdminUser, "password"); err != nil {
		t.Fatalf("failed to create non-admin user: %v", err)
	}

	// delete admin user - should still fail
	if err := svc.DeleteUser(ctx, "admin"); err == nil {
		t.Fatal("expected error when deleting admin user, got nil")
	}

	// delete non-admin user - should succeed
	if err := svc.DeleteUser(ctx, "user1"); err != nil {
		t.Fatalf("failed to delete non-admin user: %v", err)
	}
}
