package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/iyilmaz24/AWS-Cron-Lambda/internal"
)

func handler(ctx context.Context) (string, error) {
	sites := os.Getenv("SITES_TO_MONITOR")
	servers := os.Getenv("SERVERS_TO_MONITOR")

	fmt.Println("***INFO: Monitoring Sites:", sites)
	fmt.Println("***INFO: Monitoring Servers:", servers)

	sitesList := strings.Split(sites, ",")
	serversList := strings.Split(servers, ",")
	unhealthyEndpoints := []string{}
	client := &http.Client{Timeout: 7 * time.Second} // 7 second max timeout for requests

	for _, site := range sitesList {
		internal.CheckEndpointHealth(client, &unhealthyEndpoints, site, "")
	}

	for _, server := range serversList {
		internal.CheckEndpointHealth(client, &unhealthyEndpoints, server, os.Getenv("BACKEND_API_KEY")) // backend servers require X-API-KEY header
	}

	allSuccessful := len(unhealthyEndpoints) == 0

	if !allSuccessful { // optional control structure - only send notification if there is an issue
		err := internal.SendNotification(client, allSuccessful, unhealthyEndpoints)
		if err != nil {
			fmt.Println("***ERROR: Error sending notification:", err)
			return "", err // this allows AWS to recognize the Lambda invocation as a failure
		}
	} else {
		fmt.Println("***INFO: All endpoints healthy - notification not sent")
	}

	return "Health check completed", nil
}

func main() {
	lambda.Start(handler)
}
