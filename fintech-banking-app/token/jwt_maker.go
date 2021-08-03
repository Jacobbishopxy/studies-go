package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTMaker struct {
	secretKey string
}

const minSecreteKeySize = 32

// 为具体类型 JWTMaker 实现 Maker 接口
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecreteKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecreteKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// 是否过期的验证
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

// 实现 CreateToken，其中 `jwt.NewWithClaims` 的 payload 参数需要实现 `Valid()` 方法
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	// 转换 secretKey 字符串为 []byte
	return jwtToken.SignedString([]byte(maker.secretKey))
}

// 实现 VerifyToken
// 需要通过 `jwt.ParseWithClaims` 来解析 token，接着传入 token 字符串，一个空的 `Payload` 对象，与一个 key 函数
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// 通过 token.Method 来获取签名算法
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	// 为了区别错误类型，我们需要转换从 ParseWithClaims() 函数返回的错误，成为 jwt.ValidationError
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
