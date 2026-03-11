
-- 创建学习资源表
CREATE TABLE IF NOT EXISTS learn_materials (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL,  -- grammar/automaton/regex/general
    file_key VARCHAR(512) NOT NULL UNIQUE,  -- OSS中的文件路径
    file_name VARCHAR(255) NOT NULL,  -- 原始文件名
    mime_type VARCHAR(100),  -- 文件MIME类型
    size_bytes BIGINT DEFAULT 0,  -- 文件大小（字节）
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建分类索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_learn_materials_category ON learn_materials(category);

-- 创建创建时间索引以提高排序性能
CREATE INDEX IF NOT EXISTS idx_learn_materials_created_at ON learn_materials(created_at);

-- 创建file_key索引以加速OSS文件查询
CREATE INDEX IF NOT EXISTS idx_learn_materials_file_key ON learn_materials(file_key);
