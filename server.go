package main

import (
	"html/template"
	"net/http"
)

type Page struct {
	Title string
}

func main() {
	var tpl = template.Must(template.ParseFiles("templates/index.html"))

	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/about", func(res http.ResponseWriter, req *http.Request) {
		var page = &Page{"About"}
		tpl.Execute(res, page)
	})

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		var page = &Page{"MyDiary"}
		tpl.Execute(res, page)
	})

	http.ListenAndServe(":3000", nil)
}
