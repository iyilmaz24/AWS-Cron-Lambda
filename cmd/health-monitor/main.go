package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/iyilmaz24/AWS-Cron-Lambda/internal"
)

func handler(ctx context.Context) (string, error) {
	sites := os.Getenv("SITES_TO_MONITOR")
	servers := os.Getenv("SERVERS_TO_MONITOR")

	fmt.Println("Monitoring Sites:", sites)
	fmt.Println("Monitoring Servers:", servers)

	sitesList := strings.Split(sites, ",")
	serversList := strings.Split(servers, ",")
    unhealthyEndpoints := []string{}

	for _, site := range sitesList {
        internal.CheckEndpointHealth(&unhealthyEndpoints, site)
	}

	for _, server := range serversList {
		internal.CheckEndpointHealth(&unhealthyEndpoints, server)
	}

    allSuccessful := len(unhealthyEndpoints) == 0
    err := internal.SendNotification(allSuccessful, unhealthyEndpoints)
    if err != nil {
        fmt.Println("Error sending notification:", err)
        return "", err // this allows AWS to recognize the Lambda invocation as a failure
    }

    return "Health check completed", nil
}

func main() {
	lambda.Start(handler)
}
