package main

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
)

//Bot struct
type Bot struct {
	Language      string
	SubscribeList map[string]bool
}

//Topics list topics
func (b Bot) Topics(args ...string) string {
	root, err := goquery.NewDocument("https://data.gov.hk/en-data/provider/hk-hko")
	if err != nil {
		panic(err)
	}
	titles := ""
	root.Find("div.dataset-item").Find("h3.dataset-heading").Each(func(i int, s *goquery.Selection) {
		title := s.Find("a").Text() + "\n"
		titles += title
	})
	return titles
}

//Tellme shows topic
func (b Bot) Tellme(args ...string) string {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("http://rss.weather.gov.hk/rss/SeveralDaysWeatherForecast.xml")
	result := ""
	for _, v := range feed.Items {
		result += v.Description
		fmt.Println(feed.Items)
	}
	return result
}

//Subscribe the topic
func (b Bot) Subscribe(args ...string) string {
	return "a"
}

//UnSubscribe the topic
func (b Bot) UnSubscribe(args ...string) string {
	return "a"
}

//TChinese changes language to T_chinese
func (b Bot) TChinese(args ...string) string {
	return "a"
}

//SChinese changes language to S_chinese
func (b Bot) SChinese(args ...string) string {
	return "a"
}

//English changes language to English
func (b Bot) English() string {
	return "a"
}
