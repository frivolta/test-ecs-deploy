package util

import (
	"bytes"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"path/filepath"
)

func SetupFirebase(serviceKeyPath string) *auth.Client {
	serviceAccountKeyFilePath, err :=
		filepath.Abs(serviceKeyPath)
	if err != nil {
		panic("Unable to load serviceAccountKeys.json file")
	}
	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)
	//Firebase admin SDK initialization
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatal(fmt.Sprintf("cannot init firebase %s", err))
	}
	//Firebase Auth
	auth, err := app.Auth(context.Background())
	if err != nil {
		log.Fatal(fmt.Sprintf("app auth error: %s", err))
	}
	return auth
}

func GetIdTokenFromUUID(uuid string, config *Config) (string, error) {
	token, err := config.AuthClient.CustomToken(context.Background(), uuid)
	if err != nil {
		return "", err
	}
	args := map[string]interface{}{
		"token":             token,
		"returnSecureToken": true,
	}
	jd, err := json.Marshal(args)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken?key=%s", config.FirebaseWeb)
	res, err := http.Post(url, "application/json", bytes.NewBuffer(jd))
	if err != nil {
		return "", err
	}
	var data map[string]interface{}
	json.NewDecoder(res.Body).Decode(&data)
	return data["idToken"].(string), nil
}

func GenerateBearerString(token string) string {
	return fmt.Sprintf("Bearer %s", token)
}
