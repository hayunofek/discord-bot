package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/hayunofek/discord-bot/cmd"
	"github.com/hayunofek/discord-bot/music"

	"github.com/bwmarrin/discordgo"
)

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

// Discord Commands
var discordCommands []cmd.DiscordCommand = []cmd.DiscordCommand{
	{
		Name:     "play",
		Function: music.PlayCommand,
	},
}

var s *discordgo.Session

func init() {
	tokenFile, err := ioutil.ReadFile("token")
	if err != nil {
		log.Fatalln("Couldn't read token file!")
	}

	token := strings.TrimSpace(strings.ReplaceAll(string(tokenFile), "\n", ""))
	if token == "" {
		log.Fatalln("token file is empty!")
	}

	BotToken = &token
}

func main() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up and running!")
		s.UpdateGameStatus(0, "!coolstuff")
	})

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	defer s.Close()

	defer func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-stop
		log.Println("Gracefully shutdowning")
	}()

	s.AddHandler(func(s *discordgo.Session, i *discordgo.MessageCreate) {
		if strings.HasPrefix(i.Content, cmd.PREFIX_SIGN) {
			fmt.Printf("Got command in channel %s", i.ChannelID)
			s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("Author: %s, Message: %s", i.Author.Username, i.Content))

			for _, command := range discordCommands {
				if strings.HasPrefix(i.Content, fmt.Sprintf("!%s ", command.Name)) {
					resp, err := command.Function(s, i, &command)

					// Don't display error to user
					if err != nil {
						log.Printf("Got an error: %v", err)
						continue
					}

					s.ChannelMessageSend(i.ChannelID, resp)
				}
			}
		}
	})

}
