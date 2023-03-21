package client

import "go-chat-ai-server/domain/model"

type ChatAPIClient interface {
	Request(inputText string, character model.Character) (string, error)
}
