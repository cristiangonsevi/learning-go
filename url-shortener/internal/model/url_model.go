package model

import "time"

type URL struct {
	ID        string
	URL       string
	ShortURL  string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type CreateURLParams struct {
	URL       string
	ExpiresAt time.Time
}
