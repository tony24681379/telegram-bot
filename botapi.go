package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	"github.com/mmcdole/gofeed"
	"github.com/mrd0ll4r/tbotapi"
)

//Bot struct
type Bot struct {
	Topics        []Topic
	SubscribeList map[string]time.Time
	Commands      map[string]string
}

//RSSData struct
type RSSData struct {
	Title       string
	Description string
	PublishDate time.Time
}

//UpdateSubscribeList pull the RSS and send it
func (b *Bot) UpdateSubscribeList(db *gorm.DB, api *tbotapi.TelegramBotAPI) {
	for _, topic := range b.Topics {
		rssData := b.getRSSData(topic)
		if b.SubscribeList[topic.Topic+topic.Language] != rssData.PublishDate {
			topic.PublishDate = rssData.PublishDate
			var updateTopic Topic
			if b.sendUpdateToSubscribeUser(db, api, topic.Topic, topic.Language, rssData) {
				updateTopic.UpdatePublishDate(db, rssData)
			}
			b.SubscribeList[topic.Topic+topic.Language] = rssData.PublishDate
		}
	}
}

//sendUpdateToSubscribeUser send message to user who subscribe the topics
func (b *Bot) sendUpdateToSubscribeUser(db *gorm.DB, api *tbotapi.TelegramBotAPI, topic, language string, rssData *RSSData) bool {
	var users Users
	var subscribe Subscribe
	if users.FindUpdateSubscribeUser(db, topic, language, rssData) {
		for _, user := range users {
			outMsg, err := api.NewOutgoingMessage(tbotapi.NewChatRecipient(user.UserID), rssData.Title+rssData.Description).SetHTML(true).Send()
			if err != nil {
				fmt.Printf("Error sending: %s\n", err)
				return false
			}
			fmt.Printf("->%d, To:\t%s, Text: %s\n", outMsg.Message.ID, outMsg.Message.Chat, *outMsg.Message.Text)
			subscribe.UpdatePublishDate(db, user, rssData, topic)
		}
	}
	return true
}

//InitBotStatus init status
func (b *Bot) InitBotStatus(db *gorm.DB) {
	var topics []Topic
	db.Find(&topics)
	for _, topic := range topics {
		b.SubscribeList[topic.Topic+topic.Language] = topic.PublishDate
		b.Topics = append(b.Topics, topic)
	}
}

func (b *Bot) getRSSData(topic Topic) *RSSData {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(topic.Link)

	publishDate := time.Now()
	for _, item := range feed.Items {
		if item.Published != "" {
			publishDate = *item.PublishedParsed
		}
	}

	doc, err := goquery.NewDocument(topic.Link)
	if err != nil {
		log.Fatal(err)
	}
	description := doc.Find("description").Text()
	description = strings.Replace(description, "]]>", "", -1)

	rssData := &RSSData{
		Title:       feed.Title,
		Description: description,
		PublishDate: publishDate,
	}
	return rssData
}

//ListTopics list topics
func (b *Bot) ListTopics(db *gorm.DB, userID int) string {
	var user User
	var topics Topics
	user.CheckUser(db, userID)
	topics.ListTopic(db, user.Language)

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
func (b *Bot) Tellme(db *gorm.DB, userID int, args ...string) string {
	if len(args) == 0 {
		return "please enter the topic"
	}
	var user User
	var topic Topic
	user.CheckUser(db, userID)

	isExist := topic.FindTopic(db, user, args[0])
	if isExist {
		rssData := b.getRSSData(topic)
		return rssData.Title + rssData.Description
	}
	return "topic not found"
}

//Subscribe the topic
func (b *Bot) Subscribe(db *gorm.DB, userID int, args ...string) string {
	if len(args) == 0 {
		return "please enter the topic"
	}
	var user User
	var topic Topic
	var subscribe Subscribe
	user.CheckUser(db, userID)

	isExist := topic.FindTopic(db, user, args[0])
	if isExist {
		rssData := b.getRSSData(topic)
		subscribe.Subscribe(db, userID, topic, rssData)
		return user.Answer
	}
	return "error"
}

//UnSubscribe the topic
func (b *Bot) UnSubscribe(db *gorm.DB, userID int, args ...string) string {
	if len(args) == 0 {
		return "please enter the topic"
	}
	var user User
	var topic Topic
	var subscribe Subscribe

	user.CheckUser(db, userID)

	isExist := topic.FindTopic(db, user, args[0])
	if isExist {
		subscribe.UnSubscribe(db, userID, topic)
		return user.Answer
	}
	return "error"
}

//TChinese changes language to T_chinese
func (b *Bot) TChinese(db *gorm.DB, userID int) string {
	var user User
	user.UpdateUser(db, userID, "TChinese", "知道了")
	return "繁體中文"
}

//SChinese changes language to S_chinese
func (b *Bot) SChinese(db *gorm.DB, userID int) string {
	var user User
	user.UpdateUser(db, userID, "SChinese", "知")
	return "简体中文"
}

//English changes language to English
func (b *Bot) English(db *gorm.DB, userID int) string {
	var user User
	user.UpdateUser(db, userID, "English", "OK")
	return "English"
}
