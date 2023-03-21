package service

import (
	"context"
	"go-chat-ai-server/app/client"
	"go-chat-ai-server/app/repository"
	"go-chat-ai-server/domain/model"
	"go-chat-ai-server/domain/model/character"
)

func ProvideChat(client client.ChatAPIClient, characterRepo repository.Character) Chat {
	return Chat{client: client, characterRepo: characterRepo}
}

type Chat struct {
	client        client.ChatAPIClient
	characterRepo repository.Character
}

func (s Chat) Characters(ctx context.Context) ([]string, error) {
	characters, err := s.characterRepo.All(ctx)
	if err != nil {
		return []string{}, err
	}
	return characters.Names(), nil
}

func (s Chat) Talk(ctx context.Context, inputText string, name string) (string, error) {
	char, err := s.characterRepo.Of(ctx, name)
	if err != nil {
		return "", err
	}

	txt, err := s.client.Request(inputText, char)
	if err != nil {
		return "", err
	}

	return txt, nil
}

func (s Chat) CharacterCreate(ctx context.Context, name string, bio string) error {
	char := model.MakeCharacter(character.Name(name), character.Bio(bio))

	err := s.characterRepo.Save(ctx, char)
	if err != nil {
		return err
	}

	return nil
}