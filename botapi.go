package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	"github.com/mmcdole/gofeed"
)

//Bot struct
type Bot struct {
	SubscribeList map[string]time.Time
	Commands      map[string]string
}

//updateUser check user Language and SubscribeList
func (b Bot) updateUser(db *gorm.DB, userID int, language, answer string) User {
	var user User
	db.Model(&user).Where(User{UserID: userID}).
		Updates(User{UserID: userID, Language: language, Answer: answer})
	return user
}

//UpdateSubscribeList
func (b Bot) UpdateSubscribeList(db *gorm.DB) {
	var topics []Topic
	db.Find(&topics)
	for _, topic := range topics {
		fp := gofeed.NewParser()
		feed, _ := fp.ParseURL(topic.Link)
		b.SubscribeList[topic.Topic+topic.Language] = *feed.PublishedParsed
		if b.SubscribeList[topic.Topic+topic.Language] != topic.PublishDate {
			topic.PublishDate = *feed.PublishedParsed
		}
	}
}

//InitSubscribeList
func (b Bot) InitSubscribeList(db *gorm.DB) {
	var topics []Topic
	db.Find(&topics)
	for _, topic := range topics {
		b.SubscribeList[topic.Topic+topic.Language] = topic.PublishDate
	}
}

//checkUser
func (b Bot) checkUser(db *gorm.DB, userID int) User {
	var user User
	if db.Where(User{UserID: userID}).First(&user).RecordNotFound() {
		user = User{UserID: userID, Language: "English", Answer: "OK"}
		db.Create(&user)
	}
	return user
}

func (b Bot) findTopic(db *gorm.DB, user User, topicName string) (Topic, bool) {
	var topic Topic
	if db.Where("topic = ? AND language = ?", topicName, user.Language).Find(&topic).RecordNotFound() {
		return topic, false
	}
	return topic, true
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
	if len(args) == 0 {
		return "please enter the topic"
	}
	result := ""
	user := b.checkUser(db, userID)
	topic, isExist := b.findTopic(db, user, args[0])

	if isExist {
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
	return "topic not found"
}

//Subscribe the topic
func (b Bot) Subscribe(db *gorm.DB, userID int, args ...string) string {
	if len(args) == 0 {
		return "please enter the topic"
	}
	user := b.checkUser(db, userID)
	topic, isExist := b.findTopic(db, user, args[0])
	if isExist {
		fp := gofeed.NewParser()
		feed, _ := fp.ParseURL(topic.Link)
		publishDate := time.Now()
		for _, item := range feed.Items {
			if item.Published != "" {
				fmt.Println(feed.Published)
				publishDate = *item.PublishedParsed
			}
		}
		var subscribe Subscribe
		db.Where(Subscribe{TopicID: topic.Topic}).
			Attrs(Subscribe{User: userID, TopicID: topic.Topic, PublishDate: publishDate}).
			FirstOrCreate(&subscribe)
		return user.Answer
	}
	return "error"
}

//UnSubscribe the topic
func (b Bot) UnSubscribe(db *gorm.DB, userID int, args ...string) string {
	if len(args) == 0 {
		return "please enter the topic"
	}
	user := b.checkUser(db, userID)
	topic, isExist := b.findTopic(db, user, args[0])
	if isExist {
		db.Unscoped().Where(Subscribe{User: userID, TopicID: topic.Topic}).Delete(Subscribe{})
		return user.Answer
	}
	return "error"
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
