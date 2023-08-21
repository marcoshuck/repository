package repository

import (
	"context"
)

// Repository contains a set of methods to interact with a persistence layer.
type Repository[E any] interface {
	// Create creates an entity in a persistence layer.
	Create(ctx context.Context, entity E) (E, error)
	// CreateBulk creates a set of entities in a persistence layer.
	CreateBulk(ctx context.Context, entities []E) ([]E, error)
	// Get returns an entity from a persistence layer identified by its ID. It returns an error if the entity doesn't exist.
	Get(ctx context.Context, id uint) (E, error)
	// Find returns a set of entities from a persistence layer identified by their IDs. It returns an empty slice if no
	// records were found.
	Find(ctx context.Context, ids []uint) ([]E, error)
	// Update updates with the values of entity the entity identified by id.
	Update(ctx context.Context, id uint, entity E) (E, error)
	// UpdateBulk updates with the values of entity all the elements identified by the slice of ids.
	UpdateBulk(ctx context.Context, ids []uint, entity E) ([]E, error)
	// Remove removes the given id from a persistence layer.
	Remove(ctx context.Context, id uint) (E, error)
	// RemoveBulk removes a set of elements from a persistence layer.
	RemoveBulk(ctx context.Context, ids []uint) ([]E, error)
}
