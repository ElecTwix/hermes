package genai

import (
	"context"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GenAI struct {
	model  *genai.GenerativeModel
	client *genai.Client
}

func NewGenAI() *GenAI {
	// Access your API key as an environment variable (see "Set up your API key" above)

	return &GenAI{}
}

func (g *GenAI) Login(ctx context.Context, apiKey string, modelName string) error {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	model := client.GenerativeModel(modelName)

	g.client = client
	g.model = model

	return nil
}

func (g *GenAI) Generate(prompt string, ctx context.Context) (string, error) {
	resp, err := g.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	var parts []genai.Part
	for _, candidate := range resp.Candidates {
		parts = append(parts, candidate.Content.Parts...)
	}

	str := ""
	for _, part := range parts {
		switch p := part.(type) {
		case genai.Text:
			str += string(p)
		case genai.Blob:
			str += fmt.Sprintf("Blob: %s", p.MIMEType)
		default:
			str += fmt.Sprintf("Unknown type: %T", p)
		}
	}

	return str, nil
}
