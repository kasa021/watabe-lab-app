package repository

import (
	"context"

	"github.com/kasa021/watabe-lab-app/internal/domain"
	"gorm.io/gorm"
)

type AttendanceRepository interface {
	Create(ctx context.Context, log *domain.CheckInLog) error
	Update(ctx context.Context, log *domain.CheckInLog) error
	GetActiveCheckIn(ctx context.Context, userID uint) (*domain.CheckInLog, error)
	GetAllActiveCheckIns(ctx context.Context) ([]domain.CheckInLog, error)
}

type attendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendanceRepository{db: db}
}

func (r *attendanceRepository) Create(ctx context.Context, log *domain.CheckInLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *attendanceRepository) Update(ctx context.Context, log *domain.CheckInLog) error {
	return r.db.WithContext(ctx).Save(log).Error
}

func (r *attendanceRepository) GetActiveCheckIn(ctx context.Context, userID uint) (*domain.CheckInLog, error) {
	var log domain.CheckInLog
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ? AND check_out_at IS NULL", userID).
		First(&log).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *attendanceRepository) GetAllActiveCheckIns(ctx context.Context) ([]domain.CheckInLog, error) {
	var logs []domain.CheckInLog
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("check_out_at IS NULL").
		Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}
