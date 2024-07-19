package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/replmade/firebase-spells-go/fb"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	serviceAccountKeyPath := os.Getenv("SERVICE_ACCOUNT_KEY_PATH")
	apiKey := os.Getenv("FIREBASE_WEB_API_KEY")
	userEmail := os.Getenv("USER_EMAIL")
	userPassword := os.Getenv("USER_PASSWORD")

	if serviceAccountKeyPath == "" || apiKey == "" || userEmail == "" || userPassword == "" {
		log.Fatalf("Missing required environment variables")
	}

	fa := &fb.FirebaseAuth{}

	if err := fa.Initialize(serviceAccountKeyPath); err != nil {
		log.Fatalf("Failed to initialize Firebase Auth: %v", err)
	}

	fa.SetAPIKey(apiKey)

	_, err = fa.AuthenticateUser(userEmail, userPassword)
	if err != nil {
		log.Fatalf("Failed to authenticate user: %v", err)
	}

	_, err = fa.GetSessionCookie(60 * 60 * 24 * 7) // 1 week
	if err != nil {
		log.Fatalf("Failed to get session cookie: %v", err)
	}

	_, err = fa.AuthUser()
	if err != nil {
		log.Fatalf("Failed to authenticate user with session cookie: %v", err)
	}

	_, err = fa.CreateCustomToken("user_uid")
	if err != nil {
		log.Fatalf("Failed to create custom token: %v", err)
	}
}
