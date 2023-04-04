package repo

import (
	"context"

	"github.com/johnnyzhao/auth-server-api/internal/domain"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	if len(user.Nickname) == 0 {
		user.Nickname = user.UserID
	}
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepo) GetByUserID(ctx context.Context, userID string) (domain.User, error) {
	var user domain.User

	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&user).Error; err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *UserRepo) UpdateByUserID(ctx context.Context, userID string, values map[string]interface{}) (domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Model(&user).
		Where("user_id = ?", userID).
		Updates(values).Error; err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *UserRepo) DeleteByUserID(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&domain.User{}).Error
}

func (r *UserRepo) Migrate() error {
	return r.db.AutoMigrate(&domain.User{})
}
