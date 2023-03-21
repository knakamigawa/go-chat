package repository

import (
	"context"
	"go-chat-ai-server/domain/model"
)

type Character interface {
	Of(ctx context.Context, name string) (model.Character, error)
	Save(ctx context.Context, char model.Character) error
	All(ctx context.Context) (model.Characters, error)
}
