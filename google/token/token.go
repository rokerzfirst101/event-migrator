package token

import (
	"encoding/json"
	"fmt"
	"net/http"

	gofrErr "developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	GoogleCalendar = "calendar"
)

func NewClientFromEnv(ctx *gofr.Context, credentialEnvKey, tokenEnvKey string, scopes []string) (*http.Client, error) {
	credentialData, err := readEnvValue(ctx, credentialEnvKey)
	if err != nil {
		return nil, err
	}

	tokenData, err := readEnvValue(ctx, tokenEnvKey)
	if err != nil {
		return nil, err
	}

	var token *oauth2.Token

	err = json.Unmarshal(tokenData, &token)
	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(credentialData, scopes...)
	if err != nil {
		return nil, err
	}

	return config.Client(ctx, token), nil
}

func readEnvValue(ctx *gofr.Context, envKey string) ([]byte, error) {
	value := ctx.Config.Get(envKey)
	if value == "" {
		return nil, &gofrErr.Response{
			StatusCode: http.StatusInternalServerError,
			Reason:     fmt.Sprintf("Unable to read %s from environment", envKey),
		}
	}

	return []byte(value), nil
}
