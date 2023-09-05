package main

import (
	"bytes"
	"io"
	"net/url"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func Test_readMessage(t *testing.T) {
	type args struct {
		s *discordgo.Session
		m *discordgo.MessageCreate
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{name: "cat", args: args{
			s: nil,
			m: &discordgo.MessageCreate{Message: &discordgo.Message{
				Content: "https://www.reddit.com/r/AnimalsBeingBros/comments/ip89wl/possibly_the_most_patient_kitty_in_the_world_with/",
			}},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_findDashUrl(t *testing.T) {
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
		resp, err := fetchJson(Url)
		if err != nil {
			t.Errorf(`fetchJson %q failed, %q`, URL, err)
		}
		j, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("\tcould not read response, %v\n", err)
		}
		if !bytes.Contains(j, []byte(`"fallback_url`)) {
			t.Errorf(`DASH not found: %q`, Url)
		}
	}
}
