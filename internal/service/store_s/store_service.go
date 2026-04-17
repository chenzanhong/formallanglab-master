package store_s

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/chenzanhong/formallanglab-master/internal/domain/dto"
	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
	"github.com/chenzanhong/formallanglab-master/internal/domain/storage"
	"github.com/chenzanhong/formallanglab-master/internal/repository"
)

type StoreService interface {
	// 自动机
	CreateAutomaton(ctx context.Context, a *storage.AutomatonRecord) error
	FindAutomatonByUsername(ctx context.Context, username string) ([]model.Automaton, error) // sql 时只返回 content 字段
	FindAutomata(ctx context.Context, username string, lastID uint, limit int) (dto.FindAutomatonResponse, error)
	DeleteAutomaton(ctx context.Context, id string, username string) error

	// 文法
	CreateGrammar(ctx context.Context, g *storage.GrammarRecord) error
	FindGrammarByUsername(ctx context.Context, username string) ([]model.Grammar, error) // sql 时只返回 content 字段
	FindGrammars(ctx context.Context, username string, lastID uint, limit int) (dto.FindGrammarResponse, error)
	DeleteGrammar(ctx context.Context, id string, username string) error

	// 正则表达式
	CreateRegex(ctx context.Context, r *storage.RegexRecord) error
	FindRegexes(ctx context.Context, username string, lastID uint, limit int) (dto.FindRegexResponse, error)
	DeleteRegex(ctx context.Context, id string, username string) error
}

type StoreServiceImpl struct {
	storeRepo repository.StoreRepository
}

func NewStoreService(storeRepo repository.StoreRepository) StoreService {
	return &StoreServiceImpl{storeRepo: storeRepo}
}

// CreateAutomaton 创建自动机记录
func (s *StoreServiceImpl) CreateAutomaton(ctx context.Context, a *storage.AutomatonRecord) error {
	// 在保存前调用 Validate 函数确保自动机设置了 type，不需要校验 err，允许结构不完整
	_ = a.Automaton.Validate()

	data, err := json.Marshal(a.Automaton)
	if err != nil {
		return err
	}
	sum := sha256.Sum256(data)
	a.AutomatonHash = hex.EncodeToString(sum[:])

	// 根据 Validate 后的 Automaton.Type 设置 AutomatonType 字段
	switch a.Automaton.Type {
	case model.DFA:
		a.AutomatonType = model.DFA.String()
	case model.NFA:
		a.AutomatonType = model.NFA.String()
	case model.EpsilonNFA:
		a.AutomatonType = model.EpsilonNFA.String()
	default:
		a.AutomatonType = model.DefaultAutomatonType.String()
	}

	return s.storeRepo.CreateAutomaton(ctx, a)
}

// FindAutomatonByUsername 根据用户名查询自动机列表
func (s *StoreServiceImpl) FindAutomatonByUsername(ctx context.Context, username string) ([]model.Automaton, error) {
	// 调用 repository 层方法获取存储记录
	records, err := s.storeRepo.FindAutomatonByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	// 转换为 model.Automaton
	automata := make([]model.Automaton, 0, len(records))
	for _, record := range records {
		automata = append(automata, record.Automaton)
	}

	return automata, nil
}

// CreateGrammar 创建文法记录
func (s *StoreServiceImpl) CreateGrammar(ctx context.Context, g *storage.GrammarRecord) error {
	data, err := json.Marshal(g.Grammar)
	if err != nil {
		return err
	}
	sum := sha256.Sum256(data)
	g.GrammarHash = hex.EncodeToString(sum[:])

	return s.storeRepo.CreateGrammar(ctx, g)
}

// FindGrammarByUsername 根据用户名查询文法列表
func (s *StoreServiceImpl) FindGrammarByUsername(ctx context.Context, username string) ([]model.Grammar, error) {
	// 调用 repository 层方法获取存储记录
	records, err := s.storeRepo.FindGrammarByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	// 转换为 model.Grammar
	grammars := make([]model.Grammar, 0, len(records))
	for _, record := range records {
		grammars = append(grammars, record.Grammar)
	}

	return grammars, nil
}

