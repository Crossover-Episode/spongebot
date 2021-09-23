package main

import (
	// "log"
	// "os"
	// "os/signal"
	// "syscall"

	// "github.com/bwmarrin/discordgo"
	// "github.com/davidparks11/spongebot/pkg/commandhandler"
	"image/jpeg"
	"os"

	"github.com/davidparks11/spongebot/pkg/spongebob"
)

func main() {
	img, err := spongebob.GenerateMeme("hopefully this works")
	if err != nil {
		panic(err.Error())
	}

	out, err := os.Create("./output.jpg")
	if err != nil {
			panic(err.Error())
	}

	var opt jpeg.Options
	opt.Quality = 100

	err = jpeg.Encode(out, img, &opt) 
	if err != nil {
			panic(err.Error())
	}
	// if len(os.Args) != 2 {
	// 	log.Fatalln("failed to supply bot token")
	// }
	// token := os.Args[1]
	// botSession, err := discordgo.New("Bot " + token)
	// if err != nil {
	// 	log.Fatalln("error creating bot session: " + err.Error())
	// }

	// botSession.AddHandler(commandhandler.OnMessageCreate)

	// err = botSession.Open()
	// if err != nil {
	// 	log.Fatalln("Failed to open bot session: " + err.Error())
	// }

	// log.Println("Bot is now running. Press CTRL-C to exit.")
	// sc := make(chan os.Signal, 1)
	// signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	// <-sc
}
