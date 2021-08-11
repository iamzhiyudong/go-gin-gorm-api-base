package routes

import (
	"github.com/gin-gonic/gin"
	"zhiyudong.cn/gin-test/controller"
	"zhiyudong.cn/gin-test/middleware"
)

func CollectRouter(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware(), middleware.RecoveryMiddleware())
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	r.GET("/api/auth/info", middleware.AuthMiddleware(), controller.Info)

	// 路由分组
	categoryRoutes := r.Group("/categories")
	categoryController := controller.NewCategoryController()

	categoryRoutes.POST("", categoryController.Create)
	categoryRoutes.PUT("/:id", categoryController.Update) // 替换
	categoryRoutes.GET("/:id", categoryController.Show)
	categoryRoutes.DELETE("/:id", categoryController.Delete)
	// categoryRoutes.PATCH("/:id", categoryController.Delete) // 局部替换

	postRoutes := r.Group("/posts")
	postRoutes.Use(middleware.AuthMiddleware())
	postController := controller.NewPostController()

	postRoutes.POST("", postController.Create)
	postRoutes.PUT("/:id", postController.Update) // 替换
	postRoutes.GET("/:id", postController.Show)
	postRoutes.DELETE("/:id", postController.Delete)
	postRoutes.GET("/page/list", postController.PageList)

	return r
}
