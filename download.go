package main

import (
	"log"
	"os"
	"os/exec"
	"path"
)

type media struct {
	os.FileInfo
	Id   string
	Path string
}

// download handles retrieving the v.redd.it media. Initialises
// a media instance for each URL.
func download(URL string) (*media, error) {
	f := media{
		Id: path.Base(URL),
	}
	cmd := exec.Command(
		"youtube-dl",
		"-v",
		//"--id",              // use id as name
		"--output",
		f.Id,
		"--write-info-json", // save file information
		//"--restrict-filenames",
		"--merge-output-format", // Downloading the best available audio and video
		outputFormat,
		URL,
	)
	cmd.Dir = path.Join(dir, f.Id)
	err := os.MkdirAll(cmd.Dir, 0755)
	if err != nil {
		log.Printf("Error creating sub-directory: %v", err)
		cmd.Dir = dir
	}
	log.Print("...run youtube-dl...")
	err = cmd.Run()
	if err != nil {
		log.Printf("Failed process: %v", cmd.Args)
		return &f, err
	}
	f.Path = path.Join(cmd.Dir, f.Id+originalExt)
	f.FileInfo, err = os.Stat(f.Path)
	if err != nil {
		log.Print("Error finding downloaded file: ", err)
	}
	// TODO want to return output error from youtube-dl
	return &f, err
}
