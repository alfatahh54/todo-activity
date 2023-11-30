# Todo-Activity
This is application CRUD Todo-Activity with gin Golang and use database MySql
## Database mysql migration
* Set environment on file ".env" like example on file ".env.example". 
* Install golang-migrate cmd 
``` bash
$ # Go 1.15 and below
$ go get -tags 'mysql' -u github.com/golang-migrate/migrate/cmd/migrate
$ # Go 1.16+
$ go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```
* Migrate up database
``` bash
$ # For up all
$ make migrateup
$ # For up n version migration
$ make migrateup n=total_version
```
* Migrate down database
``` bash
$ # For down all
$ make migratedown
$ # For down n version migration
$ make migrateup n=total_version
```
* Migrate create
* Create new file migration
``` bash
$ make migratecreate name=file_name
```
## Development
* Install gin-bin
``` bash
$ go install github.com/codegangsta/gin@latest
```
* run dev
``` bash
$ make dev
```
## Build
``` bash
$ make build
```
# Start
``` bash
$ make start
```