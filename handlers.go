package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
)

func handleVredditLink(s *discordgo.Session, m *discordgo.MessageCreate) {
	c := m.ChannelID
	log.Printf("Message by %s  contains v.redd.it link: %s\n", m.Author.Username, m.Content)
	log.Print("Message sent on channel: ", c)
	uu := strings.Split(m.Content, " ")
	for _, u := range uu {
		if strings.Contains(u, "v.redd.it") {
			u, err := url.Parse(u)
			if err != nil {
				log.Printf("Message segment: `%s`, did not parse as URl: %v\n", u, err)
				continue
			}
			log.Printf("Message segment: `%s`, is valid URl\n", u)
			f, err := download(u.String())
			if err != nil {
				log.Println("Error downloading f")
				continue
			}
			if f.Size() > 8000000 {
				log.Printf("%v bytes is too large for Discord upload", f.Size())
				continue
			}
			blank := &discordgo.MessageEmbed{}
			_, err = s.ChannelMessageEditEmbed(c, m.ID, blank)
			if err != nil {
				log.Printf("Error removing %s's embedded v.redd.it image", m.Author.Username)
			}
			o, err := os.Open(path.Join(dir, f.Id, f.Name()))
			if err != nil {
				log.Print("Error reading ", f.Name())
			}
			_, err = s.ChannelFileSend(c, f.Name(), o)
			if err != nil {
				log.Printf("Error uploading %s's media %s, %v", m.Author.Username, f.Name(), err)
			}
			_ = o.Close()
		}
	}
}

func handleRedditLink(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Printf("Message by %s contains reddit link: %s\n", m.Author.Username, m.Content)
	// TODO check if standard reddit link contains v.redd.it media through reddit api
	uu := strings.Split(m.Content, " ")
	for _, u := range uu {
		if strings.Contains(u, "reddit.com") {
			u, err := url.Parse(u)
			if err != nil {
				log.Printf("Message segment: `%s`, did not parse as URl: %v\n", u, err)
				continue
			}
			log.Printf("Message segment: `%s`, is valid URl\n", u)
			// getVredditLink(URL)
			handleVredditLink(s, m)
		}
	}
}

func handleImgurTest(s *discordgo.Session, m *discordgo.MessageCreate) {
	// TODO imgur links seem to embed nicely, so see about using imgur api
	_, _ = s.ChannelMessageSend(m.ChannelID, "https://i.imgur.com/ZWexH7h.mp4")
}
