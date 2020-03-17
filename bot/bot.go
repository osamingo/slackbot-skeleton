package bot

import (
	"context"
	"time"

	joehttp "github.com/go-joe/http-server"
	"github.com/go-joe/joe"
	slackadpt "github.com/go-joe/slack-adapter"
	"github.com/nlopes/slack"
	"github.com/osamingo/slackbot/event"
	"github.com/pkg/errors"
	stackdriver "github.com/tommy351/zap-stackdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// A Bot wraps joe.Bot.
type Bot struct {
	bot        *joe.Bot
	slack      *slack.Client
	router     *event.Router
	routingKey string
}

// NewBot generates bot.Bot.
func NewBot(name, slackToken, path string, timeout time.Duration, debug bool) (*Bot, error) {

	lvl := zapcore.InfoLevel
	if debug {
		lvl = zapcore.DebugLevel
	}

	config := &zap.Config{
		Level:            zap.NewAtomicLevelAt(lvl),
		Encoding:         "json",
		EncoderConfig:    stackdriver.EncoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return &stackdriver.Core{
			Core: core,
		}
	}))
	if err != nil {
		return nil, err
	}

	return &Bot{
		bot: joe.New(
			name,
			slackadpt.Adapter(slackToken),
			joehttp.Server(path),
			joe.WithHandlerTimeout(timeout),
			joe.WithLogger(logger),
		),
		slack: slack.New(slackToken),
	}, nil
}

// SetRouter sets event router.
func (b *Bot) SetRouter(routingKey string, r *event.Router) {
	b.routingKey = routingKey
	b.router = r
}

// SetRespondRegex sets respond with expression.
func (b *Bot) SetRespondRegex(expr string, f func(*joe.Bot) func(joe.Message) error) {
	b.bot.RespondRegex(expr, f(b.bot))
}

// Run starts the Bot.
func (b *Bot) Run() error {
	b.bot.Brain.RegisterHandler(b.HandleHTTP)
	return b.bot.Run()
}

// HandleHTTP routes HTTP requests.
func (b *Bot) HandleHTTP(ctx context.Context, r joehttp.RequestEvent) error {

	switch r.URL.Path {
	case "/_ah/warmup":
		b.bot.Logger.Info("catch warm up request")

	case "/_events":

		eventName := r.Header.Get(b.routingKey)
		f := b.router.GetFunc(eventName)
		if f == nil {
			return errors.Errorf("bot: not found an event, event_name = %s", eventName)
		}

		target, opts, err := f(ctx, r.Body)
		if err != nil {
			return errors.Wrap(err, "bot: failed to execute event")
		}

		_, _, err = b.slack.PostMessageContext(ctx, target, opts...)
		if err != nil {
			return errors.Wrap(err, "bot: failed to send to slack")
		}
	}

	return nil
}
