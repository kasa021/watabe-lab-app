package repository

import (
	"context"

	"github.com/kasa021/watabe-lab-app/internal/domain"
	"gorm.io/gorm"
)

type SettingsRepository interface {
	GetByKey(ctx context.Context, key string) (*domain.Setting, error)
}

type settingsRepository struct {
	db *gorm.DB
}

func NewSettingsRepository(db *gorm.DB) SettingsRepository {
	return &settingsRepository{db: db}
}

func (r *settingsRepository) GetByKey(ctx context.Context, key string) (*domain.Setting, error) {
	var setting domain.Setting
	if err := r.db.WithContext(ctx).Where("key = ?", key).First(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}
