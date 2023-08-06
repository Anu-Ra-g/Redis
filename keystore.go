package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var keyStoreLock sync.Mutex

func init() {
	InitializeKeystore()
}

func containsStringElement(arr []string, target string) bool {
	for _, element := range arr {
		if element == target {
			return true
		}
	}
	return false
}

func get(c *gin.Context, ch *[]string) {
	keyStoreLock.Lock()
	defer keyStoreLock.Unlock()

	cmd := *ch

	if val, exists := keystore[cmd[1]]; exists {
		c.JSON(http.StatusOK, gin.H{
			"value": val.Value,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Key does not exits",
		})
	}

	fmt.Println(keystore)

}

func set(c *gin.Context, ch *[]string) {
	keyStoreLock.Lock()
	defer keyStoreLock.Unlock()

	cmd := *ch

	_, exists := keystore[cmd[1]]

	dummy := KeyModel{
		Value: cmd[2],
	}

	match1 := containsStringElement(cmd, "NX")
	match2 := containsStringElement(cmd, "XX")
	match3 := containsStringElement(cmd, "EX")

	if (!exists && match1) || (exists && match2) {
		exists = true
	}

	if match3 && exists {
		dummy.ExTime = cmd[4]
		dummy.InsertTime = time.Now()
	}

	keystore[cmd[1]] = dummy

	c.JSON(http.StatusOK, gin.H{
		"updated_value": dummy,
	})
}

func deleteExpiredKeys() {
	for {
		time.Sleep(10 * time.Second)

		keyStoreLock.Lock()

		currentTime := time.Now()

		for key, item := range keystore {
			if item.ExTime != "" {
				exTimeValue, err := strconv.Atoi(item.ExTime)
				if err != nil {
					continue
				}

				timeDifference := currentTime.Sub(item.InsertTime)
				if int(timeDifference.Seconds()) >= exTimeValue {
					delete(keystore, key)
				}
			}
		}

		keyStoreLock.Unlock()
	}
}
