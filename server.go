// Copyright 2015-2016 mrd0ll4r and contributors. All rights reserved.
// Use of this source code is governed by the MIT license, which can be found in
// the LICENSE file.

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/mrd0ll4r/tbotapi"
	"github.com/mrd0ll4r/tbotapi/examples/boilerplate"
)

func main() {
	apiToken := "260002715:AAE6BGznYNTeLN-3V8pz1XwfOKsYoa8_nV4"
	db, err := gorm.Open("mysql", "root:test@/test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Note: For this example to work, you'll have to enable inline queries for your bot (chat with @BotFather).

	InitDB(db)
	b := &Bot{
		SubscribeList: make(map[string]time.Time),
		Commands: map[string]string{
			"topics":      "List out the topic that our bot support",
			"tellme":      "Echo back the current info in forecast feed",
			"subscribe":   "The bot will send message to the user on new warning",
			"unsubscribe": "The bot will not send update to user from now on",
			"English":     "The bot will send content in English to user",
			"繁體中文":        "The bot will send content in traditional chinese to user",
			"简体中文":        "The bot will send content in simplified chinese to user"},
	}
	b.InitSubscribeList(db)

	updateFunc := func(update tbotapi.Update, api *tbotapi.TelegramBotAPI) {
		switch update.Type() {
		case tbotapi.MessageUpdate:
			msg := update.Message
			typ := msg.Type()
			if typ != tbotapi.TextMessage {
				// Ignore non-text messages for now.
				fmt.Println("Ignoring non-text message")
				return
			}
			// Note: Bots cannot receive from channels, at least no text messages. So we don't have to distinguish anything here.

			// Display the incoming message.
			fmt.Printf("<-%d, From:\t%s, Text: %s \n", msg.ID, msg.Chat, *msg.Text)

			result := ""
			msgs := strings.Split(*msg.Text, " ")

			if len(msgs) > 1 && msgs[0] == "@tony24681379_bot" {
				msgs = msgs[1:]
			}
			methodName := msgs[0]
			if methodName == "topics" {
				result = b.Topics(db, msg.From.ID)
			} else if methodName == "tellme" {
				result = b.Tellme(db, msg.From.ID, msgs[1:]...)
			} else if methodName == "subscribe" {
				result = b.Subscribe(db, msg.From.ID, msgs[1:]...)
			} else if methodName == "unsubscribe" {
				result = b.UnSubscribe(db, msg.From.ID, msgs[1:]...)
			} else if methodName == "English" {
				result = b.English(db, msg.From.ID)
			} else if methodName == "繁體中文" {
				result = b.TChinese(db, msg.From.ID)
			} else if methodName == "简体中文" {
				result = b.SChinese(db, msg.From.ID)
			}

			if result == "" {
				result = *msg.Text
			}
			outMsg, err := api.NewOutgoingMessage(tbotapi.NewRecipientFromChat(msg.Chat), result).SetHTML(true).Send()
			if err != nil {
				fmt.Printf("Error sending: %s\n", err)
				return
			}
			fmt.Printf("->%d, To:\t%s, Text: %s\n", outMsg.Message.ID, outMsg.Message.Chat, *outMsg.Message.Text)
		case tbotapi.InlineQueryUpdate:
			query := update.InlineQuery
			fmt.Printf("<-%s (query), From:\t%s, Query: %s \n", query.ID, query.From, query.Query)
			var results []tbotapi.InlineQueryResult

			if query.Query == "" {
				for command, info := range b.Commands {
					results = append(results, tbotapi.NewInlineQueryResultArticle(command, command, info))
				}
			} else {
				for command, info := range b.Commands {
					if strings.Contains(command, query.Query) {
						results = append(results, tbotapi.NewInlineQueryResultArticle(command, command, info))
					}
				}
			}

			err := api.NewInlineQueryAnswer(query.ID, results).Send()
			if err != nil {
				fmt.Printf("Err: %s\n", err)
			}
		case tbotapi.ChosenInlineResultUpdate:
			// id, not value.
			fmt.Println("Chosen inline query result (ID):", update.ChosenInlineResult.ID)
		default:
			fmt.Println("Ignoring unknown Update type.")
		}
	}

	// Run the bot, this will block.
	boilerplate.RunBot(apiToken, updateFunc, "InlineQuery", "Demonstrates inline queries by splitting words")
}
