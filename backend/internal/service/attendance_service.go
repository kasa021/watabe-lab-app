package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"math"

	"github.com/kasa021/watabe-lab-app/internal/domain"
	"github.com/kasa021/watabe-lab-app/internal/repository"
	"github.com/kasa021/watabe-lab-app/internal/ws"
	"gorm.io/gorm"
)

var (
	ErrAlreadyCheckedIn     = errors.New("already checked in")
	ErrNotCheckedIn         = errors.New("not checked in")
	ErrRestrictionViolation = errors.New("check-in restriction violation")
)

type AttendanceService interface {
	CheckIn(ctx context.Context, userID uint, req *CheckInRequest) error
	CheckOut(ctx context.Context, userID uint) error
	GetActiveUsers(ctx context.Context) ([]domain.CheckInLog, error)
}

type attendanceService struct {
	repo         repository.AttendanceRepository
	settingsRepo repository.SettingsRepository // Added
	hub          *ws.Hub
	achService   AchievementService
}

func NewAttendanceService(repo repository.AttendanceRepository, settingsRepo repository.SettingsRepository, hub *ws.Hub, achService AchievementService) AttendanceService {
	return &attendanceService{
		repo:         repo,
		settingsRepo: settingsRepo, // Added
		hub:          hub,
		achService:   achService,
	}
}

type CheckInRequest struct {
	CheckInMethod string   `json:"check_in_method"`
	WiFiSSID      string   `json:"wifi_ssid"`
	GPSLatitude   *float64 `json:"gps_latitude"`
	GPSLongitude  *float64 `json:"gps_longitude"`
	ClientIP      string   `json:"-"` // Added for internal use
}

func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000 // Earth radius in meters
	phi1 := lat1 * math.Pi / 180
	phi2 := lat2 * math.Pi / 180
	deltaPhi := (lat2 - lat1) * math.Pi / 180
	deltaLambda := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) +
		math.Cos(phi1)*math.Cos(phi2)*
			math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

func (s *attendanceService) CheckIn(ctx context.Context, userID uint, req *CheckInRequest) error {
	// 1. IP Address Validation
	settingIP, err := s.settingsRepo.GetByKey(ctx, "allowed_ip_range")
	if err == nil {
		// Expecting JSON: {"ips": ["1.2.3.4", ...]}
		if ips, ok := settingIP.Value["ips"].([]interface{}); ok {
			isAllowed := false
			for _, v := range ips {
				if ip, ok := v.(string); ok && ip == req.ClientIP {
					isAllowed = true
					break
				}
				// TODO: CIDR support if needed
			}
			if !isAllowed {
				return fmt.Errorf("研究室のWifiに接続してください (Your IP: %s): %w", req.ClientIP, ErrRestrictionViolation)
			}
		}
	} else {
		// Just debugging: in production we might want to log this.
		// For now, if setting is missing, we allow. This explains why it worked before seeding.
	}

	// 2. GPS Validation
	settingGPS, err := s.settingsRepo.GetByKey(ctx, "gps_location")
	if err == nil && req.GPSLatitude != nil && req.GPSLongitude != nil {
		var labLocation struct {
			Latitude     float64 `json:"latitude"`
			Longitude    float64 `json:"longitude"`
			RadiusMeters float64 `json:"radius_meters"`
		}
		// Try to unmarshal into the struct. Note: Value is JSONB map[string]interface{}.
		// We need to marshal it back to json and unmarshal to struct, or use map access.
		// Since Setting.Value is JSONB (map[string]interface{}), access directly.
		if lat, ok := settingGPS.Value["latitude"].(float64); ok {
			labLocation.Latitude = lat
		}
		if lon, ok := settingGPS.Value["longitude"].(float64); ok {
			labLocation.Longitude = lon
		}
		if rad, ok := settingGPS.Value["radius_meters"].(float64); ok {
			labLocation.RadiusMeters = rad
		}

		dist := calculateDistance(*req.GPSLatitude, *req.GPSLongitude, labLocation.Latitude, labLocation.Longitude)
		if dist > labLocation.RadiusMeters {
			return fmt.Errorf("研究室の近くでチェックインしてください。位置情報が異なります。: %w", ErrRestrictionViolation)
		}
	} else if err == nil && (req.GPSLatitude == nil || req.GPSLongitude == nil) {
		// If GPS restriction is enabled but no GPS provided
		return fmt.Errorf("位置情報の取得に失敗しました: %w", ErrRestrictionViolation)
	}

	// 3. Existing CheckIn Logic
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

	// Broadcast check-out event
	s.hub.BroadcastMessage(map[string]interface{}{
		"type":    "check_out",
		"payload": log,
	})

	// 実績解除判定 (非同期)
	go func() {
		bgCtx := context.Background()
		// 累計回数判定 (とりあえず今回の1回をトリガーに全件チェック)
		// FIXME: 本来は現在の数値を渡すべきだが、Service内でCountsを取得する実装が必要。
		// ここでは簡易的に 0 を渡して、Service側で条件と一致するか見る (Service側も実装修正が必要)
		// 一旦、トリガータイプだけ合わせておく。
		s.achService.CheckAndUnlock(bgCtx, userID, "check_in_count", nil)
		s.achService.CheckAndUnlock(bgCtx, userID, "total_duration", nil)
	}()

	return nil
}

func (s *attendanceService) GetActiveUsers(ctx context.Context) ([]domain.CheckInLog, error) {
	return s.repo.GetAllActiveCheckIns(ctx)
}
