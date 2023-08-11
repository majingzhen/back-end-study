package controller

import (
	"github.com/gin-gonic/gin/testdata/protoexample"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
