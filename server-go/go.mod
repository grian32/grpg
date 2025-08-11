module server

go 1.24.4

require grpg/data-go v0.0.0

require grpgscript v0.0.0

require (
	github.com/golang-migrate/migrate/v4 v4.18.3 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.30 // indirect
	go.uber.org/atomic v1.11.0 // indirect
)

replace grpg/data-go => ../data-go

replace grpgscript => ../grpgscript
