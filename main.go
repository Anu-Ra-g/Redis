package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func init() {
	InitializeKeystore()
	InitializeListStore()
}

//var comm1 = make(chan *[]string)
//var comm2 = make(chan *[]string)
//var comm3 = make(chan *[]string)
//var comm4 = make(chan *[]string)

func mux(c *gin.Context) {

	var action Command

	if err := c.BindJSON(&action); err != nil {
		fmt.Println(err)
		return
	}
	
	list := strings.Split(action.Command, " ")

	switch list[0] {
	case "GET":
		get(c, &list)
	case "SET":
		set(c, &list)
	case "QPUSH":
		qpush(c, &list)
	case "QPOP":
		qpop(c, &list)
	default:
		c.JSON(400, gin.H{"message": "Invalid Command or Invalid Command"})
	}

}

func main() {

	router := gin.Default()

	go deleteExpiredKeys()

	router.POST("/", mux)

	if err := router.Run(":3000"); err != nil {
		fmt.Println(err)
	}

}
