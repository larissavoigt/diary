package templates

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Templates struct {
	*template.Template
}

func New(path string) *Templates {
	t := template.New(path)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			template.Must(t.ParseFiles(path))
		}
		return err
	})
	if err != nil {
		log.Fatal(err)
	}
	return &Templates{t}
}

func (t *Templates) Render(w http.ResponseWriter, name string, data interface{}) {
	err := t.ExecuteTemplate(w, name+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t *Templates) NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	t.Render(w, "404", nil)
}

func (t *Templates) Error(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	t.Render(w, "500", err)
}
