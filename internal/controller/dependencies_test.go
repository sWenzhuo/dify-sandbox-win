package controller

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetDependencies(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 准备输入参数
	body := map[string]string{
		"language": "python3",
	}
	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, "/v1/sandbox/dependencies", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	// 调用 handler
	GetDependencies(c)
	fmt.Println(w.Body.String())

}

func TestUpdateDependencies(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// 构造 JSON 请求体
	body := map[string]string{
		"language": "python3",
	}
	jsonData, _ := json.Marshal(body)

	// 构造 HTTP 请求和响应
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 构造 gin.Context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// 调用 handler
	UpdateDependencies(c)

	// 断言状态码
	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("UpdateDependencies response: %s", w.Body.String())
}

func TestRefreshDependencies(t *testing.T) {
	fmt.Println("Hello World")
}
