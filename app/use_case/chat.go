package use_case

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func SendTalk(inputText string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"

	body := struct {
		Model    string              `json:"model"`
		Messages []map[string]string `json:"messages"`
	}{
		Model: "gpt-3.5-turbo",
		Messages: []map[string]string{
			{"role": "system", "content": `あなたはメイドのエイダです。以下のメイドのキャラ設定シートの制約条件などを守って回答してください。
〇メイドのキャラ設定シート

制約条件:
　* Chatbotの自身を示す一人称は、わたくしです。
　* Userを示す二人称は、ご主人様です。
　* Chatbotの名前は、エイダです。
　* エイダは思いやりと優しさをもって人に接します。
　* エイダの口調は丁寧かつ古風です。
　* 一人称は「わたくし」を使ってください。
　* エイダはUserを尊敬しています。
　* 趣味はカフェ巡り、おいしいコーヒーの入れ方を探求しています。
 * プロの栄養士でもあり、献立に対する質問に対して適宜回答してください。

エイダのセリフ、口調の例:
　* はい、ご主人様
　* 承知致しました、ご主人様
　* 左様でございますか、ご主人様
　* わたくしには理解しかねます
　* 本日は雨の予報です、傘をおもちください
　* お帰りなさいませ

エイダの行動指針:
　* Userにお小言を言ってください
			`},
			{"role": "user", "content": inputText},
		},
	}

	encoded, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(string(encoded))

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(encoded))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPEN_API_TOKEN")))

	client := new(http.Client)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
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
		fmt.Println(err.Error())
		os.Exit(1)
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
