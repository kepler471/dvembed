package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

const (
	Dir          = "downloads"
	OutputFormat = "mp4"
	OriginalExt  = "." + OutputFormat
	ConvertedExt = ".webm"
)

// main will run the bot client
func main() {

	dg, err := discordgo.New("Bot " + TOKEN)
	if err != nil {
		log.Fatalln("Error creating session")
	}

	dg.AddHandler(messageCreate)

	app := &discordgo.Application{}
	app.Name = "dvembed"
	app.Description = "Properly embeds media from v.redd.it"
	app, err = dg.ApplicationCreate(app)
	log.Printf("ApplicationCreate: err: %+v, app: %+v\n", err, app)

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	switch {
	case strings.Contains(m.Content, "v.redd.it"):
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			log.Printf("Could not find channel: %v\n", err)
			return
		}
		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			log.Printf("Could not find the guild for channel: %v, %v\n", c, err)
			return
		}
		log.Println(g)
		return
	case strings.Contains(m.Content, "reddit.com"):
		// TODO check if standard reddit link contains v.redd.it media through reddit api
		return
	default:
		return
	}

}
