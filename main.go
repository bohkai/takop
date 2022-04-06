package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type App struct {
	Config *DiscordConfig
	Session *discordgo.Session
}

func New(config *DiscordConfig) (*App, error) {
	dg, err := discordgo.New("Bot " + config.Token);
	if err != nil {
		return nil, err;
	}

	return &App{
		Config: config,
		Session: dg,
	}, nil
}

func (a *App) Open() error {
	 return a.Session.Open()
}

func (a *App) Close() error {
	return a.Session.Close()
}

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	app, err := New(config)
	if err != nil {
		log.Fatal(err)
	}

	channel, err := NewChannel()
	if err != nil {
		log.Fatal(err)
	}

	err = app.Open()
	if err != nil {
		log.Fatal(err)
	}

	app.Session.AddHandler(channel.Join)
	app.Session.AddHandler(channel.Leave)
	app.Session.AddHandler(channel.List)
	app.Session.AddHandler(channel.Play)

	log.Println("Bot is now running.  Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	err = app.Close()
	if err != nil {
		log.Fatal(err)
	}
}

