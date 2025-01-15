package main

import (
	"fmt"
	"math/rand"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(config *Config) error
}

type Pokemon struct {
	Name           string
	BaseExperience int
	Height         int
	Weight         int
	Types          []Types
}

var commands map[string]cliCommand

func init() {

	commands = map[string]cliCommand{
		"pokedex": {
			name:        "pokedex",
			description: "Shows all catched Pokemons",
			callback:    commandPokedex,
		},
		"inspect": {
			name:        "inspect",
			description: "param1:  {name of pokemon } - Show Details of Pokemon, If it was Catched",
			callback:    commandInspect,
		},
		"catch": {
			name:        "catch",
			description: "param1: {name of pokemon } - Try to Catch the Pokemon",
			callback:    commandCatch,
		},

		"explore": {
			name:        "explore",
			description: "param1: {cityname} - Shows all the pokemon of the location in the Pokemon word",
			callback:    commandExplore,
		},
		"mapb": {
			name:        "mapb",
			description: "Shows previous page of locations in the Pokemon word",
			callback:    commandMapBack,
		},
		"map": {
			name:        "map",
			description: "Shows the next locations in the Pokemon word",
			callback:    commandMap,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    comandHelp,
		},
	}

}

func comandHelp(config *Config) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n")

	for _, value := range commands {
		fmt.Printf("%v : %v \n", value.name, value.description)
	}
	return nil
}

func commandExit(config *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(config *Config) error {
	fmt.Println("Pokedex > map")
	citiesResponse, err := GetData[CityResponse](config, config.Next)
	if err == nil {
		config.Next = citiesResponse.Next
		config.Previous = citiesResponse.Previous
		for _, r := range citiesResponse.Results {
			fmt.Println(r.Name)
		}
	}
	return err
}

func commandMapBack(config *Config) error {
	fmt.Println("Pokedex > map")

	if config.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	citiesResponse, err := GetData[CityResponse](config, config.Previous)
	if err == nil {
		config.Next = citiesResponse.Next
		config.Previous = citiesResponse.Previous
		for _, r := range citiesResponse.Results {
			fmt.Println(r.Name)
		}
	}
	return err
}

func commandExplore(config *Config) error {
	fmt.Printf("Exploring %v...\n", config.Input.param)
	data, err := GetData[NameResponse](config, baseurl+"location-area/"+config.Input.param)
	for _, d := range data.PokemonEncounters {
		fmt.Println("- ", d.Pokemon.Name)
	}
	return err
}

func commandCatch(config *Config) error {
	url := baseurl + "pokemon/" + config.Input.param

	fmt.Printf("Throwing a Pokeball at %v...\n", config.Input.param)
	data, err := GetData[PokemonRequest](config, url)
	if err != nil {
		fmt.Printf("Did npt found %v , Did you spell it correctly ?\n", config.Input.param)
		return nil
	}
	halfBE := float64(data.BaseExperience / 2)
	catched := float64(rand.Intn(data.BaseExperience)) > halfBE

	if catched {
		fmt.Printf("%v escaped!\n", data.Name)
	} else {
		fmt.Printf("You catched %v \n", data.Name)
		fmt.Println("You may now inspect it with the inspect command.")
		config.Pokemons[data.Name] = Pokemon{Name: data.Name, BaseExperience: data.BaseExperience, Height: data.Height, Weight: data.Weight, Types: data.Types}
	}
	return err
}

func commandInspect(config *Config) error {
	if pokemon, ok := config.Pokemons[config.Input.param]; ok {
		fmt.Println("Name: ", pokemon.Name)
		fmt.Println("Height: ", pokemon.Height)
		fmt.Println("Weight: ", pokemon.Weight)
		fmt.Println("BaseExperience: ", pokemon.BaseExperience)
		fmt.Println("Types: ")
		for _, val := range pokemon.Types {
			fmt.Println(" - ", val.Type.Name)
		}
	} else {
		fmt.Println("You have not caught the Pokemon (yet?)")
	}
	return nil
}

func commandPokedex(config *Config) error {
	fmt.Println("Your Pokedex:")
	for _, val := range config.Pokemons {
		fmt.Println(" - ", val.Name)
	}
	return nil
}
