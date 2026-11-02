package session

type SessionConfig interface {
	GetSessionStoreKey() []byte
	GetMaxSessionAge() int
}
