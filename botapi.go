package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	"github.com/mmcdole/gofeed"
)

//Bot struct
type Bot struct {
	Topic    []string
	Language string
}

//CheckUser check user Language and SubscribeList
func (b Bot) CheckUser() {

}

//Topics list topics
func (b Bot) Topics(db *gorm.DB, args ...string) string {
	var topics []Topic
	db.Select("topic").Find(&topics)
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
func (b Bot) Tellme(db *gorm.DB, args ...string) string {
	result := ""

	var topic Topic
	db.Where("topic = ?", args[0]).Find(&topic)
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(topic.TChineseLink)
	result += feed.Title

	doc, err := goquery.NewDocument(topic.TChineseLink)
	if err != nil {
		log.Fatal(err)
	}
	result += doc.Find("description").Text()

	result = strings.Replace(result, "]]>", "", -1)
	fmt.Println(result)
	return result
}

//Subscribe the topic
func (b Bot) Subscribe(db *gorm.DB, args ...string) string {
	return "a"
}

//UnSubscribe the topic
func (b Bot) UnSubscribe(db *gorm.DB, args ...string) string {
	return "a"
}

//TChinese changes language to T_chinese
func (b Bot) TChinese(db *gorm.DB, args ...string) string {
	return "a"
}

//SChinese changes language to S_chinese
func (b Bot) SChinese(db *gorm.DB, args ...string) string {
	return "a"
}

//English changes language to English
func (b Bot) English(db *gorm.DB, args ...string) string {
	return "a"
}
