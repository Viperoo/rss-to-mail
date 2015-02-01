package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type Website struct {
	Id          int16
	Title       string
	Description string
	Language    string
	Url         string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

type Feed struct {
	Id        int16
	Website   Website
	WebsiteId int16
	GUID      string
}
type Date string

var db gorm.DB

func conn() (db gorm.DB) {
	db, err := gorm.Open("sqlite3", workdir+"/db/rss-to-mail.db")

	if err != nil {
		logger.Critical(err.Error())
	}

	return db
}

func makeMigrate(db gorm.DB) {
	db.DropTableIfExists(&Website{})
	db.DropTableIfExists(&Feed{})

	db.CreateTable(&Website{})
	db.CreateTable(&Feed{})

	logger.Info("Migration complete")
}
