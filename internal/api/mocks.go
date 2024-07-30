package api

func NewMockStorage() *MockStorage {
	return &MockStorage{}
}

type MockStorage struct {
	saveURL func(string) (string, error)
	getURL  func(string) (string, error)
}

func (m *MockStorage) SaveURL(originalURL string) (string, error) {
	return m.saveURL(originalURL)
}

func (m *MockStorage) GetURL(key string) (string, error) {
	return m.getURL(key)
}

func (m *MockStorage) Close() error {
	return nil
}
