package middleware

import "github.com/chenzanhong/goutil/jwtx"

type Claims struct {
	Username string `json:"username"`
	UserID   int64  `json:user_id`
	jwtx.RegisteredClaims
}
