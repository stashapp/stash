package session

import "context"

type ExternalAccessConfig interface {
	GetDangerousAllowPublicWithoutAuth() bool
	GetSecurityTripwireAccessedFromPublicInternet() string
	IsNewSystem() bool
}

type CredentialStore interface {
	LoginRequired(ctx context.Context) bool
}

type SessionConfig interface {
	GetSessionStoreKey() []byte
	GetMaxSessionAge() int
}
