package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/osamingo/slackbot/bot"
	"github.com/osamingo/slackbot/event"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	// Set environment from .env file.
	if err := godotenv.Load(); err != nil {
		return err
	}

	// app engine environment
	env := os.Getenv("NODE_ENV")
	port := os.Getenv("PORT")
	// write env_variables section in app engine settings.
	timeout := os.Getenv("TIMEOUT_SECOND")
	name := os.Getenv("BOT_NAME")
	key := os.Getenv("EVENT_ROUTING_KEY")
	channel := os.Getenv("DEFAULT_SLACK_CHANNEL")
	// token have to set in .env file.
	token := os.Getenv("SLACK_TOKEN")

	const prodEnv = "production"
	flag := env != prodEnv

	sec, err := strconv.Atoi(timeout)
	if err != nil {
		return err
	}

	b, err := bot.NewBot(name, token, ":"+port, time.Duration(sec)*time.Second, flag)
	if err != nil {
		return err
	}

	b.SetRouter(key, event.NewRouter(
		event.NewPingExecution(channel),
	))

	b.SetRespondRegex("^ping$", event.PingRespond)

	return b.Run()
}
