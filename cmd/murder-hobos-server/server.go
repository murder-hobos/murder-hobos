package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jaden-young/murder-hobos/routes"
)

func main() {
	r := routes.New()
	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(":"+port, r))
}
