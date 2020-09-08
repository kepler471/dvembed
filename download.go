package main

import (
	"log"
	"net/url"
	"os"
	"os/exec"
)

const dir = "downloads"

// Download handles retrieving the v.redd.it media
func Download(URL string) error {
	u, err := url.Parse(URL)
	//name := path.Base(URL)
	if err != nil {
		log.Fatalf("Could not parse %v: %v\n", u.Path, err)
	}
	err = os.Mkdir(dir, 0755)
	if err != nil {
		log.Printf("Directory not created: %v\n", err)
	}
	cmd := exec.Command(
		"youtube-dl",
		"--id",
		"--write-info-json",
		"--merge-output-format",
		"mp4",
		URL,
	)
	cmd.Dir = dir
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Could not run %v\n", cmd.Args)
	}
	return err
}
