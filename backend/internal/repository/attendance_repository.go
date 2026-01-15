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
	GetDailyAttendanceCounts(ctx context.Context, userID uint) ([]domain.DailyAttendance, error)
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
		Joins("JOIN users ON users.id = check_in_logs.user_id").
		Where("check_in_logs.check_out_at IS NULL AND users.is_presence_public = ?", true).
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

func (r *attendanceRepository) GetDailyAttendanceCounts(ctx context.Context, userID uint) ([]domain.DailyAttendance, error) {
	// 簡易的に check_in_logs から集計して返す
	// 本来は daily_attendances テーブルがあるならそこから取るべきだが、
	// 現状の check_in_logs から日別に集計するクエリを書く
	// PostgreSQLの DATE_TRUNC を使用

	type Result struct {
		Date         time.Time `json:"date"`
		ThinkingTime int       `json:"duration"` // minute
		Count        int       `json:"count"`
	}

	var results []Result

	// 日付ごとの滞在時間と回数を集計
	// duration_minutes が NULL の場合は 0 として扱う (COALESCE)
	err := r.db.WithContext(ctx).
		Table("check_in_logs").
		Select("DATE(check_in_at) as date, SUM(COALESCE(duration_minutes, 0)) as thinking_time, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("DATE(check_in_at)").
		Order("date DESC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// domain.DailyAttendance に変換
	var dailies []domain.DailyAttendance
	for _, res := range results {
		dailies = append(dailies, domain.DailyAttendance{
			AttendanceDate:       res.Date,
			TotalDurationMinutes: res.ThinkingTime,
			CheckInCount:         res.Count,
		})
	}

	return dailies, nil
}
