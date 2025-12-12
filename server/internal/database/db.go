package database

import (
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func NewDB() *DB {
	if err := os.MkdirAll("data", 0755); err != nil {
		// Try creating in server/data if running from root
		if err := os.MkdirAll("server/data", 0755); err == nil {
			// If successful, use server/data path
		}
	}

	dbPath := "data/nodeshell.db"
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		if _, err := os.Stat("server/data"); err == nil {
			dbPath = "server/data/nodeshell.db"
		}
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Shell{}, &History{}, &Proxy{})

	return &DB{db}
}
