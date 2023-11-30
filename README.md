# Todo-Activity
This is application CRUD Todo-Activity with gin Golang and use database MySql
## Database mysql migrtion
To migrate database please set environment on file ".env" like example on file ".env.example". 
Install golang-migrate cmd 
``` bash
$ # Go 1.15 and below
$ go get -tags 'mysql' -u github.com/golang-migrate/migrate/cmd/migrate
$ # Go 1.16+
$ go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```
For migrate up database please run on bash
``` bash
$ make migrateup
```
For migrate down please run on bash
``` bash
$ make migrateup
```