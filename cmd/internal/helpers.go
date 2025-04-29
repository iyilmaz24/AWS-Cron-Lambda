package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type Notification struct {
	NotificationEmail      bool     `json:"sendEmail"`
	NotificationSMS        bool     `json:"sendSms"`
	NotificationUrgency    string   `json:"urgency"`
	NotificationRecipients []string `json:"recipient"`
	NotificationStatus     string   `json:"status"`
	NotificationID         string   `json:"id"`
	NotificationType       string   `json:"type"`
	NotificationSource     string   `json:"source"`
	NotificationTime       string   `json:"time"`
	NotificationDate       string   `json:"date"`
	NotificationTimezone   string   `json:"timezone"`
	NotificationSubject    string   `json:"subject"`
	NotificationMessage    string   `json:"message"`
	AccessSecret           string   `json:"password"`
}

func CheckEndpointHealth(client *http.Client, unhealthyEndpoints *[]string, endpoint string, apiKey string) {

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Printf("***ERROR: Failed to create request for endpoint '%s': %v\n", endpoint, err)
		*unhealthyEndpoints = append(*unhealthyEndpoints, endpoint)
		return
	}

	if apiKey != "" {
		req.Header.Set("X-API-Key", apiKey)
	}

	isUnhealthy := true

	for attempt := 1; attempt <= 2; attempt++ { // max of 2 attempts
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("***ERROR: Attempt %d failed: network error: %v\n", attempt, err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode < 400 {
			isUnhealthy = false
			break
		}
	}

	if isUnhealthy {
		*unhealthyEndpoints = append(*unhealthyEndpoints, endpoint)
	}
}

func GetNotificationObject(allSuccessful bool, unhealthyEndpoints []string) Notification {

	now := time.Now()
	notificationDate := now.Format("2006-01-02") // "2006-01-02 15:04:05.999999999 -0700 MST" - Golang's reference timestamp
	notificationTime := now.Format("15:04:05")
	notificationTimezone := now.Format("MST")

	notificationRecipients := strings.Split(os.Getenv("EMAIL_NOTIFICATION_RECIPIENTS"), ",")
	notificationId := fmt.Sprintf("health-%d", time.Now().Unix()) // health-1711147283

	urgency := "high"
	if allSuccessful {
		urgency = "low"
	}

	messageLines := []string{}
	if !allSuccessful {
		messageLines = append(messageLines, "<strong>Unhealthy Endpoints:</strong><br/>")
		for _, endpoint := range unhealthyEndpoints {
			messageLines = append(messageLines, "• "+endpoint+"<br/>")
		}
	} else {
		messageLines = append(messageLines, "<strong>All Endpoints Healthy.</strong><br/>")
	}

	notification := Notification{
		NotificationEmail:      true,
		NotificationSMS:        true,
		NotificationUrgency:    urgency,
		NotificationRecipients: notificationRecipients,
		NotificationStatus:     "",
		NotificationID:         notificationId,
		NotificationType:       "CloudWatch CRON",
		NotificationSource:     "AWS Lambda - Health Monitor",
		NotificationTime:       notificationTime,
		NotificationDate:       notificationDate,
		NotificationTimezone:   notificationTimezone,
		NotificationSubject:    "Health Monitor Notification (AWS Lambda)",
		NotificationMessage:    strings.Join(messageLines, ""),
		AccessSecret:           "",
	}

	return notification
}

func SendNotification(client *http.Client, allSuccessful bool, unhealthyEndpoints []string) error {

	notificationServerEndpoint := os.Getenv("NOTIFICATION_SERVER_ENDPOINT")
	backendApiKey := os.Getenv("BACKEND_API_KEY")
	if notificationServerEndpoint == "" || backendApiKey == "" {
		return fmt.Errorf("Notification server endpoint and/or API key missing")
	}

	notification := GetNotificationObject(allSuccessful, unhealthyEndpoints)
	notification.NotificationStatus = "1st Attempt"

	jsonData, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	sendNotification := func() (*http.Response, error) {
		req, err := http.NewRequest("POST", notificationServerEndpoint, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-Key", backendApiKey)
		return client.Do(req)
	}

	resp, err := sendNotification()
	if err != nil {
		fmt.Println("***ERROR: First attempt to send notification failed with error:", err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode < 400 {
			fmt.Println("***INFO: Notification sent successfully")
			return nil
		}
		fmt.Printf("***ERROR: First attempt failed: received status code %d\n", resp.StatusCode)
	}

	notification.NotificationStatus = "2nd Attempt"
	jsonData, err = json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}
	time.Sleep(2 * time.Second) // wait before retrying

	resp, err = sendNotification()
	if err != nil {
		fmt.Println("***ERROR: Second attempt to send notification failed with error:", err)
		return fmt.Errorf("failed to send notification after retries: %w", err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode < 400 {
			fmt.Println("***INFO: Notification sent successfully")
			return nil
		}
		fmt.Printf("***ERROR: Second attempt failed: received status code %d\n", resp.StatusCode)
		return fmt.Errorf("failed to send notification after retries: %w", err)
	}

}
