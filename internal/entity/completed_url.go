package entity

// CompletedURL describes struct with fully completed short_url and original_url in format:
//
//	{
//	  "short_url": "http://...",
//	  "original_url": "http://..."
//	}
type CompletedURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
