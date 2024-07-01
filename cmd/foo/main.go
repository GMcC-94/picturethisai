package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/replicate/replicate-go"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	ctx := context.TODO()

	// You can also provide a token directly with
	// `replicate.NewClient(replicate.WithToken("r8_..."))`
	r8, err := replicate.NewClient(replicate.WithTokenFromEnv())
	if err != nil {
		log.Fatal(err)
	}

	model := "stability-ai/sdxl"
	version := "7762fd07cf82c948538e41f63f77d685e02b063e37e496e96eefd46c929f9bdc"

	input := replicate.PredictionInput{
		"prompt": "An astronaut riding a rainbow unicorn",
	}

	webhook := replicate.Webhook{
		URL:    "https://webhook.site/a6d17d56-f683-4f59-8ee9-7ab6bb75ec14",
		Events: []replicate.WebhookEventType{"start", "completed"},
	}

	// Run a model and wait for its output
	output, _ := r8.Run(ctx, fmt.Sprintf("%s:%s", model, version), input, &webhook)

	prediction, err := r8.CreatePrediction(ctx, version, input, &webhook, false)
	if err != nil {
		log.Fatal(err)
	}
	_ = r8.Wait(ctx, prediction) // Wait for the prediction to finish

	fmt.Println("output: ", output)
}
