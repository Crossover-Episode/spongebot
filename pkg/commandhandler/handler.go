package commandhandler

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/davidparks11/spongebot/pkg/spongebob"
)


const replyError = "oH nO! sOmEtHiNg WeNt WrOnG! tRy AgAiN lAtEr."

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
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
			s.ChannelMessageSendReply(m.MessageReference.ChannelID, spongebob.ToText(refMsg.Content, false), m.MessageReference)
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
