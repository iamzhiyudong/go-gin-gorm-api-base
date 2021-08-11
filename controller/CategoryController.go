package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"zhiyudong.cn/gin-test/model"
	"zhiyudong.cn/gin-test/repository"
	"zhiyudong.cn/gin-test/response"
	"zhiyudong.cn/gin-test/vo"
)

// 定义接口方便编辑器生成代码

type ICategoryController interface {
	IRestController
}

// 定义结构体方法 -- 重名

type CategoryController struct {
	Repository repository.CategoryRepository
}

func NewCategoryController() ICategoryController {
	repository := repository.NewCategoryRepository()
	repository.DB.AutoMigrate(model.Category{})

	return CategoryController{
		Repository: repository,
	}
}

func (c CategoryController) Create(ctx *gin.Context) {
	var requestCategory vo.CreateCategoryRequest
	if err := ctx.ShouldBind(&requestCategory); err != nil {
		response.Fail(ctx, nil, "名称必填")
		return
	}

	category, err := c.Repository.Create(requestCategory.Name)

	if err != nil {
		//response.Fail(ctx, nil, "分类已存在")
		panic(err)
		return
	}

	// TODO 时间的序列化 - github json issue 视频的 14：33 文章分类章节 - 自定义时间的返回格式
	response.Success(ctx, gin.H{"category": category}, "创建成功")
}

func (c CategoryController) Update(ctx *gin.Context) {
	// body 中的参数
	var requestCategory vo.CreateCategoryRequest
	if err := ctx.ShouldBind(&requestCategory); err != nil {
		response.Fail(ctx, nil, "名称必填")
		return
	}

	// 获取 path 中的参数
	categoryId, _ := strconv.Atoi(ctx.Params.ByName("id")) // 强制转换成 int

	updateCategory, err := c.Repository.SelectById(categoryId)
	if err != nil {
		response.Fail(ctx, nil, "分类不存在")
		return
	}

	// 更新分类
	// map
	// struct
	// name value
	category, updateErr := c.Repository.Update(*updateCategory, requestCategory.Name)
	if updateErr != nil {
		panic(err)
	}

	response.Success(ctx, gin.H{"category": category}, "修改成功")
}

func (c CategoryController) Show(ctx *gin.Context) {
	// 获取 path 中的参数
	categoryId, _ := strconv.Atoi(ctx.Params.ByName("id")) // 强制转换成 int

	category, err := c.Repository.SelectById(categoryId)
	if err != nil {
		response.Fail(ctx, nil, "分类不存在")
		return
	}

	response.Success(ctx, gin.H{"category": category}, "获取成功")
}

func (c CategoryController) Delete(ctx *gin.Context) {
	// 获取 path 中的参数
	categoryId, _ := strconv.Atoi(ctx.Params.ByName("id")) // 强制转换成 int

	err := c.Repository.DeleteById(categoryId)
	if err != nil {
		response.Fail(ctx, nil, "删除失败")
		return
	}

	response.Success(ctx, nil, "删除成功")
}
