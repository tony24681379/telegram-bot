package main

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	"github.com/mmcdole/gofeed"
)

//Bot struct
type Bot struct {
	Language string
}

//updateUser check user Language and SubscribeList
func (b Bot) updateUser(db *gorm.DB, userID int, language, answer string) User {
	var user User
	db.Model(&user).Where(User{UserID: userID}).
		Updates(User{UserID: userID, Language: language, Answer: answer})
	return user
}

//checkUser
func (b Bot) checkUser(db *gorm.DB, userID int) User {
	var user User
	if err := db.Where(User{UserID: userID}).First(&user).Error; err != nil {
		user = User{UserID: userID, Language: "English", Answer: "OK"}
		db.Create(&user)
	}
	return user
}

//Topics list topics
func (b Bot) Topics(db *gorm.DB, userID int) string {
	user := b.checkUser(db, userID)

	var topics []Topic
	db.Select("topic").Where("language = ?", user.Language).Find(&topics)
	result := ""
	for i, v := range topics {
		if i == 0 {
			result += v.Topic
		} else {
			result += ", " + v.Topic
		}
	}
	return result
}

//Tellme shows topic
func (b Bot) Tellme(db *gorm.DB, userID int, args ...string) string {
	result := ""

	user := b.checkUser(db, userID)

	var topic Topic
	db.Where("topic = ? AND language = ?", args[0], user.Language).Find(&topic)

	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(topic.Link)
	result += feed.Title

	doc, err := goquery.NewDocument(topic.Link)
	if err != nil {
		log.Fatal(err)
	}
	result += doc.Find("description").Text()

	result = strings.Replace(result, "]]>", "", -1)
	return result
}

//Subscribe the topic
func (b Bot) Subscribe(db *gorm.DB, userID int, args ...string) string {
	user := b.checkUser(db, userID)
	return user.Answer
}

//UnSubscribe the topic
func (b Bot) UnSubscribe(db *gorm.DB, userID int, args ...string) string {
	user := b.checkUser(db, userID)
	return user.Answer
}

//TChinese changes language to T_chinese
func (b Bot) TChinese(db *gorm.DB, userID int) string {
	b.updateUser(db, userID, "TChinese", "知道了")
	return "繁體中文"
}

//SChinese changes language to S_chinese
func (b Bot) SChinese(db *gorm.DB, userID int) string {
	b.updateUser(db, userID, "SChinese", "知")
	return "简体中文"
}

//English changes language to English
func (b Bot) English(db *gorm.DB, userID int) string {
	b.updateUser(db, userID, "English", "OK")
	return "English"
}
