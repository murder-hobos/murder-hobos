package main

import (
	"fmt"
	"net/http"
	"os"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello murder-hobos")
}

func main() {
	http.HandleFunc("/", hello)
	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, nil)
}
