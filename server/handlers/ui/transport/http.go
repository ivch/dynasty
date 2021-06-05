package ui

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
)

type pageConfig struct {
	APIHost    string
	PageURI    string
	PagerLimit int
}

func NewHTTPHandler(apiHost, pageURI string, pagerLimit int) http.Handler {
	r := chi.NewRouter()

	tmpl := template.Must(template.ParseFiles("../_ui/guard/index.html"))
	cfg := pageConfig{
		APIHost:    apiHost,
		PageURI:    pageURI,
		PagerLimit: pagerLimit,
	}

	r.Get("/guard", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, cfg); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	r.Get("/assets/{folder:(img|css|js)}/{filename}", func(w http.ResponseWriter, r *http.Request) {
		folder := chi.URLParam(r, "folder")
		filename := chi.URLParam(r, "filename")
		ext := filepath.Ext(filename)

		fp, err := os.Open(fmt.Sprintf("../_ui/guard/%s/%s", folder, filename))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		switch ext {
		case ".js":
			w.Header().Add("Content-Type", "text/javascript")
		case ".css":
			w.Header().Add("Content-Type", "text/css")
		}
		io.Copy(w, fp)
	})

	r.Get("/docs/oferta", func(w http.ResponseWriter, r *http.Request) {
		fp, err := os.Open("../_ui/oferta.html")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Add("Content-Type", "text/html")
		io.Copy(w, fp)
	})

	return r
}
