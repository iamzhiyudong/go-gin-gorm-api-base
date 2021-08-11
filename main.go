package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
	"zhiyudong.cn/gin-test/common"
	"zhiyudong.cn/gin-test/routes"
)

func main() {
	InitConfig()

	db := common.InitDB()

	dbInstance, dbErr := db.DB()
	if dbErr != nil {
		panic("failed to init db")
		return
	} else {
		defer dbInstance.Close()
	}

	r := gin.Default()
	r = routes.CollectRouter(r)

	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	}

	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
