# Go - Gin入门

## 创建项目

    这就不介绍了 都会，使用 Golang或者VSCode/Fleet都行

## 初始化项目

    初始化 go mod，在项目目录下执行 ```go mod init project``` , 项目目录下会自动创建 go.mod 文件，所有的依赖都会在里边

## 导入 gin 依赖

    在项目目录下执行```go get "github.com/gin-gonic/gin"``` 导入gin依赖

## 创建项目启动入口

    创建 main.go文件，编写main方法。

```go
package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    r.Run() // listen and serve on 0.0.0.0:8080
}
```

项目默认端口 8080 ， 可修改 ```r.Run(":7070") ```

## 测试

访问路径http://127.0.0.1:8080/ping

## 接口示例

### 获取路径参数

```go
package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // 此规则能够匹配/user/matuto 这种格式，但不能匹配/user/ 或 /user这种格式
    r.GET("/user/:name", func(c *gin.Context) {
        name := c.Param("name")
        c.String(http.StatusOK, "Hello %s", name)
    })

    // 但是，这个规则既能匹配/user/matuto/格式也能匹配/user/matuto/send这种格式
    // 如果没有其他路由器匹配/user/matuto，它将重定向到/user/matuto/
    r.GET("/user/:name/*action", func(c *gin.Context) {
        name := c.Param("name")
        action := c.Param("action")
        message := name + " is " + action
        c.String(http.StatusOK, message)
    })
    r.Run() 
}
```

### 获取get参数

```go
package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    //获取get参数 /user/info?firstname=Ma&lastname=tuto&age=24
    r.GET("/user/info", func(c *gin.Context) {
        firstname := c.DefaultQuery("firstname", "Ma")
        lastname := c.Query("lastname")
        age := c.Request.URL.Query().Get("age")
        c.String(http.StatusOK, "%s %s 今年%s岁了!", firstname, lastname, age)
    })
    r.Run()
}
```

### 获取post参数

```go
package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    //获取post参数 - form-data
    r.POST("/user/add", func(c *gin.Context) {
        username := c.PostForm("username")
        password := c.DefaultPostForm("password", "123456") // 此方法可以设置默认值
        c.JSON(http.StatusOK, gin.H{
            "status": "posted", 
            "message": "请求成功", 
            "username": username, 
            "password": password,
        })
    })
    r.Run()
}
```

### Get + Post 混合

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // 示例：
    // POST /post?pagesize=10&pageindex=1 HTTP/1.1
    // Content-Type: application/x-www-form-urlencoded
    // name=admin&age=18
    r.POST("/user/list", func(c *gin.Context) {
        pagesize := c.DefaultQuery("pagesize", "10")
        pageindex := c.Query("pageindex")
        name := c.PostForm("name")
        age := c.PostForm("age")
        c.String(http.StatusOK, "name:%s,age:%s,pagesize:%s,pageindex:%s", name, age, pagesize, pageindex)
    })
    r.Run()
}
```

### 文件上传

```go
package main

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    //单文件上传
    //限制上传大小 默认 32MiB
    r.MaxMultipartMemory = 8 << 20 //8MiB
    r.POST("/upload", func(c *gin.Context) {
        file, _ := c.FormFile("file")
        log.Println(file.Filename)
        // 上传文件到指定的路径
        // c.SaveUploadedFile(file, dst)

        c.String(http.StatusOK, "%s is uploaded!", file.Filename)
    })
    r.Run()
}
```

### 多文件上传

```go
package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    //限制上传大小 默认 32MiB
    r.MaxMultipartMemory = 8 << 20 //8MiB
    // 多文件上传
    r.POST("/uploads", func(c *gin.Context) {
        form, _ := c.MultipartForm()
        files := form.File["upload"]
        for _, file := range files {
            log.Println(file.Filename)

            // 上传文件到指定的路径
            // c.SaveUploadedFile(file, dst)
        }
        c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
    })
    r.Run()
}
```
