package controller

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"zhiyudong.cn/gin-test/common"
	"zhiyudong.cn/gin-test/dto"
	"zhiyudong.cn/gin-test/model"
	"zhiyudong.cn/gin-test/response"
	"zhiyudong.cn/gin-test/util"
)

func Register(ctx *gin.Context) {
	DB := common.GetDB()

	/*
		使用 map / 结构体 / Bind 来获取请求中的参数
		可应用于前端 axios 请求中的 json 格式请求体
	*/

	// 使用 map 获取请求的参数
	// var requestMap = make(map[string]string)
	// json.NewDecoder(ctx.Request.Body).Decode(&requestMap)

	// 使用结构体获取参数
	// var requestUser = model.User{}
	// json.NewDecoder(ctx.Request.Body).Decode(&requestUser)

	// 使用 gin 的 bind 获取请求的参数 -- axios 请求报错
	var requestUser = model.User{}
	ctx.Bind(&requestUser)

	// 获取参数
	name := requestUser.Name
	telephone := requestUser.Telephone
	password := requestUser.Password

	// 获取参数 -- data-form 格式的请求体
	// name := ctx.PostForm("name")
	// telephone := ctx.PostForm("telephone")
	// password := ctx.PostForm("password")

	// 数据验证
	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 40022, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 40022, nil, "密码不能少于6位")
		return
	}
	if len(name) == 0 {
		name = util.RandomString(10)
	}

	log.Println(name, telephone, password)

	// 判断手机号是否存在

	if isTelephoneExist(DB, telephone) {
		response.Response(ctx, http.StatusUnprocessableEntity, 40022, nil, "用户已经存在")
		return
	}

	// 创建用户

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 50000, nil, "密码加密错误")
		// ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "密码加密错误"})
		return
	}

	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hashedPassword),
	}

	DB.Create(&newUser)

	// 返回结果
	// response.Success(ctx, nil, "注册成功")

	// 发放 token -- 注册成功后返回 token
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 50000, nil, "系统错误")
		log.Printf("token generate error: %v", err)
		return
	}

	// 返回结果
	response.Success(ctx, gin.H{"token": token}, "注册成功")
}

func Login(ctx *gin.Context) {
	db := common.GetDB()

	// 使用 gin 的 bind 获取请求的参数 -- axios 请求报错
	var requestUser = model.User{}
	ctx.Bind(&requestUser)

	// 获取参数
	telephone := requestUser.Telephone
	password := requestUser.Password

	// 获取参数
	// telephone := ctx.PostForm("telephone")
	// password := ctx.PostForm("password")

	// 数据验证
	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 40022, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 40022, nil, "密码不能少于6位")
		return
	}

	// 判断手机号是否存在
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 40022, nil, "用户不存在")
		return
	}

	// 判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Fail(ctx, nil, "密码错误")
		return
	}

	// 发放 token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 50000, nil, "系统错误")
		log.Printf("token generate error: %v", err)
		return
	}

	// 返回结果
	response.Success(ctx, gin.H{"token": token}, "登录成功")
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user.(model.User))}}) // 类型断言 - 判断类型是否匹配
}
