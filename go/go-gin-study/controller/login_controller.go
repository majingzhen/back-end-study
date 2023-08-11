package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Login struct {
	User     string `form:"user" json:"user" xml:"user" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

func LoginJSON(c *gin.Context) {
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

func LoginXML(c *gin.Context) {
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

func LoginForm(c *gin.Context) {
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
