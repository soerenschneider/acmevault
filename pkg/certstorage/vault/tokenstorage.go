package vault

type TokenStorage interface {
	StoreToken(token string) error
	ReadToken() (string, error)
}
