package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

func init() {
	InitializeListStore()
}

var listStoreLock sync.Mutex

func qpush(c *gin.Context, ch *[]string) {
	listStoreLock.Lock()
	defer listStoreLock.Unlock()

	cmd := *ch

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

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})

	fmt.Println(liststore)
}

func qpop(c *gin.Context, ch *[]string) {
	listStoreLock.Lock()
	defer listStoreLock.Unlock()

	cmd := *ch

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
			"message": popped_value,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "queue is empty",
		})
	}

	fmt.Println(liststore)
}
