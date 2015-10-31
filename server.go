package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"

	"github.com/larissavoigt/diary/internal/auth"
	"github.com/larissavoigt/diary/internal/db"
	"github.com/larissavoigt/diary/internal/templates"
	"github.com/larissavoigt/diary/internal/user"
	"github.com/rs/xhandler"
)

func main() {
	tpl := templates.New("templates")

	// chain authenticated middleware
	c := xhandler.Chain{}
	c.UseC(func(next xhandler.HandlerC) xhandler.HandlerC {
		return auth.NewMiddleware(next)
	})

	// server static assets files
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.Handle("/entries/", c.Handler(
		xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			u := ctx.Value("user").(*user.User)

			switch r.Method {
			case "GET":

				switch r.URL.Path[len("/entries/"):] {
				case "":
					entries, err := db.FindUserEntries(u.ID)
					if err != nil {
						tpl.Error(w, err)
					} else {
						tpl.Render(w, "entries", entries)
					}
				case "new":
					tpl.Render(w, "new_entry", u)
				default:
					fmt.Fprintf(w, "yay")
				}

			case "POST":
				rate := r.FormValue("rate")
				desc := r.FormValue("description")
				_, err := db.CreateEntry(u.ID, rate, desc)
				if err != nil {
					tpl.Error(w, err)
				} else {
					http.Redirect(w, r, "/entries", http.StatusFound)
				}
			default:
				http.Error(w, "", http.StatusMethodNotAllowed)
			}
		})))

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		token, err := auth.GetToken(code)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		} else {
			id, err := db.CreateUser(token)
			if err != nil {
				tpl.Error(w, err)
			} else {
				auth.SaveSession(w, id)
				http.Redirect(w, r, "/entries/new", http.StatusFound)
			}
		}
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		auth.DestroySession(w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			if r.URL.Path != "/" {
				tpl.NotFound(w)
				return
			}
			_, err := auth.CurrenUser(r)
			if err == nil {
				http.Redirect(w, r, "/entries/new", 302)
			} else {
				p := struct {
					FacebookURL string
				}{
					auth.RedirectURL(),
				}
				tpl.Render(w, "index", p)
			}
		} else {
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":3000", nil)
}
