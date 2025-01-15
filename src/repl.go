package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func startRepl(config *Config) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := cleanInput(scanner.Text())
		if len(input) == 0 {
			continue
		}
		config.Input.input = input
		config.Input.command = input[0]
		if len(input) > 1 {
			config.Input.param = input[1]
		}
		err := findCommand(config.Input.command, config)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Your command was: %s\n", input[0])

	}
}

func cleanInput(text string) []string {
	loweredText := strings.ToLower(text)
	result := strings.Fields(loweredText)
	return result
}

func findCommand(input string, config *Config) error {
	command, ok := commands[input]
	if ok {
		err := command.callback(config)
		return err

	} else {
		return errors.New("Unknown command")
	}
}
