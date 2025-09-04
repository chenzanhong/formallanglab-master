package dto

type RegisterRequest struct {
    Name     string `json:"name" binding:"required,min=2,max=32"`
    Password string `json:"password" binding:"required,min=6"`
}

type RegisterResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	ID 		uint	`json:"id"`
	Name    string `json:"name"`
	Error   string `json:"error"`
}

type LoginRequest struct {
    Name     string `json:"name" binding:"required,min=2,max=32"`
    Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Token   string `json:"token"`
	Error   string `json:"error"`
	Name    string `json:"name"`
	ID		uint		`json:"id"`
	// Role    string `json:"role"`
}
