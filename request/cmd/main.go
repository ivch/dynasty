package main

import "net/http"

func main() {
	http.HandleFunc("/test", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("works fine"))
	})

	http.HandleFunc("/foo", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("no auth i'm fine"))
	})

	http.ListenAndServe(":9999", nil)
}
