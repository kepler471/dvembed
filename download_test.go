package main

import "testing"

func TestDownloadVRedditLink(t *testing.T) {
	URL := "https://v.redd.it/duir5tuwswl51"
	err := Download(URL)
	if err != nil {
		t.Errorf(`youtube-dl %q failed`, URL)
	}
	URL = "https://v.redd.it/5ltubsoyawl51"
	err = Download(URL)
	if err != nil {
		t.Errorf(`youtube-dl %q failed`, URL)
	}
}

func TestDownloadRedditLink(t *testing.T) {
	URL := "https://www.reddit.com/r/IdiotsInCars/comments/ioqqbf/i_know_ill_cut_in_front_of_this_semi/"
	err := Download(URL)
	if err != nil {
		t.Errorf(`youtube-dl %q failed`, URL)
	}
}
