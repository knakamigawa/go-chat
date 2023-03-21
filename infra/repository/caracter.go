package repository

import (
	"context"
	"go-chat-ai-server/domain/model"
	"go-chat-ai-server/domain/model/character"
	"go-chat-ai-server/infra/database/chat_db"
)

func ProvideDbCharacter(queries *chat_db.Queries) *DbCharacter {
	return &DbCharacter{queries: queries}
}

type DbCharacter struct {
	queries *chat_db.Queries
}

func (r DbCharacter) All(ctx context.Context) (model.Characters, error) {
	characters, err := r.queries.ListCharacter(ctx)
	if err != nil {
		return model.Characters{}, err
	}
	models := make([]model.Character, len(characters))
	for i, v := range characters {
		models[i] = model.MakeCharacter(character.Name(v.Name), character.Bio(v.Bio))
	}

	return models, nil
}

func (r DbCharacter) Of(ctx context.Context, name string) (model.Character, error) {
	findCharacter, err := r.queries.FindCharacter(ctx, name)
	if err != nil {
		return model.Character{}, err
	}
	return model.MakeCharacter(character.Name(findCharacter.Name), character.Bio(findCharacter.Bio)), nil
}

func (r DbCharacter) Save(ctx context.Context, char model.Character) error {
	param := chat_db.CreateCharacterParams{
		Name: char.Name().String(),
		Bio:  char.Bio().String(),
	}
	_, err := r.queries.CreateCharacter(ctx, param)
	if err != nil {
		return err
	}

	return nil
}
