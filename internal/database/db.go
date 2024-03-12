package database

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const defaultTimeout = 3 * time.Second

type DB struct {
	*gorm.DB
}

func New(dsn string) (*DB, error) {
	gormDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db, err := gormDb.DB()
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(2 * time.Hour)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DB{gormDb}, nil
}
