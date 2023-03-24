package gpt

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"gpt_bot/biz/conf"
	"gpt_bot/biz/model"
	"math/rand"
)

var openaiCli *openai.Client

func init() {
	openaiCli = openai.NewClient(conf.GetConf().Gpt.ApiKey)
}

func Chat(ctx context.Context, userID string, contents []*model.ChatContent) (openai.ChatCompletionMessage, error) {
	msgs := make([]openai.ChatCompletionMessage, 0, len(contents))
	for _, content := range contents {
		msgs = append(msgs, content.ToChatCompletionMessage())
	}
	req := openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		Messages:    msgs,
		Temperature: 0.5,
		TopP:        rand.New(rand.NewSource(1)).Float32(),
		User:        userID,
	}
	resp, err := openaiCli.CreateChatCompletion(ctx, req)
	if err != nil {
		return openai.ChatCompletionMessage{}, err
	}
	if len(resp.Choices) == 0 {
		return openai.ChatCompletionMessage{}, fmt.Errorf("resp.Choices is nil")
	}
	return resp.Choices[0].Message, nil
}
