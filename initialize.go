package main

import "time"

var keystore map[string]KeyModel
var liststore map[string]ListModel

func InitializeKeystore() {

	keystore = map[string]KeyModel{
		"demo_key": {
			Value:      "42",
			ExTime:     "30",
			InsertTime: time.Now(),
		},
	}
}

func InitializeListStore() {
	liststore = map[string]ListModel{
		"demo_list": {
			Value: []string{"1", "2", "3"},
		},
	}
}
