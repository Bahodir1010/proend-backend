// models/template.go
package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// Template stores information about a master document template.
type Template struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	FileID    uuid.UUID `json:"file_id" db:"file_id"`
	Version   int       `json:"version" db:"version"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Templates is a slice of Template structs
type Templates []Template