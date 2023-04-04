package domain

import "context"

type UserStorage interface {
	Create(ctx context.Context, user *User) error
	GetByUserID(ctx context.Context, userID string) (User, error)
	UpdateByUserID(ctx context.Context, userID string, values map[string]interface{}) (User, error)
	DeleteByUserID(ctx context.Context, userID string) error
}
