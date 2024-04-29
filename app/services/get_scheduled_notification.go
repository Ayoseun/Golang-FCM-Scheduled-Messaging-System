package app

import (
	"context"
	"fmt"
	"log"
	"time"
	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	
)

func FetchScheduledMessages() {
	loc, err := time.LoadLocation("Africa/Lagos")
	if err != nil {
		panic(err)
	}

	now := time.Now().In(loc)
	app, client, err := initializeFirebaseApp()
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v", err)
	}
	defer client.Close()

	collectionPath := "scheduledMessaging"
	docs, err := client.Collection(collectionPath).Documents(context.Background()).GetAll()
	if err != nil {
		log.Fatalf("Error getting documents: %v", err)
	}

	processDocuments(app, docs, now, loc)
}



func processDocuments(app *firebase.App, docs []*firestore.DocumentSnapshot, now time.Time, loc *time.Location) {
	for _, doc := range docs {
		docData := doc.Data()
		isSent, _ := docData["isSent"].(bool)

		if !isSent {
			period, ok := docData["date"].(string)
			if !ok {
				log.Printf("Error getting date field for document %s: field is not a string", doc.Ref.ID)
				continue
			}
			parsedTime, err := time.Parse(time.RFC3339, period)
			if err != nil {
				log.Printf("Error parsing date field for document %s: %v", doc.Ref.ID, err)
				continue
			}

			if parsedTime.In(loc).Before(now) || parsedTime.In(loc).Equal(now) {
				usersTokenData, ok := docData["users"].([]interface{})
				if !ok {
					log.Printf("Error getting users field for document %s: field is not a slice", doc.Ref.ID)
					continue
				}

				usersTokens := make([]string, len(usersTokenData))
				for i, v := range usersTokenData {
					token, ok := v.(string)
					if !ok {
						log.Printf("Error converting user token for document %s: value is not a string", doc.Ref.ID)
						continue
					}
					usersTokens[i] = token
				}

				header, _ := docData["header"].(string)
				body, _ := docData["body"].(string)
				image, _ := docData["image"].(string)

				_, err = SendNotification(app, header, body, image, usersTokens)
				if err != nil {
					log.Printf("Error sending notification for document %s: %v", doc.Ref.ID, err)
				} else {
					err = updateIsSentField(doc, loc)
					if err != nil {
						log.Printf("Error updating IsSent field for document %s: %v", doc.Ref.ID, err)
					}
				}
			}
		}
	}
}

func updateIsSentField(doc *firestore.DocumentSnapshot, loc *time.Location) error {
	_, err := doc.Ref.Update(context.Background(), []firestore.Update{
		{Path: "isSent", Value: true},
	})
	if err != nil {
		return fmt.Errorf("error updating IsSent field: %v", err)
	}
	return nil
}