package main

import (
	"context"
	"fmt"
	"log"

	"github.com/openai/openai-go/v3"
)

type Object struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IMG         string `json:"img"`
}

const SYSTEM_PROMPT = `
You are a precise data extraction assistant. You will receive an HTML page that decribes a consumer product with details and must extract information to fill a JSON object with exactly three fields.

## Input Format
You will receive:
1. An empty JSON template: {"title": "", "description": "", "img": ""}
2. A complete HTML page

## Field Descriptions

*title* (string):
- Extract the product title
- Return only the text content, no HTML tags
- If multiple candidates exist, choose the most descriptive and prominent one

*description* (string):
- Extract product description from the page content

*img* (string):
- Extract the product image. if more than one, then concat them all with comma

## Output Requirements
- Return ONLY the filled JSON object
- Ensure valid JSON syntax
- Use empty string "" for any field that cannot be determined

Now process the provided HTML and return the filled JSON.
`

const URL = "https://www.benda.co.il/product/36208-498-32/"

func main() {
	fmt.Println("Starting!")
	fmt.Println("Querying HTML...")
	html, err := queryHTML(URL)
	if err != nil {
		log.Fatalf("error querying page: %v", err)
	}

	fmt.Println("Queried HTML successfully!")

	// defaults to os.LookupEnv("OPENAI_API_KEY")
	client := openai.NewClient()

	ctx := context.Background()

	fmt.Println("Querying LLM...")
	res, err := queryGPT(ctx, client, SYSTEM_PROMPT, html)
	if err != nil {
		log.Fatalf("error querying LLM: %v", err)
	}
	fmt.Println("Queried LLM successfully!")

	fmt.Println(res)
}
