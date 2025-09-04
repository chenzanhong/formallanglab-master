package utils

import (
    "golang.org/x/crypto/bcrypt"
)

// HashPassword 对密码进行加密
func HashPassword(password string) (string, error) {
    // bcrypt.DefaultCost 是推荐的哈希成本（通常为 10）
    // 你也可以使用 bcrypt.MinCost 或自定义（如 12 提高安全性，但更慢）
    hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedBytes), nil
}

// CheckPasswordHash 验证明文密码与哈希值是否匹配
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}