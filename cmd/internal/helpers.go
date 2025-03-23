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

func CheckEndpointHealth(unhealthyEndpoints *[]string, endpoint string) {
	resp, err := http.Get(endpoint)
	if err != nil || resp.StatusCode >= 400 {
		resp, err = http.Get(endpoint) // retry HTTP request
		if err != nil || resp.StatusCode >= 400 {
			*unhealthyEndpoints = append(*unhealthyEndpoints, endpoint) // if unsuccessful, add endpoint to array
		}
	}
}

func SendNotification(allSuccessful bool, unhealthyEndpoints []string) error {
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
		NotificationStatus:     "Pending",
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

	notificationServerEndpoint := os.Getenv("NOTIFICATION_SERVER_ENDPOINT")
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
		req.Header.Set("X-API-Key", os.Getenv("NOTIFICATION_SERVER_API_KEY"))

		client := &http.Client{Timeout: 5 * time.Second}
		return client.Do(req)
	}

	resp, err := sendNotification()
	if err != nil || resp.StatusCode >= 400 {
		fmt.Println("First attempt to send notification failed:", err, "Status Code:", resp.StatusCode)
		notification.NotificationStatus = "Retrying"
		jsonData, _ = json.Marshal(notification)

		time.Sleep(2 * time.Second) // wait before retrying

		resp, err = sendNotification()
		if err != nil {
			fmt.Println("Second attempt to send notification failed with error:", err)
			return fmt.Errorf("failed to send notification after retries: %w", err)
		} else if resp.StatusCode >= 400 {
			fmt.Println("Second attempt to send notification failed with status code:", resp.StatusCode)
			return fmt.Errorf("failed to send notification after retries: status code %d", resp.StatusCode)
		}
	}

	defer resp.Body.Close()
	fmt.Println("Notification sent successfully")
	return nil
}
