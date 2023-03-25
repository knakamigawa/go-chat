package repository

import (
	"context"
	"go-chat-ai-server/domain/model"
)

type UserMailPassword interface {
	Of(ctx context.Context, email model.Email, password model.Password) (model.User, error)
	Save(ctx context.Context, char model.UserEmailPassword) error
}
