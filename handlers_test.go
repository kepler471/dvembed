package main

import (
	"bytes"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
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

func TestJsonGetRequest(t *testing.T) {
	client := &http.Client{}
	Url := "https://www.reddit.com/r/AnimalsBeingBros/comments/ip89wl/possibly_the_most_patient_kitty_in_the_world_with/"
	req, err := http.NewRequest("GET", Url+".json", nil)
	if err != nil {
		log.Printf("\tCould not GET with NewRequest")
	}

	req.Header.Set("User-Agent", "Dvembed/0.0")
	//resp, err := http.Get(Url+".json")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("\tFailed to reach JSON at %v, %v", Url+".json", err)
		t.Fail()
	}

	type metadata struct {
		Kind string `json:"kind"`
	}

	var v []metadata

	j, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("\tJson read? %v", string(j))
	}
	err = json.Unmarshal(j, &v)
	if err != nil {
		log.Printf("\tJson Unmarshal? %v", v)
	}
	//err = json.NewDecoder(resp.Body).Decode(&v)
	if err != nil {
		log.Printf("\t\tFailed to decode JSON from %v, %v", Url+".json", err)
		t.Fail()
	}
	log.Printf("Show interface 'v': %v", v)
	if bytes.Contains(j, []byte("fallback_url")) {
		log.Printf("Found 'fallback_url'!!")
	}
}
