package repository

import (
	"context"
	"time"

	"github.com/kasa021/watabe-lab-app/internal/domain"
	"gorm.io/gorm"
)

type AttendanceRepository interface {
	Create(ctx context.Context, log *domain.CheckInLog) error
	Update(ctx context.Context, log *domain.CheckInLog) error
	GetActiveCheckIn(ctx context.Context, userID uint) (*domain.CheckInLog, error)
	GetAllActiveCheckIns(ctx context.Context) ([]domain.CheckInLog, error)
	GetUserRanking(ctx context.Context, from, to time.Time) ([]domain.UserRanking, error)
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

func (r *attendanceRepository) GetUserRanking(ctx context.Context, from, to time.Time) ([]domain.UserRanking, error) {
	var results []domain.UserRanking
	// JOINしてUser情報も一度に取得
	if err := r.db.WithContext(ctx).
		Table("check_in_logs").
		Select("check_in_logs.user_id, SUM(check_in_logs.duration_minutes) as total_duration, users.display_name, users.username").
		Joins("JOIN users ON users.id = check_in_logs.user_id").
		Where("check_in_logs.check_in_at BETWEEN ? AND ?", from, to).
		Group("check_in_logs.user_id, users.display_name, users.username").
		Order("total_duration DESC").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
