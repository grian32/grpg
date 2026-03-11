module server

go 1.24.4

require (
	github.com/golang-migrate/migrate/v4 v4.18.3
	github.com/mattn/go-sqlite3 v1.14.30
	grpg/data-go v0.0.0
)

require (
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	go.uber.org/atomic v1.11.0 // indirect
)

replace grpg/data-go => ../data-go
