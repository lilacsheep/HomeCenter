package models

import (
	"time"
)

type DefaultModel struct {
	Id        int       `json:"id"`
	UpdatedAt time.Time `json:"updated_at" orm:"updated_at"`
	CreatedAt time.Time `json:"created_at" orm:"created_at"`
}
