package repository

import (
	"backend/internal/domain/storage"
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StoreRepository interface {
	// 只要创建和查询即可，CreateXXX、FindXXXByUsername
	// 自动机
	CreateAutomaton(ctx context.Context, a *storage.AutomatonRecord) error
	FindAutomatonByUsername(ctx context.Context, username string) ([]storage.AutomatonRecord, error) // sql时只返回content和id字段
	FindAutomatonsAfterID(ctx context.Context, username string, lastID uint, limit int) ([]storage.AutomatonRecord, bool, error)
	DeleteAutomatonByID(ctx context.Context, id string, username string) error

	// 文法
	CreateGrammar(ctx context.Context, g *storage.GrammarRecord) error
	FindGrammarByUsername(ctx context.Context, username string) ([]storage.GrammarRecord, error) // sql时只返回content和id字段
	FindGrammarsAfterID(ctx context.Context, username string, lastID uint, limit int) ([]storage.GrammarRecord, bool, error)
	DeleteGrammarByID(ctx context.Context, id string, username string) error

	// 正则表达式
	CreateRegex(ctx context.Context, r *storage.RegexRecord) error
	FindRegexesByUsername(ctx context.Context, username string) ([]storage.RegexRecord, error)
	FindRegexesAfterID(ctx context.Context, username string, lastID uint, limit int) ([]storage.RegexRecord, bool, error) // sql时只返回content和id字段
	DeleteRegexByID(ctx context.Context, id string, username string) error
}

type StoreRepositoryImpl struct {
	db *gorm.DB
}

// NewStoreRepository 创建存储仓库实例
func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &StoreRepositoryImpl{
		db: db,
	}
}

// CreateAutomaton 创建自动机记录
func (r *StoreRepositoryImpl) CreateAutomaton(ctx context.Context, a *storage.AutomatonRecord) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "username"}, {Name: "automaton_hash"}},
			DoNothing: true,
		}).
		Create(a).
		Error
}

// FindAutomatonByUsername 根据用户名查询自动机列表
func (r *StoreRepositoryImpl) FindAutomatonByUsername(ctx context.Context, username string) ([]storage.AutomatonRecord, error) {
	var records []storage.AutomatonRecord
	if err := r.db.WithContext(ctx).Where("username = ?", username).Select("id, name, automaton, created_at").Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// CreateGrammar 创建文法记录
func (r *StoreRepositoryImpl) CreateGrammar(ctx context.Context, g *storage.GrammarRecord) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "username"}, {Name: "grammar_hash"}},
			DoNothing: true,
		}).
		Create(g).
		Error
}

// FindGrammarByUsername 根据用户名查询文法列表
func (r *StoreRepositoryImpl) FindGrammarByUsername(ctx context.Context, username string) ([]storage.GrammarRecord, error) {
	var records []storage.GrammarRecord
	if err := r.db.WithContext(ctx).Where("username = ?", username).Select("id, name, grammar, created_at").Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// FindAutomataAfterID 根据用户名和lastID分页查询自动机列表
func (r *StoreRepositoryImpl) FindAutomatonsAfterID(ctx context.Context, username string, lastID uint, limit int) ([]storage.AutomatonRecord, bool, error) {
	var records []storage.AutomatonRecord
	query := r.db.WithContext(ctx).Where("username = ?", username).Select("id, name, automaton, created_at")

	if lastID > 0 {
		query = query.Where("id > ?", lastID)
	}

	// 查询limit+1条记录来判断是否有更多
	if err := query.Order("id ASC").Limit(limit + 1).Find(&records).Error; err != nil {
		return nil, false, err
	}

	hasMore := len(records) > limit
	if hasMore {
		records = records[:limit]
	}

	return records, hasMore, nil
}

// FindGrammarsAfterID 根据用户名和lastID分页查询文法列表
func (r *StoreRepositoryImpl) FindGrammarsAfterID(ctx context.Context, username string, lastID uint, limit int) ([]storage.GrammarRecord, bool, error) {
	var records []storage.GrammarRecord
	query := r.db.WithContext(ctx).Where("username = ?", username).Select("id, grammar, created_at")

	if lastID > 0 {
		query = query.Where("id > ?", lastID)
	}

	// 查询limit+1条记录来判断是否有更多
	if err := query.Order("id ASC").Limit(limit + 1).Find(&records).Error; err != nil {
		return nil, false, err
	}

	hasMore := len(records) > limit
	if hasMore {
		records = records[:limit]
	}

	return records, hasMore, nil
}

// CreateRegex 创建正则表达式记录
func (r *StoreRepositoryImpl) CreateRegex(ctx context.Context, re *storage.RegexRecord) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "username"}, {Name: "pattern"}},
			DoNothing: true,
		}).
		Create(re).
		Error
}

func (r *StoreRepositoryImpl) FindRegexesByUsername(ctx context.Context, username string) ([]storage.RegexRecord, error) {
	var records []storage.RegexRecord
	if err := r.db.WithContext(ctx).Where("username = ?", username).Select("id, pattern").Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (r *StoreRepositoryImpl) FindRegexesAfterID(
	ctx context.Context,
	username string,
	lastID uint,
	limit int,
) ([]storage.RegexRecord, bool, error) {
	var records []storage.RegexRecord
	query := r.db.WithContext(ctx).Where("username = ?", username)
	if lastID > 0 {
		query = query.Where("id < ?", lastID)
	}
	err := query.
		Order("id DESC").
		Limit(limit + 1).
		Find(&records).Error
	if err != nil {
		return nil, false, err
	}

	hasMore := len(records) > limit
	if hasMore {
		records = records[:limit]
	}
	return records, hasMore, nil
}

// DeleteAutomatonByID 根据id和用户名删除自动机记录
func (r *StoreRepositoryImpl) DeleteAutomatonByID(ctx context.Context, id string, username string) error {
	return r.db.WithContext(ctx).Where("id = ? AND username = ?", id, username).Delete(&storage.AutomatonRecord{}).Error
}

// DeleteGrammarByID 根据id和用户名删除文法记录
func (r *StoreRepositoryImpl) DeleteGrammarByID(ctx context.Context, id string, username string) error {
	return r.db.WithContext(ctx).Where("id = ? AND username = ?", id, username).Delete(&storage.GrammarRecord{}).Error
}

// DeleteRegexByID 根据id和用户名删除正则表达式记录
func (r *StoreRepositoryImpl) DeleteRegexByID(ctx context.Context, id string, username string) error {
	return r.db.WithContext(ctx).Where("id = ? AND username = ?", id, username).Delete(&storage.RegexRecord{}).Error
}
