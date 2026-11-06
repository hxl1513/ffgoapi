package model

import (
	"math"
	"time"

	"gorm.io/gorm"
)

// BaseModel contains common fields for all database models
type BaseModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// User is the example user model used in the examples package.
type User struct {
	BaseModel
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"uniqueIndex;not null"`
	Password string `json:"-" gorm:"not null"`
	Role     int    `json:"role" gorm:"default:2"`
}

// TableName overrides the default pluralised table name for User.
func (User) TableName() string {
	return "users"
}

// PaginatedResult holds paginated query results with metadata
type PaginatedResult[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	TotalPages int   `json:"total_pages"`
}

// Paginate creates a PaginatedResult from the given items and total count
func Paginate[T any](items []T, total int64, page, perPage int) PaginatedResult[T] {
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	if totalPages == 0 {
		totalPages = 1
	}

	// Calculate slice bounds
	offset := (page - 1) * perPage
	end := offset + perPage
	if end > len(items) {
		end = len(items)
	}
	if offset > len(items) {
		offset = len(items)
	}

	pagedItems := items[offset:end]

	return PaginatedResult[T]{
		Items:      pagedItems,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}
}
