package main

import (
	"net/http"

	"github.com/larissavoigt/diary/internal/auth"
	"github.com/larissavoigt/diary/internal/templates"
)

func main() {
	var tpl = templates.New("templates")
	auth.Config("1629858967301577", "36b8b62d4a6d62f3e845a2682698749d", "http://localhost:3000")

	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/auth", func(res http.ResponseWriter, req *http.Request) {
		//code := req.URL.Query().Get("code")
		//_, err := auth.GetToken(code)
		tpl.Render(res, "entry", nil)
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
