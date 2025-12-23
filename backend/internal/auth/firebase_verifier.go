package auth

import (
	"context"
	"os"

	firebase "firebase.google.com/go/v4"
	firebaseauth "firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type Verifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (string, error)
}

type FirebaseVerifier struct {
	client *firebaseauth.Client
}

func NewFirebaseVerifier(ctx context.Context) (*FirebaseVerifier, error) {
	var opts []option.ClientOption
	credPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credPath != "" {
		opts = append(opts, option.WithCredentialsFile(credPath))
	}

	app, err := firebase.NewApp(ctx, nil, opts...)
	if err != nil {
		return nil, err
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return &FirebaseVerifier{client: client}, nil
}

func (v *FirebaseVerifier) VerifyIDToken(ctx context.Context, idToken string) (string, error) {
	tok, err := v.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return "", err
	}
	return tok.UID, nil
}
