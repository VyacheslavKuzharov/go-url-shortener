package api

func NewMockStorage() *MockStorage {
	return &MockStorage{}
}

type MockStorage struct {
	saveURL func(string) (string, error)
	getURL  func(string) (string, bool)
}

func (m *MockStorage) SaveURL(originalURL string) (string, error) {
	return m.saveURL(originalURL)
}

func (m *MockStorage) GetURL(key string) (string, bool) {
	return m.getURL(key)
}
