// Copyright 2015-2016 mrd0ll4r and contributors. All rights reserved.
// Use of this source code is governed by the MIT license, which can be found in
// the LICENSE file.

package main

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/mrd0ll4r/tbotapi"
	"github.com/mrd0ll4r/tbotapi/examples/boilerplate"
)

func main() {
	apiToken := "260002715:AAE6BGznYNTeLN-3V8pz1XwfOKsYoa8_nV4"

	// Note: For this example to work, you'll have to enable inline queries for your bot (chat with @BotFather).

	b := &Bot{
		Language:      "English",
		SubscribeList: make(map[string]bool),
	}
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
			// msg.Chat implements fmt.Stringer, so it'll display nicely.
			// We know it's a text message, so we can safely use the Message.Text pointer.
			fmt.Printf("<-%d, From:\t%s, Text: %s \n", msg.ID, msg.Chat, *msg.Text)

			msgs := strings.Split(*msg.Text, " ")
			methodName := msgs[0]
			method := reflect.ValueOf(b).MethodByName(methodName)
			if method.IsValid() {
				exec := method.Interface().(func(...string) string)
				outMsg, err := api.NewOutgoingMessage(tbotapi.NewRecipientFromChat(msg.Chat), exec(msgs[1:]...)).Send()
				if err != nil {
					fmt.Printf("Error sending: %s\n", err)
					return
				}
				fmt.Printf("->%d, To:\t%s, Text: %s\n", outMsg.Message.ID, outMsg.Message.Chat, *outMsg.Message.Text)
				return
			}

			outMsg, err := api.NewOutgoingMessage(tbotapi.NewRecipientFromChat(msg.Chat), *msg.Text).Send()

			if err != nil {
				fmt.Printf("Error sending: %s\n", err)
				return
			}
			fmt.Printf("->%d, To:\t%s, Text: %s\n", outMsg.Message.ID, outMsg.Message.Chat, *outMsg.Message.Text)
		case tbotapi.InlineQueryUpdate:
			query := update.InlineQuery
			fmt.Printf("<-%s (query), From:\t%s, Query: %s \n", query.ID, query.From, query.Query)
			var results []tbotapi.InlineQueryResult

			for i, s := range query.Query {
				if len(results) >= 50 {
					// The API accepts up to 50 results.
					break
				}
				if !unicode.IsSpace(s) {
					// Don't set mandatory fields to whitespace.
					results = append(results, tbotapi.NewInlineQueryResultArticle(fmt.Sprint(i), string(s), string(s)))
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
