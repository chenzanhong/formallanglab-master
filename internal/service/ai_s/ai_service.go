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
	systemPrompt = `你是一位形式语言与自动机理论课程的教学助教，专注于自动机、正则表达式、文法和图灵机等内容。
	
目前和学生一对一交流，回答中不要用“你们”，而应该用“你”。

请用清晰、准确、循序渐进的方式回答问题，优先解释原理和“为什么”，必要时举例说明。

如果学生提供了自动机结构、文法或测试字符串，请结合具体情境以及已有对话内容进行分析。

不要虚构定理或算法，不确定时请说明“标准理论中通常……”。

注意！对于明显超出形式语言与自动机范围的问题（如编程、其他课程内容等），请礼貌回应，例如：
> “抱歉，这个问题超出了本课程《形式语言与自动机》的范围，我无法回答该问题。”
`
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
		fmt.Println(turn.AI)
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
