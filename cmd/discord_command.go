package cmd

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const PREFIX_SIGN string = "!"

func (dc *DiscordCommand) GetMyCommandPrefix() string {
	return fmt.Sprintf("%s%s", PREFIX_SIGN, dc.Name)
}

type DiscordCommand struct {
	Name     string
	Function func(*discordgo.Session, *discordgo.MessageCreate, *DiscordCommand) (string, error)
}
