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

var comm1 chan *[]string
var comm2 chan *[]string
var comm3 chan *[]string
var comm4 chan *[]string

func mux(c *gin.Context) {

	var action Command

	if err := c.ShouldBindJSON(&action); err != nil {
		fmt.Println(err)
		return
	}

	list := strings.Split(action.Command, " ")

	go get(c, comm1)
	go set(c, comm2)
	go qpush(c, comm3)
	go qpop(c, comm4)
	
	go deleteExpiredKeys()

	switch list[0] {
	case "GET":
		comm1 <- &list
	case "SET":
		comm2 <- &list
	case "QPUSH":
		comm3 <- &list
	case "QPOP":
		comm4 <- &list
	default:
		c.JSON(400, gin.H{"message": "Invalid Command or Invalid Command"})
	}
}

func main() {

	comm2 = make(chan *[]string)
	comm3 = make(chan *[]string)
	comm1 = make(chan *[]string)
	comm4 = make(chan *[]string)

	router := gin.Default()

	router.GET("/", mux)

	if err := router.Run(":3000"); err != nil {
		fmt.Println(err)
	}

}
