package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/davidparks11/spongebot/pkg/commandhandler"
	
	"github.com/davidparks11/spongebot/pkg/spongebob"
)

func main() {

	generator, err:= spongebob.NewGenerator()
	if err != nil {
		panic(err.Error())
	}

	handler := commandhandler.NewHandler(generator)

	if len(os.Args) != 2 {
		log.Fatalln("failed to supply bot token")
	}
	token := os.Args[1]
	botSession, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln("error creating bot session: " + err.Error())
	}

	botSession.AddHandler(handler.OnMessageCreate)

	err = botSession.Open()
	if err != nil {
		log.Fatalln("Failed to open bot session: " + err.Error())
	}

	log.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
