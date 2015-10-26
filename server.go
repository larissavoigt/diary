package main

import (
	"net/http"
	"text/template"
)

type Page struct {
	Title string
}

func main() {
	var tpl = template.Must(template.ParseFiles("templates/index.html"))
	http.HandleFunc("/about", func(res http.ResponseWriter, req *http.Request) {
		var page = &Page{"About"}
		tpl.Execute(res, page)
	})
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		var page = &Page{"My Diary"}
		tpl.Execute(res, page)
	})
	http.ListenAndServe(":3000", nil)
}
