package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"net/url"
	"os"
	"strings"
)

func handleVredditLink(s *discordgo.Session, m *discordgo.MessageCreate) {
	c := m.ChannelID
	log.Printf("Message by %s  contains v.redd.it link: %s", m.Author.Username, m.Content)
	log.Print("Message sent on channel: ", c)
	uu := strings.Split(m.Content, " ")
	for _, u := range uu {
		if strings.Contains(u, "v.redd.it") || strings.Contains(u, "reddit.com") {
			u, err := url.Parse(u)
			if err != nil {
				log.Printf("Message segment: `%s`, did not parse as URl: %v", u, err)
				continue
			}
			log.Printf("Message segment: `%s`, is valid URl", u)
			log.Print("Entering download process >")
			f, err := download(u.String())
			if err != nil {
				log.Print("Did not download: ", err)
				continue
			}
			if f.Size() > 8000000 {
				log.Printf("%v bytes is too large for Discord upload", f.Size())
				continue
			}

			blank := discordgo.MessageEmbed{}
			_, err = s.ChannelMessageEditEmbed(c, m.ID, &blank)
			if err != nil {
				log.Printf("Error removing %s's embedded v.redd.it image: %v", m.Author.Username, err)
			}
			o, err := os.Open(f.Path)
			if err != nil {
				log.Print("Error reading ", f.Name())
			}
			msg, err := s.ChannelFileSend(c, f.Name(), o)
			if err != nil {
				log.Printf("Error sending %s's media %s, %v", m.Author.Username, f.Name(), err)
			}
			_ = o.Close()
			log.Printf("msg: %v sent on channel: %v in server: %v", msg.ID, msg.ChannelID, msg.GuildID)
		}
	}
}

func handleRedditLink(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Printf("Message by %s contains reddit link: %s", m.Author.Username, m.Content)
	//// TODO check if standard reddit link contains v.redd.it media through reddit api
	//uu := strings.Split(m.Content, " ")
	//for _, u := range uu {
	//	if strings.Contains(u, "reddit.com") {
	//		u, err := url.Parse(u)
	//		if err != nil {
	//			log.Printf("Message segment: `%s`, did not parse as URl: %v", u, err)
	//			continue
	//		}
	//		log.Printf("Message segment: `%s`, is valid URl", u)
	//		// getVredditLink(URL)
	//	}
	//}
	handleVredditLink(s, m)
}

func handleImgurTest(s *discordgo.Session, m *discordgo.MessageCreate) {
	// TODO imgur links seem to embed nicely, so see about using imgur api
	_, _ = s.ChannelMessageSend(m.ChannelID, "https://i.imgur.com/ZWexH7h.mp4")
}
