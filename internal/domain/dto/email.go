package dto

// 发送注册验证码请求
type SendRegisterVerificationCodeRequest struct {
	Email string `json:"email" binding:"required"`
}

type SendRegisterVerificationCodeResponse struct {
	Msg    string `json:"msg"`
	Result bool   `json:"result"`
}

// 发送重置密码验证码请求
type SendResetPasswordVerificationCodeRequest struct {
	Email string `json:"email" binding:"required"`
}

type SendResetPasswordVerificationCodeResponse struct {
	Msg    string `json:"msg"`
	Result bool   `json:"result"`
}


// 验证码响应
type VerificationCodeResponse struct {
	Result bool   `json:"result"`
	Msg    string `json:"msg"`
}
