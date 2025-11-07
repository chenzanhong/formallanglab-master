package model

import (
	"backend/pkg/binding"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	AISessionTTL = 60 * time.Minute // AI 会话过期时间
	MaxTurns     = 10               // 最大对话记录数
)

type AISession struct {
	Page        binding.PageType `json:"page"`
	LastActive  int64            `json:"last_active"`            // 最后活跃时间
	RecentTurns []QAPair         `json:"recent_turns",omitempty` // 最近的对话记录
	Summary     string           `json:"summary,omitempty"`      // 会话总结
}

type QAPair struct {
	User string `json:"user"` // 用户输入
	AI   string `json:"ai"`   // AI 回复
}

func (s *AISession) Trim() {
	if len(s.RecentTurns) > MaxTurns {
		s.RecentTurns = s.RecentTurns[len(s.RecentTurns)-MaxTurns:]
	}
}

// 存入redis
func (s *AISession) Save(redisClient *redis.Client, sessionID string) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return redisClient.Set(context.Background(), "ai_sess:"+sessionID, data, 60*time.Minute).Err()
}

func SaveAISession(redisClient *redis.Client, sessionID string, session *AISession) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return redisClient.Set(context.Background(), "ai_sess:"+sessionID, data, 60*time.Minute).Err()
}

func (s AISession) LoadAISession(redisClient *redis.Client, sessionID string) (*AISession, error) {
	val, err := redisClient.Get(context.Background(), "ai_sess:"+sessionID).Result()
	if err == redis.Nil {
		return nil, nil
	}
	var sess AISession
	if err := json.Unmarshal([]byte(val), &sess); err != nil {
		return nil, err
	}
	return &sess, nil
}
