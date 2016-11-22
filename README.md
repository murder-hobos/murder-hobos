# murder-hobos

Murder Hobos is a DnD 5e Spellbook reference application... which may be extended as far as we have the will to.

## Getting Started

1. If you're keen on using windows to develop, follow the instructions on the golang 
   [webpage](https://golang.org/doc/install), using the [windows installer](https://storage.googleapis.com/golang/go1.7.3.windows-amd64.msi)
   If you're using a linux distro, go will be in the standard repos somewhere. Make sure to setup your path.

2. Follow the [heroku](https://devcenter.heroku.com/articles/getting-started-with-go#introduction) tutorial to get a feel for it.
   I already have a project set up, will add you guys when you give me account info.

3. We will be using Amazon RDS to host our db, for now we can initialize local databases
   with data using our little [init](https://github.com/jaden-young/murder-hobos/tree/master/db/initDb) tool. 
   Note that you should create a database using mysql first, then pass that db/username/password to the init program as cli arguments.
   
   To read up on databases in Go, this [tutorial](http://go-database-sql.org/) is a good place to start.
   We will be using [sqlx](https://github.com/jmoiron/sqlx), a superset of database/sql commands to make things easier on ourselves.
   A quick guide to sqlx can be found [here](http://jmoiron.github.io/sqlx/)

4. Run ```go get -u github.com/jaden-young/murder-hobos``` to download this repository into the directory ```$GOPATH/src/github.com/jaden-young/murder-hobos```
   This is where the project should reside, and where all work should be done. The entire repository is cloned down, including branches.

5. Run ```go install ./...``` from the project root to install server/init executables.

5. Create a .env file in the project root with the following entries:
    ```
    PORT="some-port"
    MYSQL_USER="db-username"
    MYSQL_PASS="db-password"
    MYSQL_DB_NAME="db-database name"
    MYSQL_ADDR="hostname:port"
    ```
    ```MYSQL_ADDR``` for testing is probably going to be ```localhost:3306```, local machine with default mysql port.
    ```PORT``` is the port the webserver will run on, so when you run the server the site can be accessed through typing for example ```localhost:8000```.
    This file is sourced by heroku when running the server, adding those values as environment variables while the server is running.

6. To run the server, run ```heroku local``` in the project root.
7. After making any changes to files, run ```go install ./...``` from the project root to update the executable that the server runs

