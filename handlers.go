package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Printf("%v `%v`", m.Author.Username, m.Content)
	if m.Author.ID == s.State.User.ID && !strings.Contains(m.Content, "BOT-READABLE") {
		return
	}
	if strings.Contains(m.Content, "v.redd.it") || strings.Contains(m.Content, "reddit.com") {
		readMessage(s, m)
	}
}

// readMessage looks through en entire message and sends any links to handleLink
func readMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Print("Message sent on channel: ", m.ChannelID)
	msg := strings.Split(m.Content, " ")
	for _, unit := range msg {

		// Protection against URL copied without trailing `/`
		if !strings.HasSuffix(unit, "/") {
			unit += "/"
		}

		// TODO: Move this to handleLink?
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
	// TODO: use regex match from youtube-dl/extractor/reddit.py
	switch {
	case strings.Contains(unit, "v.redd.it"):
		log.Printf("Message unit: `%s`, is valid URl", unit)

	case strings.Contains(unit, "reddit.com"):
		log.Printf("Message unit: `%s`, is valid URL", unit)
		log.Print("\tChecking JSON for media type...")

		dashUrl := findDashUrl(Url)
		if dashUrl == "" {
			return fmt.Errorf("%v. %v", "No dash_url found.", "Cannot confirm that content type is video. End process")
		}

		log.Printf("\t...fallback_url found, %v v.redd.it media confirmed. Continue...", dashUrl)
		Url, _ = url.Parse(dashUrl)
		Url.Path += "/"
	}

	log.Println("Url = ", Url)
	log.Println("Entering download process >")

	id := path.Base(Url.String())
	log.Println("id = ", id)
	f := media{
		Id: id,
	}
	// Create proper tmp folders
	tmp := filepath.Join(dir, f.Id)
	log.Println("tmp = ", tmp)
	err := os.MkdirAll(tmp, 0755)
	if err != nil {
		log.Printf("\tError creating sub-directory: %v", err)
	}

	f.Path = filepath.Join(tmp, id+originalExt)
	log.Println("f.Path = ", f.Path)
	// Download all dash urls
	log.Println("mpd link = ", Url.String()+"DASHPlaylist.mpd")
	mpd, _, _ := download(Url.String()+"DASHPlaylist.mpd", tmp)
	mpd, _ = os.Open(mpd.Name())
	defer mpd.Close()
	variantLinks := decode(mpd)
	fmt.Println(variantLinks)

	// Download all variants contained in the DASHPlaylist.mpd
	//var variantFiles []*os.File
	sizes := make(map[string]int64)
	for _, vl := range variantLinks {
		_, size, _ := download(Url.String()+vl, tmp)
		sizes[vl] = size
	}

	// Find best quality within limit
	var audio string
	var best string
	var bestSize int64
	for variant, size := range sizes {
		if !strings.Contains(variant, "AUDIO") && !strings.Contains(variant, "audio") {
			if size > bestSize && size < messageSizeLimit {
				// TODO: Might want to check for bestSize + audio size
				best = variant
			}
		} else {
			// TODO: make some choice of audio quality
			audio = variant
		}
	}

	// mux
	if err = mux(best, audio, id+originalExt, tmp); err != nil {
		log.Printf("\tmux failed, %v. Try with best video, no audio", err)
		f.FileInfo, err = os.Stat(filepath.Join(tmp, best))
	} else {
		f.FileInfo, err = os.Stat(f.Path)
	}

	if err != nil {
		return fmt.Errorf("\terror finding downloaded file: %v", err)
	}

	o, err := os.Open(f.Path)
	if err != nil {
		return fmt.Errorf("error reading %v, %v", f.Name(), err)
	}
	defer o.Close()

	msg, err := s.ChannelFileSend(m.ChannelID, f.Name(), o)
	if err != nil {
		return fmt.Errorf("failed sending %s's media %s, %v", m.Author.Username, f.Name(), err)
	}

	log.Printf("msg: %v sent on channel: %v in server: %v", msg.ID, msg.ChannelID, msg.GuildID)

	return nil
}

// findDashUrl searches for a DASH url and returns the first found
// This function should have multiple ways to find the DASH url. There are
// other json entries that contain the DSAH mpd link, but also the html
// contains the dash link under data-mpd-url
func findDashUrl(Url *url.URL) string {
	resp, _ := fetchJson(Url)
	j, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("\t...Error reading JSON response, %v", err)
		return ""
	}

	if !bytes.Contains(j, []byte(`"fallback_url"`)) {
		// return j, fmt.Errorf("\t...fallback_url not found.")
		fmt.Errorf("\t...fallback_url not found in JSON response.")
		return ""
	}

	var result interface{}
	json.Unmarshal(j, &result)
	var dashUrl string
	var walk func(value reflect.Value)
	walk = func(v reflect.Value) {
		if v.String() == "dash_url" {
			fmt.Printf("Visiting %v\n", v)
		}
		for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.Array, reflect.Slice:
			for i := 0; i < v.Len(); i++ {
				walk(v.Index(i))
			}
		case reflect.Map:
			for _, k := range v.MapKeys() {
				if k.String() == "url_overridden_by_dest" {
					dashUrl = fmt.Sprintf("%v", v.MapIndex(k))
					log.Printf("\tFound url_overridden_by_dest: %q\n", dashUrl)
				}
				if k.String() != "secure_media" {
					walk(v.MapIndex(k))
				}
			}
		}
	}

	walk(reflect.ValueOf(result))
	return dashUrl
}

// reddit.py regex
func isVRedditLink() {}

// reddit.py regex
func isRedditLink() {}
