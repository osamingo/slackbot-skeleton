package event

import (
	"context"
	"fmt"

	"github.com/go-joe/joe"
	"github.com/slack-go/slack"
)

func NewPingExecution(channel string) Execution {
	e := Execution{
		eventName: "PING",
		channel:   channel,
	}

	e.fn = func(ctx context.Context, body []byte) (string, []slack.MsgOption, error) {
		return e.channel, []slack.MsgOption{
			slack.MsgOptionText("pong", false),
		}, nil
	}

	return e
}

func PingRespond(_ *joe.Bot) func(joe.Message) error {
	return func(msg joe.Message) error {
		if err := msg.RespondE("pong"); err != nil {
			return fmt.Errorf("event: failed to send a response: %w", err)
		}

		return nil
	}
}
