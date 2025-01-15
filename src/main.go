package main

import (
	"time"
)

type Input struct {
	input   []string
	command string
	param   string
}

type Config struct {
	Input    Input
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Timeout  time.Duration
	Results  []struct {
		Name string `json:"name"`
	} `json:"results"`
	history  cache
	Pokemons map[string]Pokemon
}

const (
	baseurl = "https://pokeapi.co/api/v2/"
)

func main() {

	cache := NewCache(5 * time.Second)
	cache.Add("Tim", make([]byte, 0))
	cache.Add("Paul", make([]byte, 0))

	config := Config{
		Next:     baseurl + "location-area/",
		Previous: "",
		Input:    Input{input: make([]string, 0), command: "", param: ""},
		Timeout:  5 * time.Second,
		history:  NewCache(5 * time.Second),
		Pokemons: make(map[string]Pokemon),
	}

	startRepl(&config)

}
