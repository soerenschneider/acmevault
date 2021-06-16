package vault

import "io/ioutil"

type FileTokenStorage struct {
	path string
}

func NewFileTokenStorage(path string) (*FileTokenStorage, error) {
	return &FileTokenStorage{
		path: path,
	}, nil
}

func (f *FileTokenStorage) StoreToken(token string) error {
	return ioutil.WriteFile(f.path, []byte(token), 0600)
}

func (f *FileTokenStorage) ReadToken() (string, error) {
	content, err := ioutil.ReadFile(f.path)
	if err != nil {
		return "", nil
	}

	return string(content), nil
}
