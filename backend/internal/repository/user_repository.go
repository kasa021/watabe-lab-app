package repository

import (
	"github.com/kasa021/watabe-lab-app/internal/domain"
	"gorm.io/gorm"
)

// UserRepository ユーザーリポジトリのインターフェース
type UserRepository interface {
	Create(user *domain.User) error
	FindByID(id uint) (*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
	FindAll() ([]domain.User, error)
	Update(user *domain.User) error
	Delete(id uint) error
}

// userRepository ユーザーリポジトリの実装
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository ユーザーリポジトリを作成
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create ユーザーを作成
func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

// FindByID IDでユーザーを検索
func (r *userRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername ユーザー名でユーザーを検索
func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindAll すべてのユーザーを取得
func (r *userRepository) FindAll() ([]domain.User, error) {
	var users []domain.User
	err := r.db.Find(&users).Error
	return users, err
}

// Update ユーザー情報を更新
func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

// Delete ユーザーを削除
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&domain.User{}, id).Error
}

