package main

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// readMessage looks through en entire message and sends any links to handleLink
func readMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	//log.Printf("Message by %s contains a reddit link: %s", m.Author.Username, m.Content)
	log.Print("Message sent on channel: ", m.ChannelID)
	msg := strings.Split(m.Content, " ")
	for _, unit := range msg {

		// Protection against URL copied without trailing `/`
		if !strings.HasSuffix(unit, "/") {
			unit += "/"
		}

		Url, err := url.Parse(unit)
		if err != nil {
			log.Printf("Message unit: `%s`, did not parse as URl: %v", unit, err)
			continue
		}

		err = handleLink(Url, unit, s, m)
		if err != nil {
			log.Printf("%v", err)
		}
	}
}

// handleLink processes valid links, and sends a message with the downloaded
// media, to the channel the link was posted. It can determine whether a long
// form link such as
// https://www.reddit.com/r/pics/comments/lpivdq/a_sense_of_priorities/, contains
// v.redd.it media. This is checked by looking for "fallback_url" within the
// post's JSON data view. Currently this is found by reading the entire JSON
// response, then searching for the term.
func handleLink(Url *url.URL, unit string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	switch {
	case strings.Contains(unit, "v.redd.it"):
		log.Printf("Message unit: `%s`, is valid URl", unit)

	case strings.Contains(unit, "reddit.com"):
		log.Printf("Message unit: `%s`, is valid URl", unit)
		log.Print("\tChecking JSON for media type...")

		err := verifyLongLink(Url)
		if err != nil {
			return fmt.Errorf("%v. %v", err, "Cannot confirm that content type is video. End process")
		}
		log.Printf("\t...fallback_url found, v.redd.it media confirmed. Continue...")
	}

	log.Print("Entering download process >")

	f, err := download(Url)
	if err != nil {
		return fmt.Errorf("did not download: %v", err)
	}
	log.Printf("\t...download successful")

	if f.Size() > 8000000 {
		// Send to compress.go
		return fmt.Errorf("%v bytes is too large for Discord upload", f.Size())
	}

	err = removeEmbed(s, m)
	if err != nil {
		log.Printf("Error removing %s's embedded v.redd.it image: %v", m.Author.Username, err)
	}

	o, err := os.Open(f.Path)
	defer o.Close()
	if err != nil {
		log.Print("Error reading ", f.Name())
	}

	msg, err := s.ChannelFileSend(m.ChannelID, f.Name(), o)
	if err != nil {
		return fmt.Errorf("failed sending %s's media %s, %v", m.Author.Username, f.Name(), err)
	}

	log.Printf("msg: %v sent on channel: %v in server: %v", msg.ID, msg.ChannelID, msg.GuildID)

	return nil
}

func verifyLongLink(Url *url.URL) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", Url.String()+".json", nil)
	if err != nil {
		return fmt.Errorf("\t...Failed creating request for %v: %v", Url.String()+".json", err)
	}

	req.Header.Set("User-Agent", "Dvembed/0.0")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("\t...Failed `GET` request for %v, %v", Url.String()+".json", err)
	}

	j, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("\t...Error reading response, %v", err)
	}

	if !bytes.Contains(j, []byte(`"fallback_url"`)) {
		return fmt.Errorf("\t...fallback_url not found.")
	}

	return nil
}

func removeEmbed(s *discordgo.Session, m *discordgo.MessageCreate) error {
	blank := discordgo.MessageEmbed{}
	_, err := s.ChannelMessageEditEmbed(m.ChannelID, m.ID, &blank)
	return err
}
