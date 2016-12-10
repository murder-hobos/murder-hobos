package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/murder-hobos/murder-hobos/routes"
)

func main() {
	// setup DB
	dbconfig := mysql.Config{
		User:            os.Getenv("MYSQL_USER"),
		Passwd:          os.Getenv("MYSQL_PASS"),
		DBName:          os.Getenv("MYSQL_DB_NAME"),
		Net:             "tcp",
		Addr:            os.Getenv("MYSQL_ADDR"),
		MultiStatements: false,
	}

	log.Println(dbconfig.FormatDSN())
	r := routes.New(dbconfig.FormatDSN())
	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(":"+port, r))
}
