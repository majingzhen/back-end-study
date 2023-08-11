# Go - Gin进阶

## 路由分组

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
    sys := r.Group("/sys")
    //限制上传大小 默认 32MiB
    r.MaxMultipartMemory = 8 << 20 //8MiB
    {
        sys.GET("/ping", func(c *gin.Context) {
            c.JSON(200, gin.H{
                "message": "pong",
            })
        })
        //单文件上传
        sys.POST("/upload", func(c *gin.Context) {
            file, _ := c.FormFile("file")
            log.Println(file.Filename)
            // 上传文件到指定的路径
            // c.SaveUploadedFile(file, dst)

            c.String(http.StatusOK, "%s is uploaded!", file.Filename)
        })

        // 多文件上传
        sys.POST("/uploads", func(c *gin.Context) {
            form, _ := c.MultipartForm()
            files := form.File["upload"]
            for _, file := range files {
                log.Println(file.Filename)

                // 上传文件到指定的路径
                // c.SaveUploadedFile(file, dst)
            }
            c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
        })
    }

    user := r.Group("/user")
    {
        // 此规则能够匹配/user/matuto 这种格式，但不能匹配/user/ 或 /user这种格式
        user.GET("/:name", func(c *gin.Context) {
            name := c.Param("name")
            c.String(http.StatusOK, "Hello %s", name)
        })

        // 但是，这个规则既能匹配/user/matuto/格式也能匹配/user/matuto/send这种格式
        // 如果没有其他路由器匹配/user/matuto，它将重定向到/user/matuto/
        user.GET("/:name/*action", func(c *gin.Context) {
            name := c.Param("name")
            action := c.Param("action")
            message := name + " is " + action
            c.String(http.StatusOK, message)
        })

        //获取get参数 /user/info?firstname=Ma&lastname=tuto&age=24
        user.GET("/info", func(c *gin.Context) {
            firstname := c.DefaultQuery("firstname", "Ma")
            lastname := c.Query("lastname")
            age := c.Request.URL.Query().Get("age")
            c.String(http.StatusOK, "%s %s 今年%s岁了!", firstname, lastname, age)
        })

        //获取post参数
        user.POST("/add", func(c *gin.Context) {
            username := c.PostForm("username")
            password := c.DefaultPostForm("password", "123456") // 此方法可以设置默认值
            c.JSON(http.StatusOK, gin.H{
                "status": "posted", "message": "请求成功", "username": username, "password": password,
            })
        })

        // 示例：
        // POST /post?pagesize=10&pageindex=1 HTTP/1.1
        // Content-Type: application/x-www-form-urlencoded
        // name=admin&age=18
        user.POST("/list", func(c *gin.Context) {
            pagesize := c.DefaultQuery("pagesize", "10")
            pageindex := c.Query("pageindex")
            name := c.PostForm("name")
            age := c.PostForm("age")
            c.String(http.StatusOK, "name:%s,age:%s,pagesize:%s,pageindex:%s", name, age, pagesize, pageindex)
        })
    }
    r.Run()
}
```

## 无中间件启动

```go
package main

import "github.com/gin-gonic/gin"

func main() {
     //默认启动方式，包含 Logger、Recovery 中间件
    //r := gin.Default()
    r := gin.New()
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    r.Run() // listen and serve on 0.0.0.0:8080
}
```

## 使用中间件

```go
package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    // 创建一个不包含中间件的路由器
    r := gin.New()
    // 全局中间件
    //使用Logger中间件
    r.Use(gin.Logger())
    //使用Recovery中间件
    r.Use(gin.Recovery())

    // 为路由添加中间件
    //    r.GET("/hello", func(c *gin.Context) {
    //        c.JSON(200, gin.H{
    //            "message": "go",
    //        })
    //    }, AuthRequired())

    // 为路由组添加中间件
    // auth := r.Group("/", AuthRequired())
    auth := r.Group("/")
    auth.Use(AuthRequired())
    {

    }
    r.run()
}
```

## 写日志文件

```go
package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "matuto.cc/go-gin-study/filters"
    "github.com/gin-gonic/gin"
)

