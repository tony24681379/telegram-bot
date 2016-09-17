package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Topic struct {
	gorm.Model
	Language    string
	Topic       string
	Link        string
	PublishDate time.Time
}

type User struct {
	gorm.Model
	UserID    int
	Language  string
	Answer    string
	Subscribe []Subscribe
}

type Subscribe struct {
	gorm.Model
	User        int
	TopicID     string
	PublishDate time.Time
}

func InitDB(db *gorm.DB) {
	db.AutoMigrate(&Topic{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Subscribe{})
	//createData(db)
}

func createData(db *gorm.DB) {
	db.Create(&Topic{
		Language:    "English",
		Topic:       "current",
		Link:        "http://rss.weather.gov.hk/rss/CurrentWeather.xml",
		PublishDate: time.Now()})
	db.Create(&Topic{
		Language:    "TChinese",
		Topic:       "天氣報告",
		Link:        "http://rss.weather.gov.hk/rss/CurrentWeather_uc.xml",
		PublishDate: time.Now()})
	db.Create(&Topic{
		Language:    "SChinese",
		Topic:       "天气报告",
		Link:        "http://gbrss.weather.gov.hk/rss/CurrentWeather_uc.xml",
		PublishDate: time.Now()})

	db.Create(&Topic{
		Language:    "English",
		Topic:       "warning",
		Link:        "http://rss.weather.gov.hk/rss/WeatherWarningBulletin.xml",
		PublishDate: time.Now()})
	db.Create(&Topic{
		Language:    "TChinese",
		Topic:       "警告",
		Link:        "http://rss.weather.gov.hk/rss/WeatherWarningBulletin_uc.xml",
		PublishDate: time.Now()})
	db.Create(&Topic{
		Language:    "SChinese",
		Topic:       "警告",
		Link:        "http://gbrss.weather.gov.hk/rss/WeatherWarningBulletin_uc.xml",
		PublishDate: time.Now()})
}
