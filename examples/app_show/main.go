package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Scalingo/go-scalingo/v9"
)

func main() {
	ctx := context.Background()

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: ./app_show [REGION] [APP_NAME]\n")
		os.Exit(1)
	}
	region := os.Args[1]
	appName := os.Args[2]

	token := os.Getenv("SCALINGO_TOKEN")
	if token == "" {
		fmt.Fprintf(os.Stderr, "SCALINGO_TOKEN is not set, please set it before using this example\n")
		os.Exit(1)
	}

	client, err := scalingo.New(ctx, scalingo.ClientConfig{
		Region:   region,
		APIToken: token,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "fail to create scalingo client: %s\n", err.Error())
		os.Exit(1)
	}

	app, err := client.AppsShow(ctx, appName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fail to show app: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("ID: %s\n", app.ID)
	fmt.Printf("Owner: %s\n", app.Owner.Username)
	fmt.Printf("Status: %s\n", app.Status)
}
