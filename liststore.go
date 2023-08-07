package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func init() {
	InitializeListStore()
}

func qpush(c *gin.Context, ch <-chan *[]string) {
	req := <-ch
	cmd := *req

	key := cmd[1]

	val, exists := liststore[key]

	if !exists {
		liststore[key] = ListModel{
			Value: cmd[2:],
		}
	} else {
		updatedArr := append(val.Value, cmd[2:]...)
		val := ListModel{
			Value: updatedArr,
		}
		liststore[key] = val
	}

	//fmt.Println(liststore)
}

func qpop(c *gin.Context, ch <-chan *[]string) {
	req := <-ch
	cmd := *req

	key := cmd[1]

	val, exists := liststore[key]

	popped_value := ""

	if exists && len(val.Value) > 0 {
		index := len(val.Value) - 1
		popped_value = val.Value[index]
		newArr := val.Value[:index]

		liststore[key] = ListModel{
			Value: newArr,
		}

		c.JSON(http.StatusOK, gin.H{
			"value": popped_value,
		})
		return
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "queue is empty",
	})

	//fmt.Println(liststore)
}
