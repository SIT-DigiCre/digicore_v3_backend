package google

import (
	"database/sql"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Context struct {
	Config *oauth2.Config
	DB     *sql.DB
}

func CreateContext(db *sql.DB) (Context, error) {
	context := Context{DB: db}
	clientSecretJson, err := os.ReadFile("client_secret.json")
	if err != nil {
		return context, fmt.Errorf("unknow error: %w", err)
	}
	context.Config, err = google.ConfigFromJSON(clientSecretJson, "https://www.googleapis.com/auth/userinfo.email")
	if err != nil {
		return context, fmt.Errorf("unable to parse client_secret.json: %w", err)
	}

	return context, nil
}
