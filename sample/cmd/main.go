package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("works")
		io.WriteString(w, "Hello, world!\n")
	}

	fmt.Println("started sample service")
	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":9005", nil))
}
