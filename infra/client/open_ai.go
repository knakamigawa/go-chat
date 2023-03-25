package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-chat-ai-server/domain/model"
	"io"
	"net/http"
	"os"
)

func ProvideChatAPIClient() ChatAPIClient {
	return ChatAPIClient{}
}

type ChatAPIClient struct{}

func (c ChatAPIClient) Request(inputText string, character model.Character) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"

	var messages []map[string]string
	if character.Bio() != "" {
		messages = append(messages, map[string]string{"role": "system", "content": character.Bio().String()})
	}
	messages = append(messages, map[string]string{"role": "user", "content": inputText})

	body := struct {
		Model    string              `json:"model"`
		Messages []map[string]string `json:"messages"`
	}{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
	}

	encoded, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	fmt.Println(string(encoded))

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(encoded))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPEN_API_TOKEN")))

	client := new(http.Client)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(resp.Body)

	byteArray, _ := io.ReadAll(resp.Body)
	var responseBody Response
	err = json.Unmarshal(byteArray, &responseBody)
	if err != nil {
		return "", err
	}

	var msg string
	for _, v := range responseBody.Choices {
		fmt.Println(v.Index)
		fmt.Println(v.Message.Role)
		fmt.Printf("msg: %s\n", v.Message.Content)
		fmt.Printf("res: %s\n", v.FinishReason)
		msg = v.Message.Content
	}
	fmt.Printf("PromptTokens    : %d\n", responseBody.Usage.PromptTokens)
	fmt.Printf("CompletionTokens: %d\n", responseBody.Usage.CompletionTokens)
	fmt.Printf("TotalTokens     : %d\n", responseBody.Usage.TotalTokens)
	return msg, nil
}

type Response struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}
