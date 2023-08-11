package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

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
