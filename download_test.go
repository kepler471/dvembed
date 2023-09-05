package main

import (
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestDownloadDASHPlaylistFromVRedditLink(t *testing.T) {
	URLs := []string{
		"https://v.redd.it/duir5tuwswl51",
		"https://v.redd.it/5ltubsoyawl51",
		"https://v.redd.it/dttgnvp69wl51",
		"https://v.redd.it/e497qwjsh1m51",
		"https://v.redd.it/b8b0ha50l5g41",
		"https://v.redd.it/zv89llsvexdz",
		"https://v.redd.it/5wiorpupd3mb1",
		"https://v.redd.it/vxtosbjn47mb1",
	}
	for _, URL := range URLs {
		Url, _ := url.Parse(URL)
		id := path.Base(Url.String())
		tmp := filepath.Join(dir, id)
		err := os.MkdirAll(tmp, 0755)
		if err != nil {
			t.Errorf("\tError creating sub-directory: %v", err)
		}
		Url = Url.JoinPath("DASHPlaylist.mpd")
		_, _, err = download(Url.String(), tmp)
		t.Logf(`MPD file to download: %q`, Url)
		if err != nil {
			t.Errorf(`%q failed`, URL)
		}
	}
}

func TestDownloadDASHPlaylistFromRedditLink(t *testing.T) {
	URLs := []string{
		"https://www.reddit.com/r/IdiotsInCars/comments/ioqqbf/i_know_ill_cut_in_front_of_this_semi/",
		"https://www.reddit.com/r/AnimalsBeingBros/comments/ip89wl/possibly_the_most_patient_kitty_in_the_world_with/",
		"https://old.reddit.com/r/StarWars/comments/l1l7f6/finished_this_last_night_took_me_30_hours_to/",
		"https://www.reddit.com/r/videos/comments/6rrwyj/that_small_heart_attack/",
		"https://www.reddit.com/r/videos/comments/6rrwyj",
		"https://www.reddit.com/r/MadeMeSmile/comments/6t7wi5/wait_for_it/",
		"https://old.reddit.com/r/MadeMeSmile/comments/6t7wi5/wait_for_it/",
		"https://www.reddit.com/r/videos/comments/6t7sg9/comedians_hilarious_joke_about_the_guam_flag/",
		"https://www.reddit.com/r/videos/comments/6t75wq/southern_man_tries_to_speak_without_an_accent/",
		"https://nm.reddit.com/r/Cricket/comments/8idvby/lousy_cameraman_finds_himself_in_cairns_line_of/",
		"https://www.reddit.com/r/IAmTheMainCharacter/comments/169ur5y/yoko/",
	}
	for _, URL := range URLs {
		Url, _ := url.Parse(URL)
		id := path.Base(Url.String())
		_ = filepath.Join(dir, id)
		tmp := filepath.Join(dir, id)
		err := os.MkdirAll(tmp, 0755)
		if err != nil {
			t.Errorf("\tError creating sub-directory: %v", err)
		}
		_, _, err = download(Url.String(), tmp)
		if err != nil {
			t.Errorf(`%q failed, %q`, URL, err)
		}
	}
}

// Downloads JSON files
func Test_fetchJson(t *testing.T) {
	URLs := []string{
		"https://www.reddit.com/r/IdiotsInCars/comments/ioqqbf/i_know_ill_cut_in_front_of_this_semi/",
		"https://www.reddit.com/r/AnimalsBeingBros/comments/ip89wl/possibly_the_most_patient_kitty_in_the_world_with/",
		"https://old.reddit.com/r/StarWars/comments/l1l7f6/finished_this_last_night_took_me_30_hours_to/",
		"https://www.reddit.com/r/videos/comments/6rrwyj/that_small_heart_attack/",
		"https://www.reddit.com/r/videos/comments/6rrwyj",
		"https://www.reddit.com/r/MadeMeSmile/comments/6t7wi5/wait_for_it/",
		"https://old.reddit.com/r/MadeMeSmile/comments/6t7wi5/wait_for_it/",
		"https://www.reddit.com/r/videos/comments/6t7sg9/comedians_hilarious_joke_about_the_guam_flag/",
		"https://www.reddit.com/r/videos/comments/6t75wq/southern_man_tries_to_speak_without_an_accent/",
		"https://nm.reddit.com/r/Cricket/comments/8idvby/lousy_cameraman_finds_himself_in_cairns_line_of/",
		"https://www.reddit.com/r/IAmTheMainCharacter/comments/169ur5y/yoko/",
	}
	for _, URL := range URLs {
		Url, _ := url.Parse(URL)
		id := path.Base(Url.String())
		tmp := filepath.Join(dir, id)
		resp, err := fetchJson(Url)
		if err != nil {
			t.Errorf(`fetchJson %q failed, %q`, URL, err)
		}
		f, err := os.Create(filepath.Join(tmp, id+".json"))
		if err != nil {
			t.Errorf("\tcould not create file, %v\n", err)
		}
		_, err = io.Copy(f, resp.Body)
		if err != nil {
			t.Errorf(`could not write file`)
		}
	}
}
