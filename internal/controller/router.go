package controller

import (
	"dify-sandbox-win/internal/middleware"
	"dify-sandbox-win/internal/static"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Setup(Router *gin.Engine) {
	PublicGroup := Router.Group("")
	PrivateGroup := Router.Group("/v1/sandbox/")

	PrivateGroup.Use(middleware.Auth())

	{
		// health check
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})
	}

	InitRunRouter(PrivateGroup)
	InitDependencyRouter(PrivateGroup)
}

func InitDependencyRouter(Router *gin.RouterGroup) {
	dependencyRouter := Router.Group("dependencies")
	{
		//更新依赖
		dependencyRouter.GET("", GetDependencies)
		dependencyRouter.POST("/update", UpdateDependencies)
		dependencyRouter.GET("/refresh", RefreshDependencies)
	}
}

func InitRunRouter(Router *gin.RouterGroup) {
	runRouter := Router.Group("")
	{
		//判断是否达到最大请求数和最大worker数
		runRouter.POST(
			"/run",
			middleware.MaxRequest(static.GetDifySandboxGlobalConfigurations().MaxRequests),
			middleware.MaxWorker(static.GetDifySandboxGlobalConfigurations().MaxWorkers),
			RunSandboxController,
		)
	}
}
