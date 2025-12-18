package ai

import (
	"context"
	"note/config"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type AIService struct {
	client *openai.Client
	cfg    *config.Config
}

func NewAIService(cfg *config.Config) *AIService {
	aiConfig := openai.DefaultConfig(cfg.AIKey)
	if cfg.AIBaseURL != "" {
		aiConfig.BaseURL = cfg.AIBaseURL
	}
	return &AIService{
		client: openai.NewClientWithConfig(aiConfig),
		cfg:    cfg,
	}
}

// GenerateTitle 调用 AI 生成标题
func (s *AIService) GenerateTitle(content string) (string, error) {
	resp, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: s.cfg.AIModel,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "你是一个专业的笔记助手。请阅读以下内容，生成一个精简、有吸引力的标题（15字以内），直接输出标题，不要包含引号或前缀。",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "笔记内容：\n" + content,
				},
			},
			Temperature: 0.7,
		},
	)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}

// GenerateSummary 调用 AI 生成摘要
func (s *AIService) GenerateSummary(content string) (string, error) {
	resp, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: s.cfg.AIModel,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "请为以下笔记生成一段50字以内的简短摘要。",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}
