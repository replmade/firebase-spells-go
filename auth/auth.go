package fb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

type FirebaseAuth struct {
	baseURL       string
	apiKey        string
	authURL       string
	idToken       string
	sessionCookie string
	decodedToken  *auth.Token
	customToken   string
	authClient    *auth.Client
}

func (fa *FirebaseAuth) Initialize(serviceAccountKeyPath string) error {
	ctx := context.Background()
	opt := option.WithCredentialsFile(serviceAccountKeyPath)
	app, err := firebase.NewApp(ctx, nil, opt)
	// app, err := auth.Client(ctx, nil, opt)
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}
	fa.authClient, err = app.Auth(ctx)
	if err != nil {
		return fmt.Errorf("error getting Auth client: %v", err)
	}
	fa.baseURL = "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword"
	return nil
}

func (fa *FirebaseAuth) SetAPIKey(apiKey string) {
	fa.apiKey = apiKey
	fa.authURL = fmt.Sprintf("%s?key=%s", fa.baseURL, fa.apiKey)
}

func (fa *FirebaseAuth) SetIdToken(idToken string) {
	fa.idToken = idToken
}

func (fa *FirebaseAuth) SetSessionCookie(sessionCookie string) {
	fa.sessionCookie = sessionCookie
}

func (fa *FirebaseAuth) SetCustomToken(customToken string) {
	fa.customToken = customToken
}

func (fa *FirebaseAuth) AuthenticateUser(email, password string) (string, error) {
	if fa.apiKey == "" {
		return "", fmt.Errorf("API Key is not set")
	}

	payload := map[string]string{
		"email":             email,
		"password":          password,
		"returnSecureToken": "true",
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshalling payload: %v", err)
	}

	resp, err := http.Post(fa.authURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	var responseData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	if resp.StatusCode == 200 {
		idToken, ok := responseData["idToken"].(string)
		if !ok {
			return "", fmt.Errorf("ID token not found in response")
		}
		fa.idToken = idToken
		return idToken, nil
	}

	errorMessage, _ := responseData["error"].(map[string]interface{})["message"].(string)
	if errorMessage == "" {
		errorMessage = "Unknown error occurred"
	}
	fmt.Printf("Failed to authenticate. Error: %s\n", errorMessage)
	return "", fmt.Errorf("authentication failed: %s", errorMessage)
}

func (fa *FirebaseAuth) GetSessionCookie(expiresIn int64) (string, error) {
	if fa.idToken == "" {
		return "", fmt.Errorf("ID Token is not set")
	}

	ctx := context.Background()
	expiresInDuration := time.Duration(expiresIn) * time.Second
	sessionCookie, err := fa.authClient.SessionCookie(ctx, fa.idToken, expiresInDuration)
	if err != nil {
		return "", fmt.Errorf("failed to create session cookie: %v", err)
	}

	fa.sessionCookie = sessionCookie
	return sessionCookie, nil
}

func (fa *FirebaseAuth) AuthUser() (*auth.Token, error) {
	if fa.sessionCookie == "" {
		return nil, fmt.Errorf("session Cookie is not set")
	}

	ctx := context.Background()
	decodedToken, err := fa.authClient.VerifySessionCookieAndCheckRevoked(ctx, fa.sessionCookie)
	if err != nil {
		return nil, fmt.Errorf("failed to verify session cookie: %v", err)
	}

	fa.decodedToken = decodedToken
	return decodedToken, nil
}

func (fa *FirebaseAuth) GetUserByUID(uid string) (*auth.UserRecord, error) {
	ctx := context.Background()
	userRecord, err := fa.authClient.GetUser(ctx, uid)
	if err != nil {
		if auth.IsUserNotFound(err) {
			fmt.Println("User not found.")
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return userRecord, nil
}

func (fa *FirebaseAuth) CreateCustomToken(uid string) (string, error) {
	ctx := context.Background()
	customToken, err := fa.authClient.CustomToken(ctx, uid)
	if err != nil {
		return "", fmt.Errorf("failed to create custom token: %v", err)
	}

	fa.customToken = customToken
	return customToken, nil
}
