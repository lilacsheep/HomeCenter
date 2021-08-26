package models

import (
	"time"
)

type DefaultModel struct {
	Id         int
	UpdatedAt time.Time `orm:"updated_at"`
	CreatedAt time.Time `orm:"created_at"`
}
