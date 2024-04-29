package app

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

func SendNotification(app *firebase.App, header string, body string,image string, tokens []string) (int, error) {

	fcmClient, err := app.Messaging(context.Background())
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}
	// Define the notification payload
	payload := &messaging.MulticastMessage{
		Notification: &messaging.Notification{ // Creating an instance of Notification struct
			Title:    header,
			Body:     body,
			ImageURL: image,
		},
		Tokens: nil,
	}

	// Define a map to store unique fcmTokens
	uniqueTokens := make(map[string]bool)

	// Iterate over the documents
	for _, fcmToken := range tokens {
		// Check if fcmToken is not nil
		if fcmToken != "" {
			// Convert fcmToken to a string
			token := fcmToken

			// Check if the token is already in the uniqueTokens map
			if _, exists := uniqueTokens[token]; !exists {
				// Token is not in the map, add it to the payload.Tokens slice
				payload.Tokens = append(payload.Tokens, token)

				// Mark the token as seen in the uniqueTokens map
				uniqueTokens[token] = true
			}

		} else {
			log.Printf("User %s has nil fcmToken", fcmToken)
		}
	}

	// Remove any empty slots from the payload.Tokens slice
	filteredTokens := make([]string, 0, len(payload.Tokens))
	for _, token := range payload.Tokens {
		if token != "" {
			filteredTokens = append(filteredTokens, token)
		}
	}
	payload.Tokens = filteredTokens

	// Send the notification
	response, err := fcmClient.SendMulticast(context.Background(), payload)
	if err != nil {
		log.Fatalf("Error sending notification: %v", err)
	}

	log.Printf("Notification sent to %d devices with header  %s: and message %s:", response.SuccessCount, header, body)

	return response.SuccessCount, err
}
