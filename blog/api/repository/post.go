package repository

import (
	"blog/infrastructure"
	"blog/models"
)

// repository 定义业务逻辑
type PostRepository struct {
	*infrastructure.Database
}

// PostRepository 构造器
func NewPostRepository(db *infrastructure.Database) PostRepository {
	return PostRepository{db}
}

// 存储 Post
func (p *PostRepository) SavePost(post models.Post) error {
	return p.DB.Create(&post).Error
}

// 查看所有 Post
func (p *PostRepository) FindAllPosts(keyword string) (*models.Posts, int64, error) {
	var posts models.Posts
	var totalRows int64 = 0

	queryBuilder := p.DB.Order("created_at desc").Model(&models.Post{})

	// 查询参数
	if keyword != "" {
		queryKeyword := "%" + keyword + "%"
		queryBuilder = queryBuilder.Where(
			p.DB.Where("post.title LIKE ?", queryKeyword),
		)
	}

	err := queryBuilder.Find(&posts).Count(&totalRows).Error
	return &posts, totalRows, err
}

// 根据 id 查看 Post
func (p *PostRepository) FindPost(post models.Post) (models.Post, error) {
	var res models.Post
	err := p.DB.Debug().Model(&models.Post{}).Where(&post).Take(&res).Error

	return res, err
}

// 更新
func (p *PostRepository) UpdatePost(post models.Post) error {
	return p.DB.Save(&post).Error
}

// 删除 Post
func (p *PostRepository) DeletePost(post models.Post) error {
	return p.DB.Delete(&post).Error
}
