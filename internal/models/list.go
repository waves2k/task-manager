package models

import (
	"time"

	"github.com/google/uuid"
)

type List struct {
	Id          uuid.UUID `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
