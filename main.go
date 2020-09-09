package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
)

const (
	Dir          = "downloads"
	OutputFormat = "mp4"
	OriginalExt  = "." + OutputFormat
	ConvertedExt = ".webm"
)

// Run the bot client
func main() {

	dg, err := discordgo.New("Bot " + TOKEN)
	if err != nil {
		log.Fatalln("Error creating session: ", err)
	}
	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
	err = dg.Open()
	if err != nil {
		log.Print("Error opening Discord connection, ", err)
		return
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
	//log.Printf("ApplicationCreate: err: %+v, app: %+v\n", err, app)

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Printf("%v `%v` %v %v", m.Author.Username, m.Content, m.Embeds, &m.Embeds)
	if m.Author.ID == s.State.User.ID {
		return
		//c, err := s.State.Channel(m.ChannelID)
		//if err != nil {
		//	log.Printf("Could not find channel: %v\n", err)
		//	return
		//}
		//g, err := s.State.Guild(c.GuildID)
		//if err != nil {
		//	log.Printf("Could not find the guild for channel: %v, %v\n", c, err)
		//	return
		//}
		//log.Println(g)
	}
	switch {
	case strings.Contains(m.Content, "v.redd.it"):
		handleVredditLink(s, m)
	case strings.Contains(m.Content, "reddit.com"):
		handleRedditLink(m)
	case strings.Contains(m.Content, "test imgur link"):
		handleImgurTest(s, m)
	default:
		return
	}

}

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
			file, err := Download(u.String())
			if err != nil {
				log.Println("Error downloading file")
				continue
			}
			if file.Info.Size() > 8000000 {
				log.Printf("%v bytes is too large for Discord upload", file.Info.Size())
				continue
			}
			blank := &discordgo.MessageEmbed{}
			_, err = s.ChannelMessageEditEmbed(c, m.ID, blank)
			if err != nil {
				log.Printf("Error removing %s's embedded v.redd.it image", m.Author.Username)
			}
			o, err := os.Open(path.Join(Dir, file.Id, file.Info.Name()))
			if err != nil {
				log.Print("Error reading ", file.Info.Name())
			}
			_, err = s.ChannelFileSend(c, file.Id, o)
			if err != nil {
				log.Printf("Error uploading %s's v.redd.it media, %s", m.Author.Username, file.Info.Name())
			}
			_ = o.Close()
		}
	}
}

func handleRedditLink(m *discordgo.MessageCreate) {
	log.Printf("Message by %s contains reddit link: %s\n", m.Author.Username, m.Content)
	// TODO check if standard reddit link contains v.redd.it media through reddit api
}

func handleImgurTest(s *discordgo.Session, m *discordgo.MessageCreate) {
	// TODO imgur links seem to embed nicely, so see about using imgur api
	_, _ = s.ChannelMessageSend(m.ChannelID, "https://i.imgur.com/ZWexH7h.mp4")
}
