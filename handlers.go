package main

import (
	"bytes"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// readMessage can determine whether a long form link such as
//	https://www.reddit.com/r/pics/comments/lpivdq/a_sense_of_priorities/, contains
//	v.redd.it media. This is checked by looking for "fallback_url" within the
//	post's JSON data view. Currently this is found by reading the entire JSON
//	response, then searching for the term.
func readMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Printf("Message by %s contains a reddit link: %s", m.Author.Username, m.Content)
	log.Print("Message sent on channel: ", m.ChannelID)
	msg := strings.Split(m.Content, " ")
	for _, unit := range msg {
		Url, err := url.Parse(unit)
		if err != nil {
			log.Printf("Message unit: `%s`, did not parse as URl: %v", unit, err)
			continue
		}

		switch {
		case strings.Contains(unit, "v.redd.it"):
			log.Printf("Message unit: `%s`, is valid URl", unit)
			processLink(Url, s, m)

		case strings.Contains(unit, "reddit.com"):
			log.Printf("Message unit: `%s`, is valid URl", unit)
			log.Print("\tChecking JSON for media type...")

			client := &http.Client{}
			req, err := http.NewRequest("GET", Url.String()+".json", nil)
			if err != nil {
				log.Printf("\t...Failed creating request for %v", Url.String()+".json")
				log.Printf("\t...Cannot confirm that content type is video. End process")
				continue
			}

			req.Header.Set("User-Agent", "Dvembed/0.0")
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("\t...Failed `GET` request for %v, %v", Url.String()+".json", err)
				log.Printf("\t...Cannot confirm that content type is video. End process")
				continue
			}

			j, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("\t...Error reading response, %v", err)
				log.Printf("\t...Cannot confirm that content type is video. End process")
				continue
			}

			if !bytes.Contains(j, []byte(`"fallback_url"`)) {
				log.Printf("\t...No fallback_url found, file is not v.redd.it media. End Process")
				continue
			}

			log.Printf("\t...fallback_url found, v.redd.it media confirmed. Continue...")
			processLink(Url, s, m)
		}
	}

}

func processLink(Url *url.URL, s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Print("Entering download process >")

	f, err := download(Url)

	if err != nil {
		log.Print("Did not download: ", err)
		return
	}

	if f.Size() > 8000000 {
		// Send to format.go
		log.Printf("%v bytes is too large for Discord upload", f.Size())
		return
	}

	blank := discordgo.MessageEmbed{}
	_, err = s.ChannelMessageEditEmbed(m.ChannelID, m.ID, &blank)
	if err != nil {
		log.Printf("Error removing %s's embedded v.redd.it image: %v", m.Author.Username, err)
	}
	o, err := os.Open(f.Path)
	if err != nil {
		log.Print("Error reading ", f.Name())
	}
	msg, err := s.ChannelFileSend(m.ChannelID, f.Name(), o)
	if err != nil {
		log.Printf("Error sending %s's media %s, %v", m.Author.Username, f.Name(), err)
	}
	_ = o.Close()
	log.Printf("msg: %v sent on channel: %v in server: %v", msg.ID, msg.ChannelID, msg.GuildID)
}

func handleImgurTest(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, _ = s.ChannelMessageSend(m.ChannelID, "https://i.imgur.com/ZWexH7h.mp4")
}
