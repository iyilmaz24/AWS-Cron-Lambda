package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context) (string, error) {
    sites := os.Getenv("SITES_TO_MONITOR")
    servers := os.Getenv("SERVERS_TO_MONITOR")

    fmt.Println("Monitoring Sites:", sites)
    fmt.Println("Monitoring Servers:", servers)

    return "Health check completed", nil
}

func main() {
    lambda.Start(handler)
}
