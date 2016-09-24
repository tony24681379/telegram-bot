package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Topic struct {
	gorm.Model
	Topic       string
	Language    string
	Link        string
	PublishDate time.Time
}

type Topics []Topic

type User struct {
	gorm.Model
	UserID    int
	Language  string
	Answer    string
	Subscribe []Subscribe
}

type Users []User

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

func (topic *Topic) FindTopic(db *gorm.DB, user User, topicName string) bool {
	if db.Where("topic = ? AND language = ?", topicName, user.Language).Find(&topic).RecordNotFound() {
		return false
	}
	return true
}

func (topics *Topics) ListTopic(db *gorm.DB, language string) {
	db.Select("topic").Where("language = ?", language).Find(&topics)
}

func (topic *Topic) UpdatePublishDate(db *gorm.DB, rssData *RSSData) {
	db.Model(&topic).Where(Topic{Topic: topic.Topic, Language: topic.Language}).
		Update(Topic{PublishDate: rssData.PublishDate})
}

func (user *User) UpdateUser(db *gorm.DB, userID int, language, answer string) {
	db.Model(&user).Where(User{UserID: userID}).
		Updates(User{UserID: userID, Language: language, Answer: answer})
}

func (user *User) CheckUser(db *gorm.DB, userID int) {
	if db.Where(User{UserID: userID}).First(&user).RecordNotFound() {
		user = &User{UserID: userID, Language: "English", Answer: "OK"}
		db.Create(&user)
	}
}

func (users *Users) FindUpdateSubscribeUser(db *gorm.DB, topic, language string, rssData *RSSData) bool {
	if !db.Joins("JOIN subscribes s on s.user = user_id").
		Where("topic_id = ? AND language = ? AND publish_date <> ?", topic, language, rssData.PublishDate).
		Find(&users).RecordNotFound() {
		return true
	}
	return false
}

func (subscribe *Subscribe) UpdatePublishDate(db *gorm.DB, user User, rssData *RSSData, topic string) {
	db.Model(&subscribe).Where(Subscribe{User: user.UserID, TopicID: topic}).
		Update(Subscribe{PublishDate: rssData.PublishDate})
}

func (subscribe *Subscribe) Subscribe(db *gorm.DB, userID int, topic Topic, rssData *RSSData) {
	db.Where(Subscribe{TopicID: topic.Topic}).
		Attrs(Subscribe{User: userID, TopicID: topic.Topic, PublishDate: rssData.PublishDate}).
		FirstOrCreate(&subscribe)
}

func (subscribe *Subscribe) UnSubscribe(db *gorm.DB, userID int, topic Topic) {
	db.Unscoped().Where(Subscribe{User: userID, TopicID: topic.Topic}).Delete(Subscribe{})
}
