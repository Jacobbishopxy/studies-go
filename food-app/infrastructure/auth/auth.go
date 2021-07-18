package auth

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
)

type AuthInterface interface {
	CreateAuth(uint64, *TokenDetails) error
	FetchAuth(string) (uint64, error)
	DeleteRefresh(string) error
	DeleteTokens(*AccessDetails) error
}

type ClientData struct {
	client *redis.Client
}

var _ AuthInterface = &ClientData{}

func NewAuth(client *redis.Client) *ClientData {
	return &ClientData{client}
}

type AccessDetails struct {
	TokenUuid string
	UserId    uint64
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	TokenUuid    string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

// 保存 token metadata 至 Redis
func (tk *ClientData) CreateAuth(userid uint64, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) // 转换 Unix 至 UTC
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	atCreated, err := tk.client.Set(td.TokenUuid, strconv.Itoa(int(userid)), at.Sub(now)).Result()
	if err != nil {
		return err
	}
	rtCreated, err := tk.client.Set(td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Result()
	if err != nil {
		return err
	}
	if atCreated == "0" || rtCreated == "0" {
		return errors.New("no record inserted")
	}

	return nil
}

// 检查 metadata 是否已保存
func (tk *ClientData) FetchAuth(tokenUuid string) (uint64, error) {
	userid, err := tk.client.Get(tokenUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)

	return userID, nil
}

// 从 Redis 中删除 token metadata
func (tk *ClientData) DeleteTokens(authD *AccessDetails) error {
	// 获取 refresh uuid
	refreshUuid := fmt.Sprintf("%s++%d", authD.TokenUuid, authD.UserId)
	// 删除 access token
	deletedAt, err := tk.client.Del(authD.TokenUuid).Result()
	if err != nil {
		return err
	}
	// 删除 refresh token
	deletedRt, err := tk.client.Del(refreshUuid).Result()
	if err != nil {
		return err
	}
	// 当记录被删除时，返回值不为 1
	if deletedAt != 1 || deletedRt != 1 {
		return errors.New("something went wrong")
	}

	return nil
}

// 仅删除 refresh token
func (tk *ClientData) DeleteRefresh(refreshUuid string) error {
	deleted, err := tk.client.Del(refreshUuid).Result()
	if err != nil || deleted == 0 {
		return err
	}

	return nil
}
