package main

import (
	"fmt"
	"server/common"
	"server/config"
	"server/routes"
	"server/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	config.LoadConfig()
	services.InitAuthConfig()
	common.InitDb()

	// 打印加载的配置（可选）
	fmt.Printf("Loaded config: %+v\n", config.Config)

	// 创建 Gin 引擎
	r := gin.Default()

	// 注册所有路由
	routes.RegisterRoutes(r)

	// 启动服务
	serverConfig := config.Config.Server
	err := r.Run(fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port))
	if err != nil {
		panic(err)
	}
}
