package model

import "time"

type URL struct {
	ID            string     `json:"id,omitempty" db:"id,omitempty"`
	URL           string     `json:"url" db:"url"`
	ShortURL      string     `json:"short_url" db:"short_url"`
	Title         string     `json:"title,omitempty" db:"title,omitempty"`
	Description   string     `json:"description,omitempty" db:"description,omitempty"`
	UserID        string     `json:"user_id,omitempty" db:"user_id,omitempty"`
	IsActive      bool       `json:"is_active,omitempty" db:"is_active,omitempty"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty" db:"expires_at,omitempty"`
	VisitCount    int        `json:"visit_count,omitempty" db:"visit_count,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty" db:"updated_at,omitempty"`
	LastVisitedAt *time.Time `json:"last_visited_at,omitempty" db:"last_visited_at,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty" db:"created_at,omitempty"`
}

type CreateURLParams struct {
	URL         string
	Title       string
	Description string
	UserID      string
	IsActive    bool
	ExpiresAt   *time.Time
}
