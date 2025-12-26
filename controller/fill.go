package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	U "fill.ai/util"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type Body struct {
	URL string `json:"url"`
}

func Fill(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// if not post, fuckoff
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"status": "fail", "msg": "Method not allowed"})
		return
	}

	// 2. Decode the JSON body
	var body Body
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"status": "fail", "msg": "failed to parse JSON body"})
		return
	}

	if body.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"status": "fail", "msg": "URL not provided"})
		return
	}

	fmt.Println("Starting!")
	fmt.Println("Querying HTML...")
	html, err := U.QueryHTML(body.URL)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"status": "fail", "msg": "error querying page"})
		return
	}

	fmt.Println("Queried HTML successfully!")

	// defaults to os.LookupEnv("OPENAI_API_KEY")
	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
	)
	ctx := context.Background()

	fmt.Println("Querying LLM...")
	res, err := U.QueryGPT(ctx, client, html, body.URL)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"status": "fail", "msg": "error querying GPT"})
		return
	}

	fmt.Println("Queried LLM successfully!")
	fmt.Println(res)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "msg": res})
}