func main() {
    // 创建一个不包含中间件的路由器
    r := gin.New()

    // 禁用控制台颜色
    gin.DisableConsoleColor()
    //创建日志文件
    f, _ := os.Create("gin.log")
    gin.DefaultWriter = io.MultiWriter(f)
    // 如果需要将日志同时写入文件和控制台，请使用以下代码
    // gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
    // 全局中间件
    //使用Logger中间件
    r.Use(gin.Logger())
    //使用Recovery中间件
    r.Use(gin.Recovery())

    r.GET("/ping", func(c *gin.Conte
xt) {
        c.String(200, "pong")
    })
    r.run()
}
```

# 

## 创建router组件

这里我们将 router 的创建 及中间件的使用都抽离出来一个单独的文件中，使用时直接进行调用就可以了

```go
package routers

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "log"
    "matuto.cc/go-gin-study/filters"
    "net/http"
)
type Routers struct {

}
func (routers *Routers) InitRouter() *gin.Engine {
    r := gin.New()
    //使用自定义Logger中间件
    r.Use(filters.LogFilter())
    //使用Recovery中间件
    r.Use(gin.Recovery())
    r.GET("/hello", func(c *gin.Context) {
        c.String(200, "go")
    })
    sys := r.Group("/sys")
    return r
}
```

```go
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
    //启动服务
    r.Run()
}
```

## 模型绑定和验证

这里的解释直接借用官方的解释

若要将请求主体绑定到结构体中，请使用模型绑定，目前支持JSON、XML、YAML和标准表单值(foo=bar&boo=baz)的绑定。

Gin使用 [go-playground/validator.v8](https://github.com/go-playground/validator) 验证参数，[查看完整文档](https://godoc.org/gopkg.in/go-playground/validator.v8#hdr-Baked_In_Validators_and_Tags)。

需要在绑定的字段上设置tag，比如，绑定格式为json，需要这样设置 `json:"fieldname"` 。

此外，Gin还提供了两套绑定方法：

- Must bind
- - Methods - `Bind`, `BindJSON`, `BindXML`, `BindQuery`, `BindYAML`
- - Behavior - 这些方法底层使用 `MustBindWith`，如果存在绑定错误，请求将被以下指令中止 `c.AbortWithError(400, err).SetType(ErrorTypeBind)`，响应状态代码会被设置为400，请求头`Content-Type`被设置为`text/plain; charset=utf-8`。注意，如果你试图在此之后设置响应代码，将会发出一个警告 `[GIN-debug] [WARNING] Headers were already written. Wanted to override status code 400 with 422`，如果你希望更好地控制行为，请使用`ShouldBind`相关的方法
- Should bind
- - Methods - `ShouldBind`, `ShouldBindJSON`, `ShouldBindXML`, `ShouldBindQuery`, `ShouldBindYAML`
- - Behavior - 这些方法底层使用 `ShouldBindWith`，如果存在绑定错误，则返回错误，开发人员可以正确处理请求和错误。

当我们使用绑定方法时，Gin会根据Content-Type推断出使用哪种绑定器，如果你确定你绑定的是什么，你可以使用`MustBindWith`或者`BindingWith`。

你还可以给字段指定特定规则的修饰符，如果一个字段用`binding:"required"`修饰，并且在绑定时该字段的值为空，那么将返回一个错误。

```go
type Login struct {
    User     string `form:"user" json:"user" xml:"user" binding:"required"`
    Password string `form:"password" json:"password" xml:"password" binding:"required"`
}
```

```go
func loginJSON(c *gin.Context) {
    var json Login
    if err := c.ShouldBindJSON(&json); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if json.User != "admin" || json.Password != "123456" {
        c.JSON(http.StatusUnauthorized, gin.H{"stauts": "unauthorized"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"stauts": "you are logged in"})
}

func loginXML(c *gin.Context) {
    var xml Login
    if err := c.ShouldBindXML(&xml); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if xml.User != "admin" || xml.Password != "123456" {
        c.JSON(http.StatusUnauthorized, gin.H{"stauts": "unauthorized"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"stauts": "you are logged in"})
}

func loginForm(c *gin.Context) {
    var form Login
    if err := c.ShouldBind(&form); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if form.User != "admin" || form.Password != "123456" {
        c.JSON(http.StatusUnauthorized, gin.H{"stauts": "unauthorized"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"stauts": "you are logged in"})
}
```

## 自定义校验器

```go
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

    // 注册验证器
    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        v.RegisterValidation("bookabledate", bookableDate)
    }
    r.GET("/bookable", getBookable)
    //启动服务
    r.Run()
}

type Booking struct {
    CheckIn  time.Time `form:"check_in" binding:"required,bookabledate" time_format:"2006-01-02"`
    CheckOut time.Time `form:"check_out" binding:"required,gtfield=CheckIn" time_format:"2006-01-02"`
}

func bookableDate(fl validator.FieldLevel) bool {
    if date, ok := fl.Field().Interface().(time.Time); ok {
        today := time.Now()
        if today.Year() > date.Year() || today.YearDay() > date.YearDay() {
            return false
        }
    }
    return true
}

func getBookable(c *gin.Context) {
    var b Booking
    if err := c.ShouldBindWith(&b, binding.Query); err == nil {
        c.JSON(http.StatusOK, gin.H{"message": "Booking dates are valid!"})
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    }
}
```

## 只绑定get参数

```go
package main

import (
    "io"
    "log"
    "net/http"
    "os"

    "matuto.cc/go-gin-study/routers"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    r.Any("/testGet", StartPage)
    //启动服务
    r.Run()
}

type Person struct {
    Name    string `form:"name"`
    Address string `form:"address"`
}

func StartPage(c *gin.Context) {
    var person Person
    if c.ShouldBindQuery(&person) == nil {
        log.Println("====== Only Bind By Query String ======")
        log.Println(person.Name)
        log.Println(person.Address)
    }
    c.String(http.StatusOK, "Success")
}
```

## 绑定Get参数或者Post参数

```go
package main

import (
    "io"
    "log"
    "net/http"
    "os"

    "matuto.cc/go-gin-study/routers"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    r.Any("/testGet", StartPage)
    //启动服务
    r.Run()
}

type Person struct {
    Name    string `form:"name" json:"name"`
    Address string `form:"address" json:"address"`
}

func StartPage(c *gin.Context) {
    var person Person
    var personJson Person
    if c.Bind(&person) == nil {
        log.Println("====== Only Bind By Query String ======")
        log.Println(person.Name)
        log.Println(person.Address)
    }
    if c.BindJSON(&personJson) == nil {
        log.Println("====== Only Bind By JSON ======")
        log.Println(personJson.Name)
        log.Println(personJson.Address)
    }
    c.String(http.StatusOK, "Success")
}
```

## 绑定uri

```go
package main

import (
    "io"
    "net/http"
    "os"

    "matuto.cc/go-gin-study/routers"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    r.GET("/:name/:id", BindUri)
    //启动服务
    r.Run()
}

type Student struct {
    ID   string `uri:"id" binding:"required,uuid"`
    Name string `uri:"name" binding:"required"`
}

func BindUri(c *gin.Context) {
    var student Student
    if err := c.ShouldBindUri(&student); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"name": student.Name, "uuid": student.ID})
}
```

## 绑定表单复选框

```go
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
    r.POST("/subColors", SubColors)
    //启动服务
    r.Run()
}

type MyColors struct {
    Colors []string `form:"colors[]"`
}

func SubColors(c *gin.Context) {
    var mc MyColors
    c.ShouldBind(&mc)
    c.JSON(http.StatusOK, gin.H{"color": mc.Colors})
}
```

```html
<form action="http://127.0.0.1:8080/subColors" method="POST">
    <p>Check some colors</p>
    <label for="red">Red</label>
    <input type="checkbox" name="colors[]" value="red" id="red">
    <label for="green">Green</label>
    <input type="checkbox" name="colors[]" value="green" id="green">
    <label for="blue">Blue</label>
    <input type="checkbox" name="colors[]" value="blue" id="blue">
    <input type="submit">
</form>
```
