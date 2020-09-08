package main

import "testing"

func TestDownload(t *testing.T) {
	URL := "https://v.redd.it/duir5tuwswl51"
	err := Download(URL)
	if err != nil {
		t.Errorf(`youtube-dl %q failed`, URL)
	}
}
