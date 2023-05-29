package session

import "context"

type ExternalAccessConfig interface {
	HasCredentials(ctx context.Context) (bool, error)
	GetDangerousAllowPublicWithoutAuth() bool
	GetSecurityTripwireAccessedFromPublicInternet() string
	IsNewSystem() bool
}

type SessionConfig interface {
	GetUsername() string
	GetAPIKey() string

	GetSessionStoreKey() []byte
	GetMaxSessionAge() int
}
