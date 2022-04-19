package main

import (
	"encoding/json"
	"github.com/slack-go/slack"
	"github.com/traPtitech/traq-ws-bot"
	"github.com/traPtitech/traq-ws-bot/payload"
	"log"
	"os"
)

func main() {
	b, err := traqwsbot.NewBot(&traqwsbot.Options{
		AccessToken: os.Getenv("TRAQ_ACCESS_TOKEN"),
		Origin:      os.Getenv("TRAQ_ORIGIN"),
	})

	if err != nil {
		panic(err)
	}

	api := slack.New(os.Getenv("SLACK_TOKEN"))

	b.OnMessageCreated(func(p *payload.MessageCreated) {
		msg := p.Message
		u := msg.User
		_, _, err := api.PostMessage(os.Getenv("SLACK_CHANNEL_ID"), slack.MsgOptionText(msg.Text, true), slack.MsgOptionUsername(u.DisplayName))
		if err != nil {
			log.Printf("Slack message send error: %+v, message: %+v", err, msg)
		}

	})
	b.OnError(func(message string) {
		log.Println("Error", message)
	})

	if err := b.Start(); err != nil {
		panic(err)
	}
}
