package gpt

import (
	"GPT-Integrate/utils/request"
	"context"
	"encoding/json"
	"time"

	"github.com/sashabaranov/go-openai"
)

const (
	BALANCE_URL = "https://api.openai.com/dashboard/billing/credit_grants"
)

type GptChatCompletion struct {
	Key   string                       `json:"key"`
	Field openai.ChatCompletionRequest `json:"field"`
}
type GptCompletion struct {
	Key   string                   `json:"key"`
	Field openai.CompletionRequest `json:"field"`
}

type Balance struct {
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

func GptRequestChatCompletion(param GptChatCompletion, timeout int) string {
	// 创建ctx
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	var res = make(chan openai.ChatCompletionResponse)
	defer cancel()
	go func() {
		// 创建client
		client := openai.NewClient(param.Key)
		// 发送请求
		resp, err := client.CreateChatCompletion(ctx, param.Field)
		if err != nil {
			return
		}
		res <- resp
	}()
	select {
	case <-ctx.Done():
		return ""
	case re := <-res:
		return re.Choices[0].Message.Content
	}
}
func GptRequestCompletion(param GptCompletion, timeout int) string {
	// 创建ctx
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	var res = make(chan openai.CompletionResponse)
	defer cancel()
	go func() {
		// 创建client
		client := openai.NewClient(param.Key)
		// 发送请求
		resp, err := client.CreateCompletion(ctx, param.Field)
		if err != nil {
			return
		}
		res <- resp
	}()
	select {
	case <-ctx.Done():
		return ""
	case re := <-res:
		return re.Choices[0].Text
	}
}

func GptRequest(param interface{}, timeout int) string {
	switch p := param.(type) {
	case GptChatCompletion:
		return GptRequestChatCompletion(p, timeout)
	case GptCompletion:
		return GptRequestCompletion(p, timeout)
	default:
		return ""
	}
}

func GetBalance(key string,timeout int) (Balance, error) {
	var balance Balance
	header := map[string]string{
		"Authorization": "Bearer " + key,
	}

	body, err := request.HttpGet(BALANCE_URL, header, timeout)
	if err != nil {
		return balance, err
	}
	err = json.Unmarshal(body, &balance)
	if err != nil {
		return balance, err
	}
	return balance, nil
}
