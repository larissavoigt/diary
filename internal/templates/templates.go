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

func (t *Templates) Render(res http.ResponseWriter, name string, data interface{}) {
	err := t.ExecuteTemplate(res, name+".html", data)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}
