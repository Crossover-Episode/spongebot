package commandhandler

import (
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

const usage = "to return a text only message, @ this bot in a reply. To return a meme, add **--meme** when you mention the bot."

func (h *Handler) OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !mentionsContainsUser(m.Mentions, s.State.User.ID) {
		return
	}

	if m.Type == discordgo.MessageTypeReply && m.MessageReference != nil {
		msgRef, err := s.ChannelMessage(m.MessageReference.ChannelID, m.MessageReference.MessageID)
		if err != nil {
			_, _ = s.ChannelMessageSendReply(m.MessageReference.ChannelID, replyError, m.MessageReference)
			log.Println(err.Error())
		}

		if msgRef.Author.ID == s.State.User.ID {
			return
		}

		contentWithUsernames := replaceUserIdsWithNames(msgRef.Mentions, msgRef.Content)

		if strings.Contains(m.Content, "--meme") {
			err = h.replySpongebobMeme(s, m.MessageReference, contentWithUsernames)
		} else {
			err = h.replySpongebobText(s, m.MessageReference, contentWithUsernames)
		}
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		_, _ = s.ChannelMessageSendReply(m.ChannelID, usage, m.Reference())
	}
}

func (h *Handler) replySpongebobText(s *discordgo.Session, msgRef *discordgo.MessageReference, content string) error {
	_, err := s.ChannelMessageSendReply(msgRef.ChannelID, spongebob.ToText(content, false), msgRef)
	return err
}

func (h *Handler) replySpongebobMeme(s *discordgo.Session, msgRef *discordgo.MessageReference, content string) error {
	memeBuffer, err := h.memeGenerator.GenerateMeme(content)
	if err != nil {
		_, _ = s.ChannelMessageSendReply(msgRef.ChannelID, replyError, msgRef)
		log.Println("Failed to create meme")
	}

	msgWithMeme := &discordgo.MessageSend{
		File: &discordgo.File{
			Name:        makeFileName(content),
			ContentType: "image/jpeg",
			Reader:      memeBuffer,
		},
		Reference: msgRef,
	}
	_, err = s.ChannelMessageSendComplex(msgRef.ChannelID, msgWithMeme)

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
