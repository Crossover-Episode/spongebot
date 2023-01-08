package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/crossoverepisode/spongebot/pkg/commandhandler"

	"github.com/crossoverepisode/spongebot/pkg/spongebob"
)

func main() {

	generator, err := spongebob.NewGenerator()
	if err != nil {
		panic(err.Error())
	}

	if len(os.Args) != 2 {
		log.Fatalln("failed to supply bot token")
	}
	token := os.Args[1]
	botSession, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln("error creating bot session: " + err.Error())
	}
	
	handler := commandhandler.NewHandler(generator)
	handler.RegisterCommandHandler(botSession)

	err = botSession.Open()
	if err != nil {
		log.Fatalln("Failed to open bot session: " + err.Error())
	}

	log.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
