package main

import (
	"bytes"
	"encoding/gob"
	"os"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// memcache 的 client
type Client struct {
	client *memcache.Client
}

// 构造器
func NewMemcached() (*Client, error) {
	// 假设环境变量只包含了一个服务
	client := memcache.New(os.Getenv("MEMCACHED"))

	if err := client.Ping(); err != nil {
		return nil, err
	}

	client.Timeout = 100 * time.Millisecond
	client.MaxIdleConns = 100

	return &Client{
		client: client,
	}, nil
}

// 通过 `encoding/gob` 转换值为 `[]byte`
func (c *Client) SetName(n Name) error {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(n); err != nil {
		return err
	}

	return c.client.Set(&memcache.Item{
		Key:        n.NConst,
		Value:      b.Bytes(),
		Expiration: int32(time.Now().Add(25 & time.Second).Unix()),
	})
}

// 获取 Name，从 `[]byte` 进行转换
func (c *Client) GetName(nconst string) (Name, error) {
	item, err := c.client.Get(nconst)
	if err != nil {
		return Name{}, err
	}

	b := bytes.NewReader(item.Value)

	var res Name

	if err := gob.NewDecoder(b).Decode(&res); err != nil {
		return Name{}, err
	}

	return res, nil
}