// FindAutomata 分页查询自动机列表
func (s *StoreServiceImpl) FindAutomata(ctx context.Context, username string, lastID uint, limit int) (dto.FindAutomatonResponse, error) {
	records, hasMore, err := s.storeRepo.FindAutomatonsAfterID(ctx, username, lastID, limit)
	if err != nil {
		return dto.FindAutomatonResponse{}, err
	}

	automata := make([]dto.AutomatonItem, 0, len(records))
	var nextCursor uint
	for _, record := range records {
		automata = append(automata, dto.AutomatonItem{
			ID:        record.ID,
			Name:      record.Name,
			Automaton: record.Automaton,
			CreatedAt: record.CreatedAt,
		})
		nextCursor = uint(record.ID)
	}

	return dto.FindAutomatonResponse{
		Automatons: automata,
		HasMore:    hasMore,
		NextCursor: nextCursor,
	}, nil
}

// FindGrammars 分页查询文法列表
func (s *StoreServiceImpl) FindGrammars(ctx context.Context, username string, lastID uint, limit int) (dto.FindGrammarResponse, error) {
	records, hasMore, err := s.storeRepo.FindGrammarsAfterID(ctx, username, lastID, limit)
	if err != nil {
		return dto.FindGrammarResponse{}, err
	}

	grammars := make([]dto.GrammarItem, 0, len(records))
	var nextCursor uint
	for _, record := range records {
		grammars = append(grammars, dto.GrammarItem{
			ID:        record.ID,
			Name:      record.Name,
			Grammar:   record.Grammar,
			CreatedAt: record.CreatedAt,
		})
		nextCursor = uint(record.ID)
	}

	return dto.FindGrammarResponse{
		Grammars:   grammars,
		HasMore:    hasMore,
		NextCursor: nextCursor,
	}, nil
}

// CreateRegex 创建正则表达式记录
func (s *StoreServiceImpl) CreateRegex(ctx context.Context, r *storage.RegexRecord) error {
	return s.storeRepo.CreateRegex(ctx, r)
}

// FindRegexes 根据用户名查询正则表达式列表
func (s *StoreServiceImpl) FindRegexes(
	ctx context.Context,
	username string,
	lastID uint,
	limit int,
) (dto.FindRegexResponse, error) {
	records, hasMore, err := s.storeRepo.FindRegexesAfterID(ctx, username, lastID, limit)
	if err != nil {
		return dto.FindRegexResponse{}, err
	}

	items := make([]dto.RegexItem, len(records))
	for i, rec := range records {
		items[i] = dto.RegexItem{
			ID:        rec.ID,
			Name:      rec.Name,
			Pattern:   rec.Pattern, // model.Regex
			CreatedAt: rec.CreatedAt,
		}
	}

	var nextCursor uint
	if hasMore && len(items) > 0 {
		nextCursor = items[len(items)-1].ID
	}

	return dto.FindRegexResponse{
		Patterns:   items,
		HasMore:    hasMore,
		NextCursor: nextCursor,
	}, nil
}

// DeleteAutomaton 删除自动机记录
func (s *StoreServiceImpl) DeleteAutomaton(ctx context.Context, id string, username string) error {
	return s.storeRepo.DeleteAutomatonByID(ctx, id, username)
}

// DeleteGrammar 删除文法记录
func (s *StoreServiceImpl) DeleteGrammar(ctx context.Context, id string, username string) error {
	return s.storeRepo.DeleteGrammarByID(ctx, id, username)
}

// DeleteRegex 删除正则表达式记录
func (s *StoreServiceImpl) DeleteRegex(ctx context.Context, id string, username string) error {
	return s.storeRepo.DeleteRegexByID(ctx, id, username)
}
