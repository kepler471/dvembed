package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

type BaseUrl struct {
	Data string `xml:",chardata"`
}

type DashVideo int

const (
	K_600 DashVideo = iota
	M_1_2
	M_2_4
	M_4_8
	R_240
	R_360
	R_480
	R_720
)

// decode returns all variations of BaseUrl in a DASHPlaylist.mpd file.
func decode(mpd io.Reader) []string {
	dec := xml.NewDecoder(mpd)
	var stack []string
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "xmlselect: %v\n", err)
			os.Exit(1)
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			if tok.Name.Local == "BaseURL" {
				var Url BaseUrl
				_ = dec.DecodeElement(&Url, &tok)
				stack = append(stack, Url.Data)
			}
		}
	}
	return stack
}
