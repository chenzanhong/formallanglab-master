-- 创建正则表达式数据表
CREATE TABLE IF NOT EXISTS regexes (
    id SERIAL PRIMARY KEY,
    name varchar(64) DEFAULT 'none',
    pattern TEXT NOT NULL,  -- 存储正则表达式模式字符串
    username VARCHAR NOT NULL,  -- 外键关联到users表的name字段
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- 添加外键约束
    CONSTRAINT fk_regex_user FOREIGN KEY (username) REFERENCES users(name) ON DELETE CASCADE,
    CONSTRAINT unique_user_regex UNIQUE (username, pattern) -- 去重
);

-- 创建索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_regexes_username_created ON regexes(username, created_at DESC);

-- INSERT INTO regex (pattern, username) 
-- VALUES ('^a+b$', 'alice')
-- ON CONFLICT ON CONSTRAINT unique_user_regex DO NOTHING;