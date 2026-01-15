package service

import (
	"context"
	"errors"
	"time"

	"github.com/kasa021/watabe-lab-app/internal/domain"
	"github.com/kasa021/watabe-lab-app/internal/repository"
	"github.com/kasa021/watabe-lab-app/internal/ws"
	"gorm.io/gorm"
)

var (
	ErrAlreadyCheckedIn = errors.New("already checked in")
	ErrNotCheckedIn     = errors.New("not checked in")
)

type AttendanceService interface {
	CheckIn(ctx context.Context, userID uint, req *CheckInRequest) error
	CheckOut(ctx context.Context, userID uint) error
	GetActiveUsers(ctx context.Context) ([]domain.CheckInLog, error)
}

type attendanceService struct {
	repo repository.AttendanceRepository
	hub  *ws.Hub
}

func NewAttendanceService(repo repository.AttendanceRepository, hub *ws.Hub) AttendanceService {
	return &attendanceService{repo: repo, hub: hub}
}

type CheckInRequest struct {
	CheckInMethod string   `json:"check_in_method"`
	WiFiSSID      string   `json:"wifi_ssid"`
	GPSLatitude   *float64 `json:"gps_latitude"`
	GPSLongitude  *float64 `json:"gps_longitude"`
}

func (s *attendanceService) CheckIn(ctx context.Context, userID uint, req *CheckInRequest) error {
	// 既にチェックイン中か確認
	activeLog, err := s.repo.GetActiveCheckIn(ctx, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if activeLog != nil {
		return ErrAlreadyCheckedIn
	}

	// 新規チェックインログ作成
	log := &domain.CheckInLog{
		UserID:        userID,
		CheckInAt:     time.Now(),
		CheckInMethod: req.CheckInMethod,
		WiFiSSID:      req.WiFiSSID,
		GPSLatitude:   req.GPSLatitude,
		GPSLongitude:  req.GPSLongitude,
	}

	if err := s.repo.Create(ctx, log); err != nil {
		return err
	}

	// ユーザー情報を含めてブロードキャストするために、再度取得（または手動で構築）
	// ここではシンプルに、作成したログにユーザー情報をセットするためにリロードするか、
	// あるいは GetActiveCheckIn を呼ぶ。
	activeLog, err = s.repo.GetActiveCheckIn(ctx, userID)
	if err == nil {
		s.hub.BroadcastMessage(map[string]interface{}{
			"type":    "check_in",
			"payload": activeLog,
		})
	}
	return nil
}

func (s *attendanceService) CheckOut(ctx context.Context, userID uint) error {
	// チェックイン中のログを取得
	log, err := s.repo.GetActiveCheckIn(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotCheckedIn
		}
		return err
	}

	now := time.Now()
	log.CheckOutAt = &now

	// 滞在時間（分）計算
	duration := int(now.Sub(log.CheckInAt).Minutes())
	log.DurationMinutes = &duration

	if err := s.repo.Update(ctx, log); err != nil {
		return err
	}

	s.hub.BroadcastMessage(map[string]interface{}{
		"type":    "check_out",
		"payload": log,
	})
	return nil
}

func (s *attendanceService) GetActiveUsers(ctx context.Context) ([]domain.CheckInLog, error) {
	return s.repo.GetAllActiveCheckIns(ctx)
}
