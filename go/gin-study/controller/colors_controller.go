package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type MyColors struct {
	Colors []string `form:"colors[]"`
}

func SubColors(c *gin.Context) {
	var mc MyColors
	c.ShouldBind(&mc)
	c.JSON(http.StatusOK, gin.H{"color": mc.Colors})
}
