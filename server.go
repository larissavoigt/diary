package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/context"

	"github.com/larissavoigt/diary/internal/auth"
	"github.com/larissavoigt/diary/internal/db"
	"github.com/larissavoigt/diary/internal/templates"
	"github.com/rs/xhandler"
)

func main() {
	tpl := templates.New("templates")

	c := xhandler.Chain{}
	c.UseC(func(next xhandler.HandlerC) xhandler.HandlerC {
		return auth.NewMiddleware(next)
	})

	// server static assets files
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	entry := xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		u := ctx.Value("user").(*db.User)

		switch r.Method {
		case "GET":
			p := r.URL.Path[len("/entry/"):]
			if p == "" {
				tpl.Render(w, "entry", u)
			} else {
				fmt.Fprintf(w, "yay")
			}
		case "POST":
			rate := r.FormValue("rate")
			desc := r.FormValue("description")
			e, err := db.CreateEntry(u.ID, rate, desc)
			if err != nil {
				log.Println(err)
				http.Redirect(w, r, "/entry", 302)
			} else {
				http.Redirect(w, r, "/entry/"+e, 302)
			}
		default:
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	})

	http.Handle("/entry/", c.Handler(entry))

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
				http.Redirect(w, r, "/entry", http.StatusFound)
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			if r.URL.Path != "/" {
				tpl.NotFound(w)
				return
			}
			_, err := auth.CurrenUser(r)
			if err == nil {
				http.Redirect(w, r, "/entry", 302)
			} else {
				p := struct {
					FacebookURL string
				}{
					auth.RedirectURL(),
				}
				tpl.Render(w, "index", p)
			}
		default:
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":3000", nil)
}
