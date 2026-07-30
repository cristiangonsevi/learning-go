package model

import "time"

type TaskModel struct {
	ID        int       `db:"id,omitempty"`
	Name      string    `db:"name"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at,omitempty"`
}
