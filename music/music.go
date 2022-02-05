package music

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/hayunofek/discord-bot/cmd"
	"github.com/kkdai/youtube/v2"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

type YoutubeSong struct {
	URL string
}

// A function which does the play music command. It downloads the song from youtube to an mp4 file,
// converts the mp4 file to a dca (Discord Audio File) and then plays it
func PlayCommand(s *discordgo.Session, i *discordgo.MessageCreate, dc *cmd.DiscordCommand) (string, error) {
	ys := YoutubeSong{
		URL: strings.TrimPrefix(i.Content, dc.GetMyCommandPrefix()+" "),
	}

	log.Printf("\nDownloading song: %s\n", ys.URL)

	downloadedVideoFileName, err := ys.download()
	if err != nil {
		log.Printf("Unable to download video, url: %s, error: %v", ys.URL, err)
		return "", err
	}

	defer os.Remove(downloadedVideoFileName)

	vc, err := joinVoiceChannel(s, i)
	if err != nil {
		log.Printf("\nGot an error trying to join voice channel, error: %v", err)
		return "", err
	}

	defer vc.Disconnect()

	// Sleep for a little while before playing the sound
	time.Sleep(250 * time.Millisecond)

	vc.Speaking(true)
	defer vc.Speaking(false)

	dgvoice.PlayAudioFile(vc, downloadedVideoFileName, make(chan bool))

	// Sleep for a little while before exiting
	time.Sleep(250 * time.Millisecond)

	log.Printf("\nFinishing...\n")
	return fmt.Sprintf("You chose to play music my friend. Your song name: %s", ys.URL), nil
}

func joinVoiceChannel(s *discordgo.Session, i *discordgo.MessageCreate) (*discordgo.VoiceConnection, error) {
	// Find the channel that the message came from
	channel, err := s.State.Channel(i.ChannelID)
	if err != nil {
		log.Printf("Error getting channel, err: %v", err)
	}

	// Find the guild for that channel
	guild, err := s.State.Guild(channel.GuildID)
	if err != nil {
		log.Printf("Error getting guild, err: %v", err)
	}

	guildID := ""
	channelID := ""

	for _, vs := range guild.VoiceStates {
		if vs.UserID == i.Author.ID {
			guildID = guild.ID
			channelID = vs.ChannelID
			break
		}
	}

	if guildID == "" || channelID == "" {
		err = errors.New("couldn't find channel id and guild id for the requesting user")
		log.Printf("\n%v\n", err)
		return nil, err
	}

	log.Printf("\nJoining Voice Channel %s\n", i.ChannelID)

	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, false)
	if err != nil {
		log.Printf("Error joining voice channel, err: %v", err)
		return nil, err
	}

	return vc, nil
}

// This function downloads the youtube song, returns the name of the filename it downloaded
// and an error if it failed or nil
func (ys *YoutubeSong) download() (string, error) {
	parsedURL, err := url.Parse(ys.URL)
	if err != nil {
		log.Printf("Unable to decode url: %s, err: %v", ys.URL, err)
		return "", err
	}

	videoID := parsedURL.Query().Get("v")
	if videoID == "" {
		log.Printf("Unable to get video id from query: %s", parsedURL.RawQuery)
		return "", err
	}

	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		log.Printf("Unable to get video from youtube, video id: %s, err: %v", videoID, err)
		return "", err
	}

	formats := video.Formats.WithAudioChannels() // only get videos with audio
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		log.Printf("Unable to get stream of video from youtube, video id: %s, err: %v", videoID, err)
		return "", err
	}

	fileName := fmt.Sprintf("%s.mp4", videoID)

	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Unable to create file, video id: %s, err: %v", videoID, err)
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		log.Printf("Unable to copy stream to file, video id: %s, err: %v", videoID, err)
		return "", err
	}

	return fileName, nil
}

// Unused
func convertMP4ToOpus(filename string) (string, error) {
	opusFilename := "DB" + strings.Split(filename, ".")[0] + ".opus"
	ffmpegArgs := fmt.Sprintf("-y -i %s -strict -2 %s", filename, opusFilename)
	ffmpegArgsSplitted := strings.Split(ffmpegArgs, " ")
	cmd := exec.Command(
		"ffmpeg",
		ffmpegArgsSplitted...,
	)

	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Got an error trying to convert MP4 to DCA, error: %v", err)
		return "", err
	}

	return opusFilename, nil
}

// Unused
func fillBufferFromOpus(file *os.File) ([][]byte, error) {
	buffer := [][]byte{}
	// var opuslen int16
	FRAME_SIZE := 960
	CHANNELS := 2
	var err error

	for {
		// // Read opus frame length from dca file.
		// err := binary.Read(file, binary.LittleEndian, &opuslen)

		// // If this is the end of the file, just return.
		// if err == io.EOF || err == io.ErrUnexpectedEOF {
		// 	err := file.Close()
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	return buffer, nil
		// }

		// if err != nil {
		// 	fmt.Println("Error reading from dca file :", err)
		// 	return nil, err
		// }
		// fmt.Printf("\nlen: %d\n", opuslen)

		// if opuslen < 0 {
		// 	return buffer, nil
		// }

		// Read encoded pcm from dca file.
		InBuf := make([]byte, FRAME_SIZE*CHANNELS)
		err = binary.Read(file, binary.LittleEndian, &InBuf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			log.Printf("\nClosing!\n")
			err := file.Close()
			if err != nil {
				return nil, err
			}
			return buffer, nil
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}
}
