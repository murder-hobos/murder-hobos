# murder-hobos

Murder Hobos is a DnD 5e Spellbook reference application... which may be extended as far as we have the will to.

## Getting Started

1. If you're keen on using windows to develop, follow the instructions on the golang 
   [webpage](https://golang.org/doc/install), using the [windows installer](https://storage.googleapis.com/golang/go1.7.3.windows-amd64.msi)
   If you're using a linux distro, go will be in the standard repos somewhere. Make sure to setup your path.

2. Follow the [heroku](https://devcenter.heroku.com/articles/getting-started-with-go#introduction) tutorial to get a feel for it.
   I already have a project set up, will add you guys when you give me account info.

3. We will be using Amazon RDS to host our db, for now we can initialize local databases
   with data using our little [init](https://github.com/jaden-young/murder-hobos/tree/master/db/initDb) tool. To read up on databases in Go, this [tutorial](http://go-database-sql.org/) is a good place to start.
   We will be using [sqlx](https://github.com/jmoiron/sqlx), a superset of database/sql commands to make things easier on ourselves.
   A quick guide to sqlx can be found [here](http://jmoiron.github.io/sqlx/)