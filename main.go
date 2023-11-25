package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type Input struct {
	Body      string `json:"body"`
	From      string `json:"from"`
	Timestamp int64  `json:"timestamp"`
}

type Action struct {
	Action string `json:"action"`
	Label  string `json:"label"`
	URL    string `json:"url"`
}

type Output struct {
	Topic    string   `json:"topic"`
	Message  string   `json:"message"`
	Title    string   `json:"title"`
	Actions  []Action `json:"actions"`
	Priority int      `json:"priority"`
}

func convertInputToOutput(input Input, topic string, priority int) Output {
	callAction := Action{
		Action: "view",
		Label:  "Call",
		URL:    "tel://" + input.From,
	}

	return Output{
		Topic:    topic,
		Message:  input.Body,
		Title:    "From: " + input.From,
		Actions:  []Action{callAction},
		Priority: priority,
	}
}

func sendOutputToWebhook(output Output, webhookURL string) error {
	outputJSON, err := json.Marshal(output)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(outputJSON))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Webhook returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	// Parse JSON from request body
	var input Input
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Failed to parse input JSON", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if input.Body == "" || input.From == "" || input.Timestamp == 0 {
		http.Error(w, "Missing required fields in input JSON", http.StatusBadRequest)
		return
	}

	// Get topic and priority from environment variables
	topic := os.Getenv("NTFY_TOPIC")
	if topic == "" {
		topic = "Messages" // Default value if not provided in environment variables
	}

	priorityStr := os.Getenv("NTFY_PRIORITY")
	priority, err := strconv.Atoi(priorityStr)
	if err != nil {
		priority = 1 // Default value if not provided or invalid
	}

	// Convert input to output
	output := convertInputToOutput(input, topic, priority)

	// Print output to terminal as formatted JSON
	outputJSON, err := json.MarshalIndent(output, "", "    ")
	if err != nil {
		http.Error(w, "Failed to marshal output to JSON", http.StatusInternalServerError)
		return
	}
	fmt.Printf("Received Input:\n%+v\n\nConverted Output:\n%s\n", input, outputJSON)

	// Send output to another webhook if the URL is provided
	webhookURL := os.Getenv("NTFY_URL")
	if webhookURL != "" {
		err := sendOutputToWebhook(output, webhookURL)
		if err != nil {
			fmt.Printf("Error sending output to webhook: %v\n", err)
		} else {
			fmt.Printf("Output sent to webhook successfully\n")
		}
	}

	// You can do further processing with the output here

	// Respond to the webhook request
	w.WriteHeader(http.StatusOK)
}

func main() {
	// Define the endpoint for the webhook
	http.HandleFunc("/webhook", webhookHandler)

	// Start the server
	port := 8080
	fmt.Printf("Server listening on :%d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
