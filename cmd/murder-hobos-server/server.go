package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
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
		Addr:            os.Getenv("MYSQL_ADDR") + ":" + os.Getenv("MYSQL_PORT"),
		MultiStatements: false,
	}

	if os.Getenv("ENVIRONMENT") == "production" {
		rootCertPool := x509.NewCertPool()
		pem, err := ioutil.ReadFile("db/amazon-rds-ca-cert.pem")
		if err != nil {
			log.Fatal(err.Error())
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			log.Fatal("Failed to append PEM")
		}
		err = mysql.RegisterTLSConfig("amazon", &tls.Config{
			RootCAs:    rootCertPool,
			ServerName: os.Getenv("MYSQL_ADDR"),
		})
		if err != nil {
			log.Fatal(err.Error())
		}
		dbconfig.TLSConfig = "amazon"
	}

	log.Println(dbconfig.FormatDSN())
	r := routes.New(dbconfig.FormatDSN())
	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(":"+port, r))
}
