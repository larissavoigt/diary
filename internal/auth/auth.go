package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/larissavoigt/diary/internal/db"
	"github.com/larissavoigt/diary/internal/user"
	"github.com/rs/xhandler"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

var config *oauth2.Config

func Config(domain, port, client, secret string) {
	config = &oauth2.Config{
		Endpoint:     facebook.Endpoint,
		ClientID:     client,
		ClientSecret: secret,
		RedirectURL:  fmt.Sprintf("%s:%s/auth", domain, port),
	}
}

type Middleware struct {
	next xhandler.HandlerC
}

func (h Middleware) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user, err := CurrenUser(r)
	if err == nil {
		ctx = context.WithValue(ctx, "user", user)
		h.next.ServeHTTPC(ctx, w, r)
	} else {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

func NewMiddleware(next xhandler.HandlerC) *Middleware {
	return &Middleware{next: next}
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

func SaveSession(w http.ResponseWriter, id string) {
	cookie := &http.Cookie{
		Name:     "id",
		Value:    id,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func DestroySession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "id",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func CurrenUser(r *http.Request) (*user.User, error) {
	cookie, err := r.Cookie("id")
	if err == nil {
		id := cookie.Value
		if id != "" {
			user, err := db.FindUser(id)
			if err == nil {
				return user, nil
			}
		}
	}
	return nil, errors.New("User not found")
}
