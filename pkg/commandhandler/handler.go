package commandhandler

import (
	"bytes"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/davidparks11/spongebot/pkg/spongebob"
)

type Handler struct {
	memeGenerator *spongebob.MemeGenerator
}

func NewHandler(memeGenerator *spongebob.MemeGenerator) *Handler {
	return &Handler{memeGenerator: memeGenerator}
}

const replyError = "oH nO! sOmEtHiNg WeNt WrOnG! tRy AgAiN lAtEr."

func (h *Handler) OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	//check if message
	if m.Type == discordgo.MessageTypeReply && m.MessageReference != nil {
		h.replySpongebobText(s, m)
	}
}

func (h *Handler) replySpongebobText(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !mentionsContainsUser(m.Mentions, s.State.User.ID) {
		return nil
	}

	refMsg, err := s.ChannelMessage(m.MessageReference.ChannelID, m.MessageReference.MessageID)
	if err != nil {
		s.ChannelMessageSendReply(m.MessageReference.ChannelID, replyError, m.MessageReference)
		return err
	}
	
	if refMsg.Author.ID != s.State.User.ID {
		contentWithUsernames := replaceUserIdsWithNames(refMsg.Mentions, refMsg.Content)
		
		var memeBuffer *bytes.Buffer
		memeBuffer, err = h.memeGenerator.GenerateMeme(contentWithUsernames)
		if err != nil {
			_, _ = s.ChannelMessageSendReply(m.MessageReference.ChannelID, replyError, m.MessageReference)
			log.Println("Failed to create meme")
		}

		msgWithMeme := &discordgo.MessageSend{
			File: &discordgo.File{
				Name:        makeFileName(refMsg.Content),
				ContentType: "image/jpeg",
				Reader:      memeBuffer,
			},
			Reference: m.MessageReference,
		}
		_, err = s.ChannelMessageSendComplex(m.MessageReference.ChannelID, msgWithMeme)
	}
	return err
}

const fileNameMaxLength = 20

func makeFileName(content string) string {
	if len(content) > fileNameMaxLength {
		content = content[:fileNameMaxLength-3] + "..."
	}
	return spongebob.ToText(content, false) + ".jpg"
}

func replaceUserIdsWithNames(users []*discordgo.User, content string) string {
	for _, user := range users {
		content = strings.ReplaceAll(content, "<@!"+user.ID+">", "@"+user.Username)
	}
	return content
}

func mentionsContainsUser(users []*discordgo.User, value string) bool {
	for _, v := range users {
		if v.ID == value {
			return true
		}
	}
	return false
}
