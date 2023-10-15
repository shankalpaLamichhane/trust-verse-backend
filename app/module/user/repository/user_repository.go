package repository

import "time"

type User struct {
	ID        int       `json:"id,omitempty"`
	Username  string    `json:"userName,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty""`
}
