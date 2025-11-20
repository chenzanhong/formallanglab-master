-- 创建自动机数据表
CREATE TABLE IF NOT EXISTS automatons (
    id SERIAL PRIMARY KEY,
    name varchar(64) DEFAULT 'none',
    automaton JSONB NOT NULL,  -- 存储自动机的详细内容，使用JSONB类型支持结构化查询
    automaton_hash CHAR(64) NOT NULL, -- SHA256 哈希值，用于内容去重
    -- automaton_type VARCHAR(50) NOT NULL,  -- 自动机类型，如DFA, NFA等
    username VARCHAR NOT NULL,  -- 外键关联到users表的name字段
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- 添加外键约束
    CONSTRAINT fk_automaton_user FOREIGN KEY (username) REFERENCES users(name) ON DELETE CASCADE,
    -- 👇 关键：同一用户不能重复存储相同内容的自动机
    CONSTRAINT unique_username_automaton_hash UNIQUE (username, automaton_hash)
);

-- 创建索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_automatons_username_created ON automatons(username, created_at DESC);