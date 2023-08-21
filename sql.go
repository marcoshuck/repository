package repository

import (
	"context"
	"gorm.io/gorm"
)

// SQL implements Repository for GORM.
type SQL[E any] struct {
	db *gorm.DB
}

func (r *SQL[E]) Create(ctx context.Context, entity E) (E, error) {
	if err := r.db.WithContext(ctx).Model(new(E)).Create(&entity).Error; err != nil {
		var zero E
		return zero, err
	}
	return entity, nil
}

func (r *SQL[E]) CreateBulk(ctx context.Context, entities []E) ([]E, error) {
	if err := r.db.WithContext(ctx).Model(new(E)).Create(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *SQL[E]) Get(ctx context.Context, id uint) (E, error) {
	var out E
	if err := r.db.WithContext(ctx).Model(new(E)).Where("id = ?", id).First(&out).Error; err != nil {
		var zero E
		return zero, err
	}
	return out, nil
}

func (r *SQL[E]) Find(ctx context.Context, ids []uint) ([]E, error) {
	var out []E
	if err := r.db.WithContext(ctx).Model(new(E)).Where("id IN (?)", ids).Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *SQL[E]) Update(ctx context.Context, id uint, entity E) (E, error) {
	q := r.db.WithContext(ctx).Model(new(E)).Where("id = ?", id).Updates(&entity)
	if err := q.Error; err != nil {
		var zero E
		return zero, err
	}
	if q.RowsAffected == 0 {
		var zero E
		return zero, gorm.ErrRecordNotFound
	}
	return entity, nil
}

func (r *SQL[E]) UpdateBulk(ctx context.Context, ids []uint, entity E) ([]E, error) {
	if err := r.db.WithContext(ctx).Model(new(E)).Where("id IN (?)", ids).Updates(&entity).Error; err != nil {
		return nil, err
	}

	result, err := r.Find(ctx, ids)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *SQL[E]) Remove(ctx context.Context, id uint) (E, error) {
	entity, err := r.Get(ctx, id)
	if err != nil {
		var zero E
		return zero, err
	}

	if err := r.db.WithContext(ctx).Model(new(E)).Where("id = ?", id).Delete(&entity).Error; err != nil {
		var zero E
		return zero, err
	}

	return entity, nil
}

func (r *SQL[E]) RemoveBulk(ctx context.Context, ids []uint) ([]E, error) {
	result, err := r.Find(ctx, ids)
	if err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Model(new(E)).Where("id IN (?)", ids).Delete(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func NewRepositorySQL[E any]() Repository[E] {
	return &SQL[E]{}
}
