package repository

import (
	"context"

	"github.com/kasa021/watabe-lab-app/internal/domain"
	"gorm.io/gorm"
)

type AchievementRepository interface {
	FindAll(ctx context.Context) ([]domain.Achievement, error)
	FindByCode(ctx context.Context, code string) (*domain.Achievement, error)
	CreateUserAchievement(ctx context.Context, ua *domain.UserAchievement) error
	GetUserAchievements(ctx context.Context, userID uint) ([]domain.UserAchievement, error)
	HasUnlocked(ctx context.Context, userID uint, achievementID uint) (bool, error)
}

type achievementRepository struct {
	db *gorm.DB
}

func NewAchievementRepository(db *gorm.DB) AchievementRepository {
	return &achievementRepository{db: db}
}

func (r *achievementRepository) FindAll(ctx context.Context) ([]domain.Achievement, error) {
	var achievements []domain.Achievement
	if err := r.db.WithContext(ctx).Find(&achievements).Error; err != nil {
		return nil, err
	}
	return achievements, nil
}

func (r *achievementRepository) FindByCode(ctx context.Context, code string) (*domain.Achievement, error) {
	var achievement domain.Achievement
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&achievement).Error; err != nil {
		return nil, err
	}
	return &achievement, nil
}

func (r *achievementRepository) CreateUserAchievement(ctx context.Context, ua *domain.UserAchievement) error {
	return r.db.WithContext(ctx).Create(ua).Error
}

func (r *achievementRepository) GetUserAchievements(ctx context.Context, userID uint) ([]domain.UserAchievement, error) {
	var uas []domain.UserAchievement
	if err := r.db.WithContext(ctx).
		Preload("Achievement").
		Where("user_id = ?", userID).
		Find(&uas).Error; err != nil {
		return nil, err
	}
	return uas, nil
}

func (r *achievementRepository) HasUnlocked(ctx context.Context, userID uint, achievementID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&domain.UserAchievement{}).
		Where("user_id = ? AND achievement_id = ?", userID, achievementID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
