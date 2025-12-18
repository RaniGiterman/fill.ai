package main

import (
	"context"

	"github.com/openai/openai-go/v3"
)

func queryGPT(ctx context.Context, client openai.Client, prompt, html string) (string, error) {
	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
			openai.UserMessage(html),
		},
		Seed:  openai.Int(0),
		Model: openai.ChatModelGPT5_1,
	}

	// completion, err := client.Chat.Completions.New(ctx, params)
	completion, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		return "", err
	}

	return completion.Choices[0].Message.Content, nil
}

// resp, err := client.Responses.New(ctx, responses.ResponseNewParams{
// 	Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(prompt)},
// 	Model: openai.ChatModelGPT4,
// })
// if err != nil {
// 	return "", err
// }
//
// // println(resp.OutputText())
// return resp.OutputText(), nil
