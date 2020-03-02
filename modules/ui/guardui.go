package ui

import (
	"html/template"
	"io"
	"net/http"
	"os"

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

	r.Get("/assets/img/logo.png", func(w http.ResponseWriter, r *http.Request) {
		// fp, _ := os.Open("../ui/guard/img/" + chi.URLParam(r, "filename"))
		fp, err := os.Open("../_ui/guard/img/logo.png")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		io.Copy(w, fp)
	})

	return r
}
