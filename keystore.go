package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

//var keyStoreLock sync.Mutex

func init() {
	InitializeKeystore()
}

func get(c *gin.Context, ch <-chan *[]string) {
	req := <-ch
	cmd := *req

	keystore.mu.Lock()
	defer keystore.mu.Unlock()

	val, exists := keystore.keys[cmd[1]]

	fmt.Println(exists)
	if exists {
		c.JSON(http.StatusOK, gin.H{
			"value": val.Value,
		})
		return
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "key not found"})

	//fmt.Println("updated", keystore.keys)

}

func set(c *gin.Context, ch <-chan *[]string) {
	req := <-ch
	cmd := *req

	keystore.mu.Lock()
	defer keystore.mu.Unlock()

	_, exists := keystore.keys[cmd[1]]

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

	keystore.keys[cmd[1]] = dummy

	//fmt.Println(keystore.keys)
}

func deleteExpiredKeys() {
	for {
		time.Sleep(10 * time.Second)

		currentTime := time.Now()

		keystore.mu.Lock()

		for key, item := range keystore.keys {
			if item.ExTime != "" {
				exTimeValue, err := strconv.Atoi(item.ExTime)
				if err != nil {
					continue
				}

				timeDifference := currentTime.Sub(item.InsertTime)
				if int(timeDifference.Seconds()) >= exTimeValue {
					delete(keystore.keys, key)
				}
			}
		}

		keystore.mu.Unlock()

	}
}

func containsStringElement(arr []string, target string) bool {
	for _, element := range arr {
		if element == target {
			return true
		}
	}
	return false
}
