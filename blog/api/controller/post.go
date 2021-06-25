package controller

import (
	"blog/api/service"
	"blog/models"
	"blog/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	service *service.PostService
}

// PostController 构造器
func NewPostController(s *service.PostService) PostController {
	return PostController{s}
}

// GetPosts
func (p *PostController) GetPosts(ctx *gin.Context) {

	keyword := ctx.Query("keyword")

	data, total, err := p.service.FindAllPosts(keyword)

	if err != nil {
		util.ErrorJSON(ctx, http.StatusBadRequest, "Failed to find posts")
		return
	}
	respArr := make([]map[string]interface{}, 0)

	for _, n := range *data {
		resp := n.ResponseMap()
		respArr = append(respArr, resp)
	}

	ctx.JSON(http.StatusOK, &util.Response{
		Success: true,
		Message: "Post result set",
		Data: map[string]interface{}{
			"rows":       respArr,
			"total_rows": total,
		},
	})
}

// GetPost
func (p *PostController) GetPost(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		util.ErrorJSON(ctx, http.StatusBadRequest, "id invalid")
		return
	}

	var post models.Post
	post.ID = id
	foundPost, err := p.service.FindPost(post)
	if err != nil {
		util.ErrorJSON(ctx, http.StatusBadRequest, "Error finding post")
		return
	}
	resp := foundPost.ResponseMap()

	ctx.JSON(http.StatusOK, &util.Response{
		Success: true,
		Message: "Result set of Post",
		Data:    &resp,
	})
}

// AddPost
func (p *PostController) AddPost(ctx *gin.Context) {
	var post models.Post
	ctx.ShouldBindJSON(&post)

	if post.Title == "" {
		util.ErrorJSON(ctx, http.StatusBadRequest, "Title is required")
		return
	}
	if post.Body == "" {
		util.ErrorJSON(ctx, http.StatusBadRequest, "Body is required")
		return
	}

	err := p.service.SavePost(post)
	if err != nil {
		util.ErrorJSON(ctx, http.StatusBadRequest, "Failed to create post")
		return
	}
	util.SuccessJSON(ctx, http.StatusCreated, "Successfully Created post")
}

// DeletePost
func (p *PostController) DeletePost(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		util.ErrorJSON(ctx, http.StatusBadRequest, "id invalid")
		return
	}
	var post models.Post
	post.ID = id

	err = p.service.DeletePost(post)

	if err != nil {
		util.ErrorJSON(ctx, http.StatusBadRequest, "Failed to delete post")
		return
	}
	resp := &util.Response{
		Success: true,
		Message: "Deleted Successfully",
	}
	ctx.JSON(http.StatusOK, resp)
}

// UpdatePost
func (p *PostController) UpdatePost(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		util.ErrorJSON(ctx, http.StatusBadRequest, "id invalid")
		return
	}
	var post models.Post
	post.ID = id

	postRecord, err := p.service.FindPost(post)

	if err != nil {
		util.ErrorJSON(ctx, http.StatusBadRequest, "post with given id not found")
		return
	}
	ctx.ShouldBindJSON(&postRecord)

	if postRecord.Title == "" {
		util.ErrorJSON(ctx, http.StatusBadRequest, "Title is required")
		return
	}
	if postRecord.Body == "" {
		util.ErrorJSON(ctx, http.StatusBadRequest, "Body is required")
		return
	}

	if err := p.service.UpdatePost(postRecord); err != nil {
		util.ErrorJSON(ctx, http.StatusBadRequest, "Failed to store Post")
		return
	}
	resp := postRecord.ResponseMap()

	ctx.JSON(http.StatusOK, &util.Response{
		Success: true,
		Message: "Successfully updated Post",
		Data:    resp,
	})
}
