package main

import (
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

type Output struct {
	Topic    string `json:"topic"`
	Message  string `json:"message"`
	Title    string `json:"title"`
	Call     string `json:"call"`
	Priority int    `json:"priority"`
}

func convertInputToOutput(input Input, topic string, priority int) Output {
	return Output{
		Topic:    topic,
		Message:  input.Body,
		Title:    "From: " + input.From,
		Call:     input.From,
		Priority: priority,
	}
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
