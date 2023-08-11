package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"io"
	"matuto.cc/go-gin-study/routers"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func main() {
	// 禁用控制台颜色
	gin.DisableConsoleColor()
	// 创建一个不包含中间件的路由器
	r := new(routers.Routers).InitRouter()

	//创建日志文件
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)
	// 如果需要将日志同时写入文件和控制台，请使用以下代码
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	//r := gin.Default()
	// 加载html路径
	// 使用LoadHTMLGlob() 或者 LoadHTMLFiles()
	//	r.LoadHTMLGlob("templates/*")
	//	r.GET("/index", func(c *gin.Context) {
	//		c.HTML(http.StatusOK, "index.html", gin.H{"title": "Main website"})
	//	})

	r.Run()
}

func TestHelloRoute(t *testing.T) {
	r := new(routers.Routers).InitRouter()
	w := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/hello", nil)
	r.ServeHTTP(w, request)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "go", w.Body.String())
}
