package music

import (
	"fmt"
	"strings"

	"github.com/hayunofek/discord-bot/cmd"

	"github.com/bwmarrin/discordgo"
)

type Song struct {
	Name string
}

func PlayCommand(s *discordgo.Session, i *discordgo.MessageCreate, dc *cmd.DiscordCommand) (string, error) {
	songName := strings.Split(strings.TrimPrefix(dc.GetMyCommandPrefix(), i.Content), " ")
	return fmt.Sprintf("You chose to play music my friend. Your song name: %s", songName), nil
}
