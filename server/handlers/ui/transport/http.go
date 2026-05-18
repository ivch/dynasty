package ui

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	chi "github.com/go-chi/chi/v5"
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

		// Clean the filename to prevent path traversal
		cleanFilename := filepath.Clean(filename)
		// Ensure the cleaned filename doesn't contain path traversal or absolute paths
		if filepath.IsAbs(cleanFilename) || strings.Contains(cleanFilename, "..") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Build the safe path using filepath.Join
		basePath := filepath.Join("..", "_ui", "guard")
		fullPath := filepath.Join(basePath, folder, cleanFilename)

		// Double-check the resolved path is within the base directory
		absBase, err := filepath.Abs(basePath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		absFullPath, err := filepath.Abs(fullPath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !strings.HasPrefix(absFullPath, absBase) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ext := filepath.Ext(cleanFilename)
		// #nosec G703 -- Path traversal is prevented by multiple layers of validation above:
		// 1. filename is cleaned with filepath.Clean()
		// 2. absolute paths are rejected
		// 3. paths containing ".." are rejected
		// 4. resolved absolute path is verified to be within basePath
		fp, err := os.Open(fullPath)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		defer fp.Close()

		switch ext {
		case ".js":
			w.Header().Add("Content-Type", "text/javascript")
		case ".css":
			w.Header().Add("Content-Type", "text/css")
		}
		if _, err := io.Copy(w, fp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	r.Get("/docs/oferta", func(w http.ResponseWriter, r *http.Request) {
		fp, err := os.Open("../_ui/oferta.html")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		defer fp.Close()

		w.Header().Add("Content-Type", "text/html")
		if _, err := io.Copy(w, fp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	return r
}
