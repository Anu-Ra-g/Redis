package main

import (
	"time"
)

type Command struct {
	Command string `json:"command"`
}

type KeyModel struct {
	Value      string
	ExTime     string
	InsertTime time.Time
}

type ListModel struct {
	Value []string
}
