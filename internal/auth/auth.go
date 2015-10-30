package auth

import (
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

var config *oauth2.Config

var id = 0
var users = map[string]string{
	"1": "Larissa",
	"2": "Luiz",
	"3": "Cafe",
}

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

func User(res http.ResponseWriter, req *http.Request) (string, error) {
	cookie, err := req.Cookie("id")
	if err == nil {
		user, _ := users[cookie.Value]
		return user, nil
	} else {
		id++
		cookie := &http.Cookie{
			Name:     "id",
			Value:    strconv.Itoa(id),
			HttpOnly: true,
		}
		http.SetCookie(res, cookie)
		return "Visitor", nil
	}
}
