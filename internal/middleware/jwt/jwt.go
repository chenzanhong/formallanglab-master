package jwt

import "github.com/chenzanhong/goutil/jwtx"

type Claims struct {
	Username string `json:"username" inject:"username"`
	UserID   int64  `json:"user_id"       inject:"user_id"`
	jwtx.RegisteredClaims
}
