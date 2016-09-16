package main

import (
	"github.com/jinzhu/gorm"
)

type Topic struct {
	gorm.Model
	Topic        string
	EnglishLink  string
	TChineseLink string
	SChineseLink string
}

type User struct {
	gorm.Model
	User      uint
	Language  string
	Subscribe Subscribe
}

type Subscribe struct {
	gorm.Model
	User    uint
	TopicID string
}

func InitDB(db *gorm.DB) {
	db.AutoMigrate(&Topic{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Subscribe{})
}

func createData(db *gorm.DB) {
	db.Create(&Topic{Topic: "current",
		EnglishLink:  "http://rss.weather.gov.hk/rss/CurrentWeather.xml",
		TChineseLink: "http://rss.weather.gov.hk/rss/CurrentWeather_uc.xml",
		SChineseLink: "http://gbrss.weather.gov.hk/rss/CurrentWeather_uc.xml"})
	db.Create(&Topic{Topic: "warning",
		EnglishLink:  "http://rss.weather.gov.hk/rss/WeatherWarningBulletin.xml",
		TChineseLink: "http://rss.weather.gov.hk/rss/WeatherWarningBulletin_uc.xml",
		SChineseLink: "http://gbrss.weather.gov.hk/rss/WeatherWarningBulletin_uc.xml"})
}
