package main

import (
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
)

type Media struct {
	Id   string
	Url  string
	Info os.FileInfo
}

// Download handles retrieving the v.redd.it media
func Download(URL string) (Media, error) {
	u, err := url.Parse(URL)
	if err != nil {
		log.Fatalf("Could not parse %v: %v\n", u.Path, err)
	}
	var m Media
	m.Url = URL
	m.Id = path.Base(URL)
	err = os.MkdirAll(Dir, 0755)
	if err != nil {
		log.Printf("Error creating directory: %v\n", err)
	}
	err = os.MkdirAll(path.Join(Dir, m.Id), 0755)
	if err != nil {
		log.Printf("Error creating sub-directory: %v\n", err)
	}
	cmd := exec.Command(
		"youtube-dl",
		"-v",
		"--id",              // use id as name
		"--write-info-json", // save file information
		"--restrict-filenames",
		"--merge-output-format", // Downloading the best available audio and video
		OutputFormat,
		URL,
	)
	cmd.Dir = path.Join(Dir, m.Id)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Failed: %v\n", cmd.Args)
	}
	m.Info, err = os.Stat(path.Join(cmd.Dir, m.Id+OriginalExt))
	if err != nil {
		log.Println("Error finding downloaded file")
	}
	// TODO want to return output error from youtube-dl
	return m, nil
}
