package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	openai "github.com/sashabaranov/go-openai"
)

func OpenAI_35(issue, key string) string {
	msg := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: issue,
		},
	}
	return OpenAI_35_Context(msg, key)
	// result := make(chan string)
	// go func() {
	// 	client := openai.NewClient(key)
	// 	resp, err := client.CreateChatCompletion(
	// 		context.Background(),
	// 		openai.ChatCompletionRequest{
	// 			Model: openai.GPT3Dot5Turbo,
	// 			Messages: []openai.ChatCompletionMessage{
	// 				{
	// 					Role:    openai.ChatMessageRoleUser,
	// 					Content: issue,
	// 				},
	// 			},
	// 		},
	// 	)

	// 	if err != nil {
	// 		fmt.Printf("ChatCompletion error: %v\n", err)
	// 		return
	// 	}
	// 	result <- resp.Choices[0].Message.Content
	// }()
	// select {
	// case <-time.After(60 * time.Second):
	// 	fmt.Println("timeout")
	// 	return ""
	// case res := <-result:
	// 	return res
	// }
}
func OpenAI_35_Context(msg []openai.ChatCompletionMessage, key string) string {
	result := make(chan string)
	go func() {
		client := openai.NewClient(key)
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:    openai.GPT3Dot5Turbo,
				Messages: msg,
			},
		)

		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			return
		}
		result <- resp.Choices[0].Message.Content
	}()
	select {
	case <-time.After(60 * time.Second):
		fmt.Println("timeout")
		return ""
	case res := <-result:
		return res
	}
}

type Billing struct {
	Object         string  `json:"object"`
	TotalGranted   float64 `json:"total_granted"`
	TotalUsed      float64 `json:"total_used"`
	TotalAvailable float64 `json:"total_available"`
	Grants         struct {
		Object string `json:"object"`
		Data   []struct {
			Object      string  `json:"object"`
			ID          string  `json:"id"`
			GrantAmount float64 `json:"grant_amount"`
			UsedAmount  float64 `json:"used_amount"`
			EffectiveAt float64 `json:"effective_at"`
			ExpiresAt   float64 `json:"expires_at"`
		} `json:"data"`
	} `json:"grants"`
}

// func OpenAI_Balance(key string) string {

// 	url := "https://api.openai.com/dashboard/billing/credit_grants"
// 	var data Billing
// 	body, _ := HttpGet(url)
// 	err := json.Unmarshal(body, &data)
// 	if err != nil {
// 		panic(err)
// 	}
// 	t1 := time.Unix(int64(data.Grants.Data[0].EffectiveAt), 0)
// 	t2 := time.Unix(int64(data.Grants.Data[0].ExpiresAt), 0)
// 	msg := fmt.Sprintf("💵 已用: 💲%v\n💵 剩余: 💲%v\n⏳ 有效时间: 从 %v 到 %v\n", fmt.Sprintf("%.2f", data.TotalUsed), fmt.Sprintf("%.2f", data.TotalAvailable), t1.Format("2006-01-02 15:04:05"), t2.Format("2006-01-02 15:04:05"))
// 	fmt.Println(msg)
// 	return msg
// }

func InitAiCli(key string) *resty.Client {
	// return resty.New().SetTimeout(10*time.Second).SetHeader("Authorization", fmt.Sprintf("Bearer %s", Config.ApiKey)).SetProxy(Config.HttpProxy).SetRetryCount(3).SetRetryWaitTime(2 * time.Second)

	return resty.New().SetTimeout(10*time.Second).SetHeader("Authorization", fmt.Sprintf("Bearer %s", key)).SetRetryCount(3).SetRetryWaitTime(2 * time.Second)
}
func OpenAI_Balance(key string) string {
	var data Billing
	url := "https://api.openai.com/dashboard/billing/credit_grants"
	resp, err := InitAiCli(key).R().Get(url)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		panic(err)
	}
	t1 := time.Unix(int64(data.Grants.Data[0].EffectiveAt), 0)
	t2 := time.Unix(int64(data.Grants.Data[0].ExpiresAt), 0)
	msg := fmt.Sprintf("💵 已用: 💲%v\n💵 剩余: 💲%v\n⏳ 有效时间: 从 %v 到 %v\n", fmt.Sprintf("%.2f", data.TotalUsed), fmt.Sprintf("%.2f", data.TotalAvailable), t1.Format("2006-01-02 15:04:05"), t2.Format("2006-01-02 15:04:05"))
	return msg
}
