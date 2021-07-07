package token

import "time"

// 接口用作于创建与验证 token
type Maker interface {
	CreateToken(username string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
