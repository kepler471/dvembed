package main

import (
	"log"
	"net/url"
	"os/exec"
)

// Download handles retrieving the v.redd.it media
func Download(URL string) error {
	u, err := url.Parse(URL)
	if err != nil {
		log.Fatalf("Could not parse %v: %v", u.Path, err)
	}
	cmd := exec.Command("youtube-dl", "--merge-output-format", "mp4", URL)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Could not run %v", cmd.Args)
	}
	return err
}
