package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

/* Types */
/*Respone Types___________________________________________________________________________________________ */

type CityResponse struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type NameResponse struct {
	ID                   int64                 `json:"id"`
	Name                 string                `json:"name"`
	GameIndex            int64                 `json:"game_index"`
	EncounterMethodRates []EncounterMethodRate `json:"encounter_method_rates"`
	Location             Location              `json:"location"`
	Names                []Name                `json:"names"`
	PokemonEncounters    []PokemonEncounter    `json:"pokemon_encounters"`
}

type EncounterMethodRate struct {
	EncounterMethod Location                           `json:"encounter_method"`
	VersionDetails  []EncounterMethodRateVersionDetail `json:"version_details"`
}

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type EncounterMethodRateVersionDetail struct {
	Rate    int64    `json:"rate"`
	Version Location `json:"version"`
}

type Name struct {
	Name     string   `json:"name"`
	Language Location `json:"language"`
}

type PokemonEncounter struct {
	Pokemon        Location                        `json:"pokemon"`
	VersionDetails []PokemonEncounterVersionDetail `json:"version_details"`
}

type PokemonEncounterVersionDetail struct {
	Version          Location          `json:"version"`
	MaxChance        int64             `json:"max_chance"`
	EncounterDetails []EncounterDetail `json:"encounter_details"`
}

type EncounterDetail struct {
	MinLevel        int64         `json:"min_level"`
	MaxLevel        int64         `json:"max_level"`
	ConditionValues []interface{} `json:"condition_values"`
	Chance          int64         `json:"chance"`
	Method          Location      `json:"method"`
}

type Stats struct {
	BaseStat int `json:"base_stat"`
	Effort   int `json:"effort"`
}

type Type struct {
	Name string `json:"name"`
}
type Types struct {
	Type Type `json:"type"`
}

type PokemonRequest struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	BaseExperience int     `json:"base_experience"`
	Height         int     `json:"height"`
	Weight         int     `json:"weight"`
	Types          []Types `json:"types"`
}

/*Functions*/
func unmarshal[T any](data []byte) (T, error) {
	var result T
	if errDecode := json.Unmarshal(data, &result); errDecode != nil {
		return result, errDecode
	} else {

	}
	return result, nil
}

func fetch(url string, timeout time.Duration) ([]byte, error) {
	client := http.Client{Timeout: timeout}
	res, errGet := client.Get(url)

	if errGet != nil {
		return nil, errGet
	}

	defer res.Body.Close()
	resBody, errRead := io.ReadAll(res.Body)
	if errRead != nil {
		return nil, errGet
	}
	return resBody, nil
}

func GetData[T any](cfg *Config, url string) (T, error) {
	var result T
	if val, ok := cfg.history.Get(url); ok {
		return unmarshal[T](val)
	} else {
		data, errFetch := fetch(url, 2*time.Second)
		if errFetch != nil {
			return result, errFetch
		}
		cfg.history.Add(url, data)

		result, errUM := unmarshal[T](data)
		if errUM != nil {
			return result, errUM
		}
		return result, nil
	}

}

/*
	func decode[T any](r io.Reader, resultType *T) (T, error) {
		var result T
		if errDecode := json.NewDecoder(r).Decode(&resultType); errDecode != nil {
			return result, errDecode
		}

		return result, nil
	}
*/
