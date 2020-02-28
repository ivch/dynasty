package ui

import (
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

func NewHTTPHandler() http.Handler {
	r := chi.NewRouter()

	r.Get("/guard", func(w http.ResponseWriter, r *http.Request) {
		fp, err := os.Open("../_ui/guard/index.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		io.Copy(w, fp)
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
