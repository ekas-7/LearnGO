package model

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	Price       float64   `json:"price" validate:"required,gt=0"`
	Stock       int       `json:"stock" validate:"required,gte=0"`
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
	Category    *Category `json:"category,omitempty"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductCreateRequest struct {
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	Price       float64   `json:"price" validate:"required,gt=0"`
	Stock       int       `json:"stock" validate:"required,gte=0"`
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
	ImageURL    string    `json:"image_url"`
}

type ProductUpdateRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"omitempty,gt=0"`
	Stock       int     `json:"stock" validate:"omitempty,gte=0"`
	ImageURL    string  `json:"image_url"`
}

type ProductQueryParams struct {
	Page       int
	PageSize   int
	CategoryID uuid.UUID
	MinPrice   float64
	MaxPrice   float64
	Search     string
}
