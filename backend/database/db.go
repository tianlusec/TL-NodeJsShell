package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func NewDB() *DB {
	db, err := gorm.Open(sqlite.Open("data/nodeshell.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	
	db.AutoMigrate(&Shell{}, &History{}, &Proxy{})
	
	return &DB{db}
}

