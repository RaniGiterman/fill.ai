package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Object struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IMG         string `json:"img"`
}

const URL = "https://www.benda.co.il/product/99400-016-60/"

func main() {
	fmt.Println("Starting!")
	fmt.Println("Querying HTML...")
	html, err := queryHTML(URL)
	if err != nil {
		log.Fatalf("error querying page: %v", err)
	}

	fmt.Println("Queried HTML successfully!")

	obj := Object{Title: "", Description: "", IMG: ""}
	json, err := json.Marshal(obj)
	if err != nil {
		log.Fatalf("error marshaling json: %v", err)
	}

	fmt.Println("Querying LLM...")
	res, err := queryLLM(string(json), html)
	if err != nil {
		log.Fatalf("error querying LLM: %v", err)
	}

	fmt.Println("Queried LLM successfully!")

	fmt.Println(res)
}
