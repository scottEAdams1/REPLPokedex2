package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/scottEAdams1/REPLPokedex2/internal/pokecache"
)

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays next 20 locations from the map",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 locations from the map",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore <location name>",
			description: "Displays available pokemon in an area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch <pokemon name>",
			description: "Chance to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect <pokemon name>",
			description: "Inspect a caught pokemon",
			callback:    commandInspect,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}

func cleanInput(text string) []string {
	loweredText := strings.ToLower(text)
	slice := strings.Fields(loweredText)
	return slice
}

func startREPL(cache pokecache.Cache) {
	config := config{
		next:   "https://pokeapi.co/api/v2/location-area//",
		prev:   "",
		cache:  cache,
		caught: map[string]Pokemon{},
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		words := cleanInput(input)
		if len(words) == 0 {
			continue
		}
		command := words[0]
		if structure, ok := getCommands()[command]; ok {
			err := structure.callback(&config, words[1:])
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}
