-- 创建文法数据表
CREATE TABLE IF NOT EXISTS grammars (
    id SERIAL PRIMARY KEY,
    name varchar(64) DEFAULT 'none',
    grammar JSONB NOT NULL,  -- 存储文法的详细内容，使用JSONB类型支持结构化查询
    grammar_hash CHAR(64) NOT NULL, -- SHA256 哈希值，用于内容去重
    -- grammar_type VARCHAR(50),  -- 文法类型，如CFG, RG等
    username VARCHAR NOT NULL,  -- 外键关联到users表的name字段
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- 添加外键约束
    CONSTRAINT fk_grammar_user FOREIGN KEY (username) REFERENCES users(name) ON DELETE CASCADE,
    CONSTRAINT unique_username_grammar_hash UNIQUE (username, grammar_hash)
);

-- 创建索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_grammars_username_created ON grammars(username, created_at DESC);