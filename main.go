package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"github.com/bwmarrin/discordgo"
)

func main() {
	var c, err = NewConfig();
	if err != nil {
		log.Fatal(err)
		return
	}

	dg, err := discordgo.New("Bot " + c.Token);
	if err != nil {
		log.Println("error creating Discord session,", err)
		return;
	}

	err = dg.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}

	dg.AddHandler(ChannelVoiceJoin);

	log.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

func ChannelVoiceJoin(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}


	for _, g := range s.State.Guilds {
		for _, vs := range g.VoiceStates {
			if m.Author.ID == vs.UserID {
				s.ChannelVoiceJoin(g.ID, vs.ChannelID, false, false)
			}
		}
	}
}

