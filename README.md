# AWS Cron Lambda Health Monitor
This project implements an AWS Lambda function that monitors the health of specified websites and servers. It sends email and SMS notifications if any endpoints are unhealthy. The Lambda function is configured to run on a schedule using CloudWatch Events (Cron).  The application integrates with AWS Systems Manager Parameter Store for environment variable management.

## Features
* Monitors the health of multiple websites and servers.
* Sends email and SMS notifications upon detecting unhealthy endpoints.
* Uses AWS Lambda for serverless execution.
* Leverages CloudWatch Events for scheduled execution.
* Implements retry mechanisms for notification sending.
* Utilizes AWS Systems Manager Parameter Store for secure configuration management.

## Usage
1. **Configure Environment Variables:**  Set the following environment variables in AWS Systems Manager Parameter Store:
    * `SITES_TO_MONITOR`: A comma-separated list of website URLs to monitor (e.g., `https://www.example.com,https://www.google.com`).
    * `SERVERS_TO_MONITOR`: A comma-separated list of server URLs to monitor (e.g., `http://192.168.1.100:8080,http://localhost:3000`).
    * `EMAIL_NOTIFICATION_RECIPIENTS`: A comma-separated list of email addresses for notifications (e.g., `user1@example.com,user2@example.com`).
    * `NOTIFICATION_SERVER_ENDPOINT`: The endpoint URL of your notification server.
    * `NOTIFICATION_SERVER_API_KEY`: The API key for your notification server.

2. **Deploy the Lambda Function:** Use the provided `Makefile` to deploy the Lambda function to AWS:

   ```bash
   make deploy
   ```

3. **Configure CloudWatch Events:** Create a CloudWatch Events rule to trigger the Lambda function on your desired schedule (e.g., every 5 minutes).

## Installation
1. **Clone the repository:**
   ```bash
   git clone https://github.com/iyilmaz24/AWS-Cron-Lambda.git
   ```

2. **Navigate to the project directory:**
   ```bash
   cd AWS-Cron-Lambda/cmd
   ```

3. **Install dependencies:**
   ```bash
   go mod download
   ```

## Technologies Used
* **Go:** The programming language used for developing the Lambda function.
* **AWS Lambda:** A serverless compute service that runs the health check function.
* **AWS CloudWatch Events:** A service used to schedule the Lambda function's execution.
* **AWS Systems Manager Parameter Store:** Securely stores and retrieves environment variables.
* **`aws-lambda-go`:**  The AWS SDK for Go, used for interacting with AWS Lambda.
* **`net/http`:** The Go standard library package for making HTTP requests.
* **`encoding/json`:** The Go standard library package for handling JSON data.

## Configuration
The application is configured primarily through environment variables stored in AWS Systems Manager Parameter Store.  Refer to the "Usage" section for details on the required environment variables.

## Dependencies
The project integrates with the notification server at [https://github.com/iyilmaz24/Golang-Notification-Server](https://github.com/iyilmaz24/Golang-Notification-Server) for sending email and SMS notifications.  The `NOTIFICATION_SERVER_ENDPOINT` and `NOTIFICATION_SERVER_API_KEY` environment variables must be configured to point to your instance of this server.  The application uses the server's API to send notifications regarding the health of monitored websites and servers.

*README.md was made with [Etchr](https://etchr.dev)*