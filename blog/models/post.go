package models

import "time"

// Post
type Post struct {
	ID        int64     `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:200" json:"title"`
	Body      string    `gorm:"size:3000" json:"body"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type Posts = []Post

// 表名称
func (post *Post) TableName() string {
	return "post"
}

// 响应 Post 的 map 结构
func (post *Post) ResponseMap() map[string]interface{} {
	resp := make(map[string]interface{})
	resp["id"] = post.ID
	resp["title"] = post.Title
	resp["body"] = post.Body
	resp["created_at"] = post.CreatedAt
	resp["updated_at"] = post.UpdatedAt

	return resp
}
