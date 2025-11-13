CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	name VARCHAR UNIQUE NOT NULL,
	password VARCHAR NOT NULL,
	email VARCHAR NOT NULL,
	token VARCHAR,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
);

CREATE INDEX IF NOT EXISTS idx_users_name on users(name);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- 创建学习资源表
CREATE TABLE IF NOT EXISTS learn_materials (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    category VARCHAR(50) NOT NULL,
    description TEXT,
    file_url VARCHAR(255) NOT NULL,
    file_size BIGINT DEFAULT 0,
    file_type VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建分类索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_learn_materials_category ON learn_materials(category);

-- 创建创建时间索引以提高排序性能
CREATE INDEX IF NOT EXISTS idx_learn_materials_created_at ON learn_materials(created_at);

-- 创建唯一索引以避免重复的学习资源（可选，根据业务需求）
CREATE UNIQUE INDEX IF NOT EXISTS idx_learn_materials_unique ON learn_materials(title, category);

