package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"unicode"

	"github.com/bwmarrin/discordgo"
)

const replyError = "oH nO! sOmEtHiNg WeNt WrOnG! tRy AgAiN lAtEr."

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("failed to supply bot token")
	}
	token := os.Args[1]
	botSession, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln("error creating bot session: " + err.Error())
	}

	botSession.AddHandler(onMessageCreate)

	err = botSession.Open()
	if err != nil {
		log.Fatalln("Failed to open bot session: " + err.Error())
	}

	log.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	//check if message
	if m.Type == discordgo.MessageTypeReply && m.MessageReference != nil {
		replySpongebobText(s, m)
	}
}

func replySpongebobText(s *discordgo.Session, m *discordgo.MessageCreate) {
		if mentionsContainsUser(m.Mentions, s.State.User.ID) {
			refMsg, msgErr := s.ChannelMessage(m.MessageReference.ChannelID, m.MessageReference.MessageID)
			if msgErr != nil {
				s.ChannelMessageSendReply(m.MessageReference.ChannelID, replyError, m.MessageReference)
				log.Println("Failed to retrieve message")
			} else if refMsg.Author.ID != s.State.User.ID && !strings.Contains(refMsg.Content, s.State.User.ID) {
				s.ChannelMessageSendReply(m.MessageReference.ChannelID, ToSpongebobText(refMsg.Content, false), m.MessageReference)
			}
		}
}

func mentionsContainsUser(arr []*discordgo.User, value string) bool {
	for _, v := range arr {
		if v.ID == value {
			return true
		}
	}
	return false
}

const (
	firstCapital = 65
	lastCapital  = 90
	firstLower   = 97
	lastLower    = 122
	caseDelta    = firstCapital - firstLower
)

func ToSpongebobText(s string, startCapital bool) string {
	isLower := func(char int32) bool {
		return char >= firstLower && char <= lastLower
	}

	isUpper := func(char int32) bool {
		return char >= firstCapital && char <= lastCapital
	}

	toCapitol := startCapital
	var sbString string
	for _, v := range s {
		if v > unicode.MaxASCII {
			return s
		}
		if isLower(v) || isUpper(v) {
			if toCapitol && isLower(v) {
				sbString += string(v + caseDelta)
			} else if !toCapitol && isUpper(v) {
				sbString += string(v - caseDelta)
			} else {
				sbString += string(v)
			}
			toCapitol = !toCapitol
		} else {
			sbString += string(v)
		}
	}
	return sbString
}
