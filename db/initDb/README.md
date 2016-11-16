# InitDb

This package provides a command ```murder-hobos-init-db``` that initializes our database
to an initial state. In this state all spells and classes from PHB, EE, and SCAG are included
with necessary relationships between them.

To use this command, either ```go build``` in this directory and run the produced executable,
or ```go install``` to have the program installed to your ```$GOBIN```

***WARNING:*** Running this command wipes the tables. 

Usage:
```
murder-hobos-init-db -D database-name -u username -p password -h hostname -P port
```