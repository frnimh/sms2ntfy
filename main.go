package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func convertInputToOutput(input Input) Output {
	return Output{
		Topic:    "Messages",
		Message:  input.Body,
		Title:    "From: " + input.From,
		Call:     input.From,
		Priority: 1,
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

	// Convert input to output
	output := convertInputToOutput(input)

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
