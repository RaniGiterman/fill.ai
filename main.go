package main

import (
	"fmt"
	"log"
	"net/http"

	"fill.ai/controller"
)

// const URL = "https://www.benda.co.il/product/36208-498-32/"
// const URL = "https://bconnect.co.il/product/%d7%9e%d7%98%d7%a2%d7%9f-%d7%9e%d7%94%d7%99%d7%a8-%d7%a0%d7%99%d7%99%d7%93-%d7%9c-boost%e2%86%91charge-pro-apple-watch/"

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// serve static files
	http.HandleFunc("/fill", controller.Fill)

	fmt.Println("Server starting on :8080...")
	// Listen and serve on port 8080
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
