package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type employee struct {
	ID      int
	EmpName string
	Tel     string
	Email   string
}

func main() {
	e := employee{}
	err := json.Unmarshal([]byte(`{"ID":101,"EmpName":"thuu","Tel":"099999888","Email":"E@gmial.com"}`), &e)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(e)
	fmt.Println(e.Email)
}
