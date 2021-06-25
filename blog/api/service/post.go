package service

import (
	"blog/api/repository"
	"blog/models"
)

// Service 层，用于内层与外层的连接（即 Repository 与 Controller）
type PostService struct {
	repository *repository.PostRepository
}

// PostService 构造器
func NewPostService(r *repository.PostRepository) PostService {
	return PostService{r}
}

// 调用 repository 的 SavePost
func (p *PostService) SavePost(post models.Post) error {
	return p.repository.SavePost(post)
}

// 调用 repository 的 FindAllPosts
func (p *PostService) FindAllPosts(keyword string) (*models.Posts, int64, error) {
	return p.repository.FindAllPosts(keyword)
}

// 调用 repository 的 FindPost
func (p *PostService) FindPost(post models.Post) (models.Post, error) {
	return p.repository.FindPost(post)
}

// 调用 repository 的 UpdatePost
func (p *PostService) UpdatePost(post models.Post) error {
	return p.repository.UpdatePost(post)
}

// 调用 repository 的 DeletePost
func (p *PostService) DeletePost(post models.Post) error {
	return p.repository.DeletePost(post)
}
