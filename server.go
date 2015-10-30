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

	http.HandleFunc("/entry", func(res http.ResponseWriter, req *http.Request) {
		id, err := auth.CurrentUser(req)
		if err != nil {
			log.Println(err)
			http.Redirect(res, req, "/", 302)
			return
		}
		user, err := db.FindUser(id)
		if err != nil {
			log.Println(err)
			http.Redirect(res, req, "/", 302)
			return
		}
		switch req.Method {
		case "GET":
			tpl.Render(res, "entry", user)
		case "POST":
			rate := req.FormValue("rate")
			desc := req.FormValue("description")
			e, err := db.CreateEntry(id, rate, desc)
			if err != nil {
				log.Println(err)
				http.Redirect(res, req, "/entry", 302)
			} else {
				http.Redirect(res, req, "/entry/"+e, 302)
			}
		default:
			http.Error(res, "", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/auth", func(res http.ResponseWriter, req *http.Request) {
		code := req.URL.Query().Get("code")
		token, err := auth.GetToken(code)
		if err != nil {
			log.Println(err)
			http.Redirect(res, req, "/", 302)
			return
		}
		id, err := db.CreateUser(token)
		if err != nil {
			log.Println(err)
			http.Redirect(res, req, "/", 302)
		} else {
			auth.SaveSession(res, id)
			http.Redirect(res, req, "/entry", 302)
		}
	})

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		p := struct {
			FacebookURL string
		}{
			auth.RedirectURL(),
		}
		tpl.Render(res, "index", p)
	})
	http.ListenAndServe(":3000", nil)
}
