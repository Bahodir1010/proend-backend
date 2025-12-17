// models/document.go
package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// Document struct
type Document struct {
	ID         uuid.UUID `json:"id" db:"id"`
	TemplateID uuid.UUID `json:"template_id" db:"template_id"`
	Filename   string    `json:"filename" db:"filename"`
	Status     string    `json:"status" db:"status"`

	// --- ADDED FIELDS ---
	FIO       string `json:"fio,omitempty" db:"fio"`
	Lavozim   string `json:"lavozim,omitempty" db:"lavozim"`
	Oylik     string `json:"oylik,omitempty" db:"oylik"`
	Stavka    string `json:"stavka,omitempty" db:"stavka"`
	Username  string `json:"username,omitempty" db:"username"`
	OrderType string `json:"order_type,omitempty" db:"order_type"`
	// --------------------

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Documents is a slice of Document structs
type Documents []Document

func (d *Document) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}