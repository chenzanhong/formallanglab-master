package service

import (
	"backend/internal/domain/dto"
	"backend/internal/domain/model"
	"backend/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/packages/ssestream"
)

const (
	systemPrompt = `你是一位形式语言与自动机理论课程的专业教学助教。请严格遵守以下规则回答问题：
1. **领域限定**：仅回答与形式语言、自动机（DFA/NFA/PDA）、正则表达式、上下文无关文法（CFG）、图灵机等相关的问题；拒绝回答无关话题。
2. **教学风格**：使用清晰、准确、循序渐进的语言；必要时分步骤说明；避免使用未定义或超出本科课程范围的术语。
3. **结合上下文**：若提供了当前自动机结构、测试字符串或错误信息，请基于这些具体情境作答，给出针对性解释或修正建议。
4. **不编造内容**：若不确定答案，请明确说明“根据标准理论……”，不要虚构定理、算法或状态转移。
5. **鼓励理解**：优先解释“为什么”，而非仅给出结论；可举例辅助理解。`
)

type AIService interface {
	StreamChat(ctx context.Context, username string, req *dto.AIChatRequest) (*ssestream.Stream[openai.ChatCompletionChunk], *model.AISession, error)
	SaveSession(ctx context.Context, username string, session *model.AISession) error
}

type AIServiceImpl struct {
	client *openai.Client
	repo   repository.AIRepository
}

func NewAIService(client *openai.Client, repo repository.AIRepository) AIService {
	return &AIServiceImpl{client: client, repo: repo}
}

func (s *AIServiceImpl) StreamChat(ctx context.Context, username string, req *dto.AIChatRequest) (*ssestream.Stream[openai.ChatCompletionChunk], *model.AISession, error) {
	// 1. 加载会话（不变）
	session, err := s.repo.GetSession(ctx, username, string(req.Page))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load session: %w", err)
	}
	if session == nil {
		session = &model.AISession{Page: req.Page}
	}
	session.LastActive = time.Now().Unix()

	// 2. 构造消息：严格遵循 [system] → [history] → [context] → [current question]
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
	}

	// 历史对话：直接交替添加
	for _, turn := range session.RecentTurns {
		messages = append(messages, openai.UserMessage(turn.User))
		messages = append(messages, openai.AssistantMessage(turn.AI))
	}

	// 当前页面上下文（如自动机/文法/正则）
	if ctxPrompt := s.buildContextualPrompt(req); ctxPrompt != "" {
		messages = append(messages, openai.UserMessage(ctxPrompt))
	}

	// 用户当前问题
	messages = append(messages, openai.UserMessage(req.Question))

	// 3. 调用流式 API
	stream := s.client.Chat.Completions.NewStreaming(
		ctx, openai.ChatCompletionNewParams{
			Messages: messages,
			Model:    os.Getenv("DASHSCOPE_MODEL"),
		},
	)

	return stream, session, nil
}

func (s *AIServiceImpl) SaveSession(ctx context.Context, username string, session *model.AISession) error {
	return s.repo.SaveSession(ctx, username, session)
}

func (s *AIServiceImpl) buildContextualPrompt(req *dto.AIChatRequest) string {
	var prompt string
	switch req.Page {
	case "automaton":
		if req.Automaton != nil {
			data, _ := json.MarshalIndent(req.Automaton, "", "  ")
			prompt = fmt.Sprintf("当前自动机定义如下：\n```json\n%s\n```", string(data))
		}
	case "grammar":
		if req.Grammar != nil {
			data, _ := json.MarshalIndent(req.Grammar, "", "  ")
			prompt = fmt.Sprintf("当前上下文无关文法定义如下：\n```json\n%s\n```", string(data))
		}
	case "regex":
		if req.Regex != nil {
			prompt = fmt.Sprintf("当前正则表达式为：%s", *req.Regex)
		}
	default:
	}

	return prompt
}
