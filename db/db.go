package db

import "encore.dev/storage/sqldb"

var db = sqldb.NewDatabase("database", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
