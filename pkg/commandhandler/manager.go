package commandhandler

import (
	"errors"
	"log"
	// "strings"

	"github.com/bwmarrin/discordgo"
	"github.com/crossoverepisode/spongebot/pkg/spongebob"
)

type HandlerManager struct {
	memeGenerator *spongebob.MemeGenerator
}

func NewHandler(memeGenerator *spongebob.MemeGenerator) *HandlerManager {
	return &HandlerManager{memeGenerator: memeGenerator}
}

const replyError = "oH nO! sOmEtHiNg WeNt WrOnG! tRy AgAiN lAtEr."

var memeCommand = &discordgo.ApplicationCommand{
	Name:        "meme",
	Type:        discordgo.MessageApplicationCommand,
}

func (h *HandlerManager) RegisterCommandHandler(s *discordgo.Session) {
	s.AddHandler(h.onServerJoin)
	s.AddHandler(h.onInteractionCreate)
}

func (h *HandlerManager) onServerJoin(session *discordgo.Session, interaction *discordgo.GuildCreate) {
	_, err := session.ApplicationCommandCreate(session.State.User.ID, interaction.Guild.ID, memeCommand)
	if err != nil {
		log.Fatalf("could not register command %s with guild %s: %v", memeCommand.Name, interaction.Guild.ID, err)
	}
	log.Printf("registered command: %s on guild: %s [%s]\n", memeCommand.Name, interaction.Guild.Name, interaction.Guild.ID)
}

func (h *HandlerManager) onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name != memeCommand.Name {
		return
	}

	err := h.meme(s, i)
	if err != nil {
		log.Printf("error calling %s: %v", memeCommand.Name, err)
	}
}

func (h *HandlerManager) meme(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	cmdData := i.ApplicationCommandData()
	targetMessage := cmdData.Resolved.Messages[cmdData.TargetID]
	if targetMessage.Author.ID == s.State.User.ID {
		return respond(s, i, "yOu CaN't MeMe ThIs BoT lOl\n", errors.New("user tried to call spongebot on spongebot lol"))
	}

	contentWithUserNames, err := targetMessage.ContentWithMoreMentionsReplaced(s)
	if err != nil {
		return respond(s, i, replyError, err)
	}

	err = h.replySpongebobMeme(s, targetMessage.Reference(), contentWithUserNames)
	if err != nil {
		return respond(s, i, replyError, err)
	}
	// err = h.replySpongebobText(s, i.MessageReference, contentWithUserNames)

	return respond(s, i, "your meme has been sent!", nil)
}

// func (h *HandlerManager) replySpongebobText(s *discordgo.Session, msgRef *discordgo.MessageReference, content string) error {
// 	spongebobText := spongebob.ToText(content, false)
// 	_, err := s.ChannelMessageSendReply(msgRef.ChannelID, spongebobText, msgRef)
// 	return err
// }

func (h *HandlerManager) replySpongebobMeme(s *discordgo.Session, msgRef *discordgo.MessageReference, content string) error {
	memeBuffer, err := h.memeGenerator.GenerateMeme(content)
	if err != nil {
		_, _ = s.ChannelMessageSendReply(msgRef.ChannelID, replyError, msgRef)
		log.Println("Failed to create meme")
	}

	msgWithMeme := &discordgo.MessageSend{
		File: &discordgo.File{
			Name:        "meme-"+msgRef.MessageID+".jpg",
			ContentType: "image/jpeg",
			Reader:      memeBuffer,
		},
		Reference: msgRef,
	}
	_, err = s.ChannelMessageSendComplex(msgRef.ChannelID, msgWithMeme)

	return err
}

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, msg string, err error) error {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	return err
}
