package service

import (
	"context"
	"time"

	"github.com/kasa021/watabe-lab-app/internal/domain"
	"github.com/kasa021/watabe-lab-app/internal/repository"
)

type AchievementService interface {
	GetAchievements(ctx context.Context) ([]domain.Achievement, error)
	GetUserAchievements(ctx context.Context, userID uint) ([]domain.UserAchievement, error)
	CheckAndUnlock(ctx context.Context, userID uint, triggerType string, value interface{}) ([]domain.Achievement, error)
}

type achievementService struct {
	repo     repository.AchievementRepository
	userRepo repository.UserRepository
	// logRepo  repository.AttendanceRepository // 複雑な条件判定に必要なら追加
}

func NewAchievementService(repo repository.AchievementRepository, userRepo repository.UserRepository) AchievementService {
	return &achievementService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *achievementService) GetAchievements(ctx context.Context) ([]domain.Achievement, error) {
	return s.repo.FindAll(ctx)
}

func (s *achievementService) GetUserAchievements(ctx context.Context, userID uint) ([]domain.UserAchievement, error) {
	return s.repo.GetUserAchievements(ctx, userID)
}

// CheckAndUnlock は指定されたトリガー（例: "check_in_count"）に基づいて実績を判定し、解除する
func (s *achievementService) CheckAndUnlock(ctx context.Context, userID uint, triggerType string, value interface{}) ([]domain.Achievement, error) {
	allAchievements, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var unlocked []domain.Achievement

	for _, ach := range allAchievements {
		// 既に解除済みかチェック
		hasUnlocked, err := s.repo.HasUnlocked(ctx, userID, ach.ID)
		if err != nil {
			continue // エラーでも他を続行
		}
		if hasUnlocked {
			continue
		}

		// 条件判定（簡易実装）
		// マスタデータの ConditionType と、引数の triggerType が一致する場合のみ判定
		if ach.ConditionType != triggerType {
			continue
		}

		shouldUnlock := false

		// ConditionValue は JSONB (map[string]interface{})
		targetVal, ok := ach.ConditionValue["target"].(float64) // JSONの数値はfloat64で来ることが多い
		if !ok {
			continue
		}
		intTarget := int(targetVal)

		// トリガーごとの判定ロジック
		switch triggerType {
		case "total_check_in": // 累計回数
			currentCount, ok := value.(int)
			if ok && currentCount >= intTarget {
				shouldUnlock = true
			}
		case "total_duration": // 累計時間（分）
			currentDuration, ok := value.(int)
			if ok && currentDuration >= intTarget {
				shouldUnlock = true
			}
			// case "streak": // ストリーク（別途DB参照が必要かも）
		}

		if shouldUnlock {
			ua := &domain.UserAchievement{
				UserID:        userID,
				AchievementID: ach.ID,
				AchievedAt:    time.Now(),
			}
			if err := s.repo.CreateUserAchievement(ctx, ua); err == nil {
				unlocked = append(unlocked, ach)
			}
		}
	}

	return unlocked, nil
}
