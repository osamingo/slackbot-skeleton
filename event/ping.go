package event

import (
	"context"

	"github.com/go-joe/joe"
	"github.com/nlopes/slack"
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
		return msg.RespondE("pong")
	}
}
