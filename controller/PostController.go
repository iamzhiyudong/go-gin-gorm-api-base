package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"zhiyudong.cn/gin-test/common"
	"zhiyudong.cn/gin-test/model"
	"zhiyudong.cn/gin-test/response"
	"zhiyudong.cn/gin-test/vo"
)

type IPostController interface {
	IRestController
	PageList(ctx *gin.Context)
}

type PostController struct {
	DB *gorm.DB
}

func NewPostController() IPostController {
	db := common.GetDB()
	db.AutoMigrate(model.Post{})
	return PostController{DB: db}
}

func (p PostController) PageList(ctx *gin.Context) {
	// 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	// 分页
	var posts []model.Post
	p.DB.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&posts)

	// 总数
	var total int64
	p.DB.Model(model.Post{}).Count(&total)

	response.Success(ctx, gin.H{"data": posts, "total": total}, "获取成功")

}

func (p PostController) Create(ctx *gin.Context) {
	var requestPost vo.CreatePostRequest
	// 数据验证
	if err := ctx.ShouldBind(&requestPost); err != nil {
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// 获取登录用户
	user, _ := ctx.Get("user")

	// 创建文章
	post := model.Post{
		UserId:     user.(model.User).ID,
		CategoryId: requestPost.CategoryId,
		Title:      requestPost.Title,
		HeadImg:    requestPost.HeadImg,
		Content:    requestPost.Content,
	}

	if err := p.DB.Create(&post).Error; err != nil {
		panic(err)
		return
	}

	response.Success(ctx, gin.H{"post": post}, "创建成功")
}

func (p PostController) Update(ctx *gin.Context) {
	var requestPost vo.CreatePostRequest
	// 数据验证
	if err := ctx.ShouldBind(&requestPost); err != nil {
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// 获取 path 中的 id
	postId := ctx.Params.ByName("id")

	var post model.Post
	if p.DB.Where("id = ?", postId).First(&post).Error != nil {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	// 判断当前用户是否为文章的作者
	// 获取登录用户
	user, _ := ctx.Get("user")
	userId := user.(model.User).ID
	if userId != post.UserId {
		response.Fail(ctx, nil, "文章不属于您，请勿非法操作")
		return
	}

	// 更新文章  TODO updates 无法智能选择字段 需要先构建出来
	if err := p.DB.Model(&post).Updates(&model.Post{CategoryId: requestPost.CategoryId, Title: requestPost.Title, HeadImg: requestPost.HeadImg, Content: requestPost.Content}).Error; err != nil {
		response.Fail(ctx, nil, "更新失败")
		return
	}

	response.Success(ctx, gin.H{"post": post}, "更新成功")
}

func (p PostController) Show(ctx *gin.Context) {
	// 获取 path 中的 id
	postId := ctx.Params.ByName("id")

	var post model.Post
	if p.DB.Preload("Category").Where("id = ?", postId).First(&post).Error != nil {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	response.Success(ctx, gin.H{"post": post}, "获取成功")
}

func (p PostController) Delete(ctx *gin.Context) {
	// 获取 path 中的 id
	postId := ctx.Params.ByName("id")

	var post model.Post
	if p.DB.Where("id = ?", postId).First(&post).Error != nil {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	// 判断当前用户是否为文章的作者
	// 获取登录用户
	user, _ := ctx.Get("user")
	userId := user.(model.User).ID
	if userId != post.UserId {
		response.Fail(ctx, nil, "文章不属于您，请勿非法操作")
		return
	}

	result := p.DB.Delete(&post)
	if result.Error != nil || result.RowsAffected == 0 {
		response.Fail(ctx, nil, "删除失败")
		return
	}

	response.Success(ctx, gin.H{"post": post}, "删除成功")
}
