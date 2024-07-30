package entity

type ShortenURL struct {
	ShortKey    string `json:"short_key"`
	OriginalURL string `json:"original_url"`
}
