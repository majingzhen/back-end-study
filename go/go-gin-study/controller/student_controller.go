package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

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
