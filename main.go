package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	dir          = "downloads"
	outputFormat = "mp4"
	originalExt  = "." + outputFormat
	convertedExt = ".webm"
)

// Run the bot client
func main() {
	if TOKEN == "" {
		log.Fatal("Secret token (in token.go) is missing.\nThe token can be found at https://discord.com/developers/applications/")
	}
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Printf("Error creating base download directory: %v", err)
	}

	dg, err := discordgo.New("Bot " + TOKEN)
	if err != nil {
		log.Fatal("Error creating session: ", err)
	}
	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening Discord connection, ", err)
	}
	log.Print("dvembed bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	_ = dg.Close()

	//app := &discordgo.Application{}
	//app.Name = "dvembed"
	//app.Description = "Properly embeds media from v.redd.it"
	//app, err = dg.ApplicationCreate(app)
	//log.Printf("ApplicationCreate: err: %+v, app: %+v", err, app)

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Printf("%v `%v` %v %v", m.Author.Username, m.Content, m.Embeds, &m.Embeds)
	if m.Author.ID == s.State.User.ID {
		return
		//c, err := s.State.Channel(m.ChannelID)
		//if err != nil {
		//	log.Printf("Could not find channel: %v", err)
		//	return
		//}
		//g, err := s.State.Guild(c.GuildID)
		//if err != nil {
		//	log.Printf("Could not find the guild for channel: %v, %v", c, err)
		//	return
		//}
		//log.Print(g)
	}
	switch {
	case strings.Contains(m.Content, "v.redd.it"):
		handleVredditLink(s, m)
	case strings.Contains(m.Content, "reddit.com"):
		handleRedditLink(s, m)
	case strings.Contains(m.Content, "test imgur link"):
		handleImgurTest(s, m)
	default:
		return
	}

}
