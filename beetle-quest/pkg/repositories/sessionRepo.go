package repositories

type SessionRepo interface {
	CreateSession(token string) error
	RevokeToken(token string) error
	FindToken(token string) (string, error)
}
