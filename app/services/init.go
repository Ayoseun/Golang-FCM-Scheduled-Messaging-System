package app

import (
	"context"
	"encoding/json"
	"firebase-fcm-cron-job/app/config"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func createFirebaseCredentials(cfg *config.Config) ([]byte, error) {
	credentials := map[string]interface{}{
		"type":                        "service_account",
		"project_id":                  cfg.PROJECT_ID,
		"private_key_id":              cfg.PRIVATE_KEY_ID,
		"private_key":                 cfg.PRIVATE_KEY,
		"client_email":                cfg.CLIENT_EMAIL,
		"client_id":                   cfg.CLIENT_ID,
		"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
		"token_uri":                   "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url":        cfg.CLIENT_URL,
		"universe_domain":             "googleapis.com",
		"messagingSenderId":           cfg.MESSAGING_SENDER_ID,
		"appId":                       cfg.APP_ID,
	}

	return json.Marshal(credentials)
}

func initializeFirebaseApp() (*firebase.App, *firestore.Client, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load configuration: %v", err)
	}

	credentials, err := createFirebaseCredentials(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("error marshaling credentials: %v", err)
	}

	opt := option.WithCredentialsJSON(credentials)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing Firebase app: %v", err)
	}

	client, err := app.Firestore(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("error getting Firestore client: %v", err)
	}

	return app, client, nil
}