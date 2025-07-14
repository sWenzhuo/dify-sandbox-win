package controller

import (
	runner_types "dify-sandbox-win/internal/core/runner/types"
	"dify-sandbox-win/internal/service"
	"dify-sandbox-win/internal/types"
	"github.com/gin-gonic/gin"
)

func RunSandboxController(c *gin.Context) {
	BindRequest(c, func(req struct {
		Language      string `json:"language" form:"language" binding:"required"`
		Code          string `json:"code" form:"code" binding:"required"`
		Preload       string `json:"preload" form:"preload"`
		EnableNetwork bool   `json:"enable_network" form:"enable_network"`
	}) {
		switch req.Language {
		case "python3":
			c.JSON(200, service.RunPython3Code(req.Code, req.Preload, &runner_types.RunnerOptions{
				EnableNetwork: req.EnableNetwork,
			}))
		default:
			c.JSON(400, types.ErrorResponse(-400, "unsupported language"))
		}
	})
}

func GetDependencies(c *gin.Context) {
	BindRequest(c, func(req struct {
		Language string `json:"language" form:"language" binding:"required"`
	}) {
		switch req.Language {
		case "python3":
			c.JSON(200, service.ListPython3Dependencies())
		default:
			c.JSON(400, types.ErrorResponse(-400, "unsupported language"))
		}
	})
}

func UpdateDependencies(c *gin.Context) {
	BindRequest(c, func(req struct {
		Language string `json:"language" form:"language" binding:"required"`
	}) {
		switch req.Language {
		case "python3":
			c.JSON(200, service.UpdateDependencies())
		default:
			c.JSON(400, types.ErrorResponse(-400, "unsupported language"))
		}
	})
}

func RefreshDependencies(c *gin.Context) {
	BindRequest(c, func(req struct {
		Language string `json:"language" form:"language" binding:"required"`
	}) {
		switch req.Language {
		case "python3":
			c.JSON(200, service.RefreshPython3Dependencies())
		default:
			c.JSON(400, types.ErrorResponse(-400, "unsupported language"))
		}
	})
}
