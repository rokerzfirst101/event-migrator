package token

import (
	"fmt"
	"net/http"

	gofrErr "developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func NewClientFromEnv(ctx *gofr.Context, credentialEnvKey, tokenEnvKey string, scopes []string) (*http.Client, error) {
	credentialData, err := readEnvValue(ctx, credentialEnvKey)
	if err != nil {
		return nil, err
	}

	credentialBytes := []byte(credentialData)

	refreshToken, err := readEnvValue(ctx, tokenEnvKey)
	if err != nil {
		return nil, err
	}

	token := &oauth2.Token{RefreshToken: refreshToken}

	config, err := google.ConfigFromJSON(credentialBytes, scopes...)
	if err != nil {
		return nil, err
	}

	return config.Client(ctx, token), nil
}

func readEnvValue(ctx *gofr.Context, envKey string) (string, error) {
	value := ctx.Config.Get(envKey)
	if value == "" {
		return "", &gofrErr.Response{
			StatusCode: http.StatusInternalServerError,
			Reason:     fmt.Sprintf("Unable to read %s from environment", envKey),
		}
	}

	return value, nil
}
