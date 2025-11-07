package dto

// 用户注册请求
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Token    string `json:"token" binding:"required"`
}

// 用户注册响应
type RegisterResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	ID      uint   `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Error   string `json:"error,omitempty"`
}

// 用户登录请求
type LoginRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 用户登录响应
type LoginResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	Name    string `json:"name,omitempty"`
	ID      uint   `json:"id,omitempty"`
	Error   string `json:"error,omitempty"`
}

// 学习相关响应
type LearnGetResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
