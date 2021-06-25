package database

import "postapi/app/models"

// 业务逻辑，接受 Post 数据存储于数据库
func (d *DB) CreatePost(p *models.Post) error {
	res, err := d.db.Exec(insertPostSchema, p.Title, p.Content, p.Author)
	if err != nil {
		return err
	}

	res.LastInsertId()
	return err
}

// 业务逻辑，从数据库中获取所有 Post 数据
func (d *DB) GetPosts() ([]*models.Post, error) {
	var posts []*models.Post
	err := d.db.Select(&posts, "SELECT * FROM posts")
	if err != nil {
		return posts, err
	}

	return posts, nil
}
