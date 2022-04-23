package main

import (
	"context"
	"github.com/slack-go/slack/socketmode"
	"log"
	"os"
	"fmt"

	"github.com/slack-go/slack"
	traq "github.com/traPtitech/go-traq"
)

func main() {
	appToken := os.Getenv("SLACK_APP_TOKEN")
	botToken := os.Getenv("SLACK_BOT_TOKEN")

	api := slack.New(
		botToken,
		slack.OptionAppLevelToken(appToken),
	)

	client := socketmode.New(api)

	config := traq.NewConfiguration()
	config.Host = os.Getenv("TRAQ_HOST")
	config.Scheme = os.Getenv("TRAQ_HOST_SCHEME")
	traqclient := traq.NewAPIClient(config)
	go func() {
		for evt := range client.Events {
			switch evt.Type {
			case socketmode.EventTypeSlashCommand:
				cmd, ok := evt.Data.(slack.SlashCommand)

				if !ok {
					log.Printf("Ignore slash command: %+v\n", evt)
					continue
				}

				client.Ack(*evt.Request, "hoge")

				text := fmt.Sprintf("%s: %s", cmd.UserName, cmd.Text)

				postmesreq := traq.NewPostMessageRequest(text)

				str := os.Getenv("TRAQ_ACCESS_TOKEN")

				traqauth := context.WithValue(context.Background(), traq.ContextAccessToken, str)
				_, _, err := traqclient.MessageApi.PostMessage(traqauth, os.Getenv("TRAQ_CHANNEL_ID")).PostMessageRequest(*postmesreq).Execute()
				if err != nil {
					log.Printf("%+v", err)
					panic(err)
				}

			default:
				log.Printf("Event: %+v", evt)
			}
		}
	}()

	client.Run()
}
