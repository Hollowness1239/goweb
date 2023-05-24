package main

import (
	"encoding/json"
	"fmt"
)

type emp struct {
	ID      int
	EmpName string
	Tel     string
	Email   string
}

func main() {
	data, _ := json.Marshal(&emp{1, "N", "0988899", "as@Email.com"})
	fmt.Println(string(data))
}
