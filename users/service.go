// Service users keeps track of which users to monitor.
package users

import (
	"encore.dev/storage/sqldb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//encore:service
type Service struct {
	db *gorm.DB
}
type User struct {
	ID             string
	Email          string
	UserName       string
	HashedPassword string
}

// Define a database named 'site', using the database migrations
// in the "./migrations" folder. Encore automatically provisions,
// migrates, and connects to the database.
var db = sqldb.Named("database")

// initService initializes the site service.
// It is automatically called by Encore on service startup.
func initService() (*Service, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db.Stdlib(),
	}))
	if err != nil {
		return nil, err
	}
	return &Service{db: db}, nil
}
