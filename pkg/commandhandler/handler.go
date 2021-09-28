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

const (
	replyError = "oH nO! sOmEtHiNg WeNt WrOnG! tRy AgAiN lAtEr."
	usage      = "to return a text only message, @ this bot in a reply. To return a meme, add **--meme** when you mention the bot."
	helpReminder = "to get the usage for this bot, @ it in a channel."
)

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

		contentWithUserNames, err := msgRef.ContentWithMoreMentionsReplaced(s)
		if err != nil {
			log.Println(err.Error())
			return
		}

		if strings.Contains(m.Content, "--meme") {
			err = h.replySpongebobMeme(s, m.MessageReference, contentWithUserNames)
		} else {
			err = h.replySpongebobText(s, m.MessageReference, contentWithUserNames)
		}
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		err := h.directMessageUsage(s, m)
		if err != nil {
			log.Println(err.Error())
		}
	}
}


func (h *Handler) replySpongebobText(s *discordgo.Session, msgRef *discordgo.MessageReference, content string) error {
	spongebobText := spongebob.ToText(content, false) + "\n\n" + helpReminder
	_, err := s.ChannelMessageSendReply(msgRef.ChannelID, spongebobText, msgRef)
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
		Content: helpReminder,
	}
	_, err = s.ChannelMessageSendComplex(msgRef.ChannelID, msgWithMeme)

	return err
}


func (h *Handler) directMessageUsage(s *discordgo.Session, m *discordgo.MessageCreate) error {
	directMsgChannel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		_, _ = s.ChannelMessageSendReply(m.ChannelID, usage, m.Reference()) //on failed dm, send as reply in original channel
		return err
	}
	_, _ = s.ChannelMessageSend(directMsgChannel.ID, usage)
	return nil
}


const fileNameMaxLength = 20

func makeFileName(content string) string {
	if len(content) > fileNameMaxLength {
		content = content[:fileNameMaxLength-3] + "..."
	}
	return spongebob.ToText(content, false) + ".jpg"
}


func mentionsContainsUser(users []*discordgo.User, value string) bool {
	for _, v := range users {
		if v.ID == value {
			return true
		}
	}
	return false
}
