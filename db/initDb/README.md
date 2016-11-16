# InitDb

This package provides a command ```murder-hobos-init-db``` that initializes our database
to a base state. In this state all spells and classes from PHB, EE, and SCAG are included
with necessary relationships between them.

This exists essentially to parse our magic xml file that we found. Once we have achieved inital
data population, a mysqldump file will be much more efficient for creating this inital state.

To use this command, either ```go build``` in this directory and run the produced executable,
or ```go install``` to have the program installed to your ```$GOBIN```

***WARNING:*** Running this command wipes the tables. 

Usage:
```
murder-hobos-init-db -D database-name -u username -p password -h hostname -P port
```