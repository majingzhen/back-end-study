# Go - Gin 模板渲染

## XML、JSON、YAML和ProtoBuf 渲染（输出格式）

```go
func SomeJson(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "OK", "status": http.StatusOK})
}

func MoreJson(c *gin.Context) {
    var msg struct {
        Name    string `json:"user"`
        Message string
        Number  int
    }
    msg.Name = "Matuto"
    msg.Message = "hello"
    msg.Number = 25
    c.JSON(http.StatusOK, msg)
}

func SomeXml(c *gin.Context) {
    c.XML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
}

func SomeYaml(c *gin.Context) {
    c.YAML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
}

func SomeProtoBuf(c *gin.Context) {
    reps := []int64{int64(1), int64(2)}
    label := "test"
    data := &protoexample.Test{
        Label: &label,
        Reps:  reps,
    }
    c.ProtoBuf(http.StatusOK, data)
}
```

```go
render := r.Group("render")
    {
        render.GET("/someJson", controller.SomeJson)
        render.GET("/moreJson", controller.MoreJson)
        render.GET("/someXml", controller.SomeXml)
        render.GET("/someYaml", controller.SomeYaml)
        render.GET("/someProtoBuf", controller.SomeProtoBuf)
    }
```

## 设置静态文件路径

```go
func main() {
    r := gin.Default()

    // 只能展示文件，例如图片等
    r.Static("/assets", "./assets")
    //展示目录+文件
    r.StaticFS("/more_static", gin.Dir("my_file_system", true))
    // 静态资源文件
    r.StaticFile("/html.png", "./resources/html.png")
    //启动服务
    r.Run()
}
```

## HTML 渲染

```go
func main() {
    r := gin.Default()
    // 加载html路径
    // 使用LoadHTMLGlob() 或者 LoadHTMLFiles()
    r.LoadHTMLGlob("templates/*")
    r.GET("/index", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", gin.H{"title": "Main website"})
    })

    //启动服务
    r.Run()
}
```

```html
<html>
    <h1>
        {{ .title }}
    </h1>
</html>
```

### 在不同目录中使用具有相同名称的模板

```go
// 不同目录下具有相同名称的模板
    r.LoadHTMLGlob("templates/**/*")
    r.GET("/a/index", func(c *gin.Context) {
        c.HTML(http.StatusOK, "a/index.tmpl", gin.H{"title": "a"})
    })
    r.GET("/b/index", func(c *gin.Context) {
        c.HTML(http.StatusOK, "b/index.tmpl", gin.H{"title": "b"})
    })
```

```html
{{ define "b/index.tmpl" }}
<html>
    <h1>
        {{ .title }}
    </h1>
<p>Using b/index.tmpl</p>
</html>
{{ end }}
```

```html
{{ define "a/index.tmpl" }}
<html>
    <h1>
        {{ .title }}
    </h1>
<p>Using a/index.tmpl</p>
</html>
{{ end }}
```

### 自定义函数模板

```go
r.Delims("{[{", "}]}")
    r.SetFuncMap(template.FuncMap{
        "formatAsDate": formatAsDate,
    })
    // 不同目录下具有相同名称的模板
    //r.LoadHTMLGlob("./templates/*")
    r.LoadHTMLFiles("./templates/raw/raw.tmpl")
    r.GET("/raw", func(c *gin.Context) {
        c.HTML(http.StatusOK, "raw.tmpl", map[string]interface{}{"now": time.Date(2023, 06, 05, 0, 0, 0, 0, time.UTC)})
    })
```

```go
func formatAsDate(t time.Time) string {
    year, month, day := t.Date()
    return fmt.Sprintf("%d/%02d/%02d", year, month, day)
}
```

```html
<html>
    <h1>
        {[{ .now | formatAsDate }]}
    </h1>
</html>
```

## 重定向

```go
// 重定向外部链接
r.GET("/testRedirect", func(c *gin.Context) {
    c.Redirect(http.StatusMovedPermanently, "https://www.baidu.com")
})
```

```go
// 路由重定向
    r.GET("/test1", func(c *gin.Context) {
        c.Request.URL.Path = "/test2"
        r.HandleContext(c)
    })
    r.GET("/test2", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"hello": "world"})
    })
```
