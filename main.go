package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const (
	dir              = "downloads"
	outputFormat     = "mp4"
	originalExt      = "." + outputFormat
	convertedExt     = ".webm"
	messageSizeLimit = 8_000_000
)

// TODO: use go embed for token

// Run the bot client
func main() {
	flag.StringVar(&TOKEN, "t", TOKEN, "Bot token")
	flag.Parse()
	if //goland:noinspection GoBoolExpressions
	TOKEN == "" {
		log.Fatal(
			"Secret token is missing. The token can be stored in token.go, or can be passed as\n " +
				"a command line argument when the bot is run.\n" +
				"Your token can be found at https://discord.com/developers/applications/")
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
	dg.Identify.Intents = discordgo.IntentsGuildMessages
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening Discord connection, %v", err)
	}

	log.Print("dvembed bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sc
	_ = dg.Close()

	//app := &discordgo.Application{}
	//app.Name = "dvembed"
	//app.Description = "Properly embeds media from v.redd.it"
	//app, err = dg.ApplicationCreate(app)
	//log.Printf("ApplicationCreate: err: %+v, app: %+v", err, app)
}
