package main

import (
	"log"
	"net/http"

	"github.com/larissavoigt/diary/internal/auth"
	"github.com/larissavoigt/diary/internal/db"
	"github.com/larissavoigt/diary/internal/templates"
)

func main() {
	var tpl = templates.New("templates")
	auth.Config("1629858967301577", "36b8b62d4a6d62f3e845a2682698749d", "http://localhost:3000")

	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/auth", func(res http.ResponseWriter, req *http.Request) {
		code := req.URL.Query().Get("code")
		token, err := auth.GetToken(code)
		if err != nil {
			log.Println(err)
			http.Redirect(res, req, "/", 302)
			return
		}
		_, err = db.CreateUser(token)
		if err != nil {
			log.Println(err)
			http.Redirect(res, req, "/", 302)
		} else {
			tpl.Render(res, "entry", nil)
		}
	})

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		p := struct {
			FacebookURL string
			User        string
		}{
			auth.RedirectURL(),
			"Visitor",
		}
		tpl.Render(res, "index", p)
	})
	http.ListenAndServe(":3000", nil)
}
