# Go - Gin 使用自定义验证器 V8 V10

今天在学习 gin 的自定义验证器时，idea 报错，无法注册验证器，经排查后发现是版本问题导致。

V10:

```go
package main

import (
    "io"
    "net/http"
    "os"
    "time"

    "github.com/gin-gonic/gin/binding"

    "github.com/go-playground/validator/v10"

    "matuto.cc/go-gin-study/routers"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

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

V8:

```go
package main

import (
    "io"
    "net/http"
    "os"
    "reflect"
    "time"

    "github.com/gin-gonic/gin/binding"

    "github.com/go-playground/validator/v10"

    "matuto.cc/go-gin-study/routers"

    "github.com/gin-gonic/gin"
)

func main() {


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

func bookableDate(
    v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
    field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
    if date, ok := field.Interface().(time.Time); ok {
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
