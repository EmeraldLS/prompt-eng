package llm_sync

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/mistral"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
)

func UseLlama(prompt string, chunkChan chan<- []byte) {
	llm, err := ollama.New(ollama.WithModel("llama3"), ollama.WithServerURL(""))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	completion, err := llm.Call(ctx, prompt, llms.WithTemperature(0.8), llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		fmt.Println("Gotten: ", string(chunk))
		chunkChan <- chunk
		return nil
	}))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Completion: ", completion)

}

func UseMistral(prompt string, chunkChan chan<- []byte) {
	llm, err := mistral.New(mistral.WithModel("open-mistral-7b"))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	completion, err := llm.Call(ctx, prompt, llms.WithTemperature(0.8), llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		chunkChan <- chunk
		return nil
	}))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Completion: ", completion)
}

func UseNvidia(prompt string, chunkChan chan<- []byte) {
	key := os.Getenv("NVIDIA_API_KEY")
	llm, err := openai.New(
		openai.WithBaseURL("https://integrate.api.nvidia.com/v1/"),
		openai.WithModel("mistralai/mixtral-8x7b-instruct-v0.1"),
		openai.WithToken(key),
		// openai.WithHTTPClient(httputil.DebugHTTPClient),
	)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "You are a golang expert"),
		llms.TextParts(llms.ChatMessageTypeHuman, prompt),
	}

	if _, err = llm.GenerateContent(ctx, content,
		llms.WithMaxTokens(4096),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		})); err != nil {
		log.Fatal(err)
	}
}

func UseGoogle(prompt string, chunkChan chan<- []byte) {
	ctx := context.Background()
	llm, err := googleai.New(ctx, googleai.WithAPIKey(os.Getenv("GOOGLE_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}

	// completion, err := llm.Call(ctx, prompt, llms.WithMetadata(map[string]interface{}{}), llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
	// 	fmt.Println(string(chunk))
	// 	return nil
	// }))

	// if err != nil {
	// 	log.Fatal(err)
	// }

	answer, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)

	if err != nil {
		log.Fatal(err)
	}

	chunkChan <- []byte(answer)

}
