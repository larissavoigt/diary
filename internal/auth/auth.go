package auth

import (
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

var config *oauth2.Config

// Config sets the Facebook configuration
func Config(id, secret, domain string) {
	config = &oauth2.Config{
		Endpoint:     facebook.Endpoint,
		ClientID:     id,
		ClientSecret: secret,
		RedirectURL:  fmt.Sprintf("%s/auth", domain),
	}
}

// RedirectURL returns an auth url for Facebook.
func RedirectURL() string {
	return config.AuthCodeURL("state")
}

// GetToken parses and return the Facebook token.
func GetToken(code string) (string, error) {
	t, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return "", err
	}
	return t.AccessToken, nil
}
