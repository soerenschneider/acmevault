package vault

type InMemoryTokenStorage struct {
	token string
}

func NewPopulatedInMemoryTokenStorage(vaultToken string) *InMemoryTokenStorage {
	return &InMemoryTokenStorage{
		token: vaultToken,
	}
}

func (f *InMemoryTokenStorage) StoreToken(token string) error {
	f.token = token
	return nil
}

func (f *InMemoryTokenStorage) ReadToken() (string, error) {
	return f.token, nil
}
