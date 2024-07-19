# firebase-spells-go

Go library for Firebase authentication

Because sometimes you don't want to drag an interpreter/runtime into your environment.

## Auth usage

### Prerequisites

You need to have a Firebase API key and a service account key file and it's corresponding gpg key.

Import the library. The default name is `fb` for `firebase`.

In your app, bring in the API key and service account key path through .env.

```go
// main.go

func main() {
    // ...
	serviceAccountKeyPath := os.Getenv("SERVICE_ACCOUNT_KEY_PATH")
	apiKey := os.Getenv("FIREBASE_WEB_API_KEY")
}
```

### Functions

**Initialize**: Initialize the library with the service account key file
```go
	fa := &fb.FirebaseAuth{}

	if err := fa.Initialize(serviceAccountKeyPath); err != nil {
		log.Fatalf("Failed to initialize Firebase Auth: %v", err)
	}
```

Set the API key
```go
fa.SetAPIKey(apiKey)
```

**AuthenticateUser(userEmail, userPassword string) (string, error)**: You can get an ID token with a user's email and password registered in the Firebase app

```go
	idToken, err := fa.AuthenticateUser(userEmail, userPassword)
	if err != nil {
		log.Fatalf("Failed to authenticate user: %v", err)
	}
```

**GetSessionCookie(expiresIn int64) (string, nil)**: The FirebaseAuth struct stores the ID token from the command above and uses it to get a session cookie. You just need to provide the expiry time as an argument.

```go
	sessionCookie, err := fa.GetSessionCookie(60 * 60 * 24 * 7) // 1 week
	if err != nil {
		log.Fatalf("Failed to get session cookie: %v", err)
	}
	fmt.Println("Session cookie:", sessionCookie)
```

**AuthUser() (*auth.Token, error)**: Get some information about the user from their decoded token such as their user_id.

```go
	decodedToken, err := fa.AuthUser()
	if err != nil {
		log.Fatalf("Failed to authenticate user with session cookie: %v", err)
	}
	fmt.Println("Decoded Token:", decodedToken)
```

**CreateCustomToken(uid string) (string, error)**: Create a custom token with the user uid.

```go
	customToken, err := fa.CreateCustomToken(user_id)
	if err != nil {
		log.Fatalf("Failed to create custom token: %v", err)
	}
```
