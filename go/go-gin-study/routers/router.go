package routers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"matuto.cc/go-gin-study/controller"
	"matuto.cc/go-gin-study/validators"

	"github.com/gin-gonic/gin"
	"matuto.cc/go-gin-study/filters"
)

type Routers struct {
}

var secrets = gin.H{
	"zhangsan": gin.H{"account": "zhangsan", "email": "zhangsan@qq.com"},
	"lisi":     gin.H{"account": "lisi", "email": "lisi@qq.com"},
	"wangwu":   gin.H{"account": "wangwu", "email": "wangwu@qq.com"},
}

func (routers *Routers) InitRouter() *gin.Engine {
	r := gin.New()
	//使用自定义Logger中间件
	r.Use(filters.LogFilter())
	//使用Recovery中间件
	r.Use(gin.Recovery())
	r.LoadHTMLGlob("templates/**/*")
	// 只能展示文件，例如图片等
	r.Static("/assets", "./assets")
	//展示目录+文件
	r.StaticFS("/more_static", gin.Dir("my_file_system", true))
	// 静态资源文件
	r.StaticFile("/html.png", "./resources/html.png")

	r.GET("/long_async", func(c *gin.Context) {
		// 创建要在goroutine中 使用的副本
		cCp := c.Copy()
		go func() {
			time.Sleep(5 * time.Second)
			// 这里使用创建的副本
			log.Println("Done! in path " + cCp.Request.URL.Path)
		}()
	})
	r.GET("/long_sync", func(c *gin.Context) {
		// simulate a long task with time.Sleep(). 5 seconds
		time.Sleep(5 * time.Second)

		// 这里没有使用goroutine，所以不用使用副本
		log.Println("Done! in path " + c.Request.URL.Path)
	})
	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"zhangsan": "123456",
		"lisi":     "123456",
		"wangwu":   "123456",
	}))
	authorized.GET("/secrets", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		if secret, ok := secrets[user]; ok {
			c.JSON(http.StatusOK, gin.H{
				"user":   user,
				"secret": secret,
			})
		} else {
			c.JSON(200, gin.H{
				"user":   user,
				"secret": "No Secret",
			})
		}
	})

	// 重定向外部链接
	r.GET("/testRedirect", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://www.baidu.com")
	})
	// 路由重定向
	r.GET("/test1", func(c *gin.Context) {
		c.Request.URL.Path = "/test2"
		r.HandleContext(c)
	})
	r.GET("/test2", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"hello": "world"})
	})

	r.Delims("{[{", "}]}")
	r.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})
	r.LoadHTMLFiles("./templates/raw/raw.tmpl")
	r.GET("/raw", func(c *gin.Context) {
		c.HTML(http.StatusOK, "raw.tmpl", map[string]interface{}{"now": time.Date(2023, 06, 05, 0, 0, 0, 0, time.UTC)})
	})

	sys := r.Group("/sys")
	{
		sys.POST("/loginJSON", controller.LoginJSON)
		sys.POST("/loginXML", controller.LoginXML)
		sys.POST("/loginForm", controller.LoginForm)
	}
	r.GET("index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{"title": "a"})
	})
	r.GET("/a/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "a/index.tmpl", gin.H{"title": "a"})
	})
	r.GET("/b/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "b/index.tmpl", gin.H{"title": "b"})
	})
	// 注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("bookabledate", validators.BookableDate)
	}
	r.GET("/bookable", controller.GetBookable)

	r.Any("/testGet", controller.StartPage)

	r.GET("/:name/:id", controller.BindUri)

	r.POST("/subColors", controller.SubColors)

	render := r.Group("render")
	{
		render.GET("/someJson", controller.SomeJson)
		render.GET("/moreJson", controller.MoreJson)
		render.GET("/someXml", controller.SomeXml)
		render.GET("/someYaml", controller.SomeYaml)
		render.GET("/someProtoBuf", controller.SomeProtoBuf)
	}

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
	r.GET("/cookie", func(c *gin.Context) {
		cookie, err := c.Cookie("gin_cookie")
		if err != nil {
			cookie = "notSet"
			c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
		}
		log.Printf("Cookie value: %s \n", cookie)
	})

	r.Use(MiddlewareOne())
	r.Use(MiddlewareTwo())
	r.Use(MiddlewareThree())

	r.GET("/hello", func(c *gin.Context) {
		c.String(200, "go")
	})

	return r
}

func MiddlewareOne() func(c *gin.Context) {
	return func(c *gin.Context) {
		log.Println("MiddlewareOne Start")
		c.Next()
		log.Println("MiddlewareOne End")
	}
}

func MiddlewareTwo() func(c *gin.Context) {
	return func(c *gin.Context) {
		log.Println("MiddlewareTwo Start")
		c.Next()
		log.Println("MiddlewareTwo End")
	}
}
func MiddlewareThree() func(c *gin.Context) {
	return func(c *gin.Context) {
		log.Println("MiddlewareThree Start")
		c.Next()
		log.Println("MiddlewareThree End")
	}
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d/%02d/%02d", year, month, day)
}
