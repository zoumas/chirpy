package database

type RevokedTokensRepository interface {
	Revoke(token string) error
	IsRevoked(token string) (bool, error)
}
