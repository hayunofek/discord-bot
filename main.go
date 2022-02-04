package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var s *discordgo.Session

func init() {
	flag.Parse()

	if *BotToken == "" {
		log.Fatalln("Token is missing!")
	}

	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	defer func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-stop
		log.Println("Gracefully shutdowning")
	}()

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	defer s.Close()
}

func main() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.PresenceUpdate) {
		s.ChannelMessageSend("", fmt.Sprintf("Hey %s, your state is: %s", i.User.Username, i.Status))
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.MessageCreate) {
		if strings.HasPrefix(i.Content, "!") {
			fmt.Printf("Got command in channel %s", i.ChannelID)
			s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("Author: %s, Message: %s", i.Author.Username, i.Content))
		}
	})

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})

}
