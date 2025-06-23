package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	next string
	prev string
}

type Locations struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

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

func startREPL() {
	config := config{
		next: "https://pokeapi.co/api/v2/location-area//",
		prev: "",
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
			err := structure.callback(&config)
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

func commandExit(pointer *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(pointer *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:\n")
	for _, v := range getCommands() {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}
	return nil
}

func commandMap(pointer *config) error {
	res, err := http.Get(pointer.next)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return err
	}
	locations := Locations{}
	err = json.Unmarshal(body, &locations)
	if err != nil {
		return err
	}
	pointer.next = locations.Next
	pointer.prev = locations.Previous
	for _, v := range locations.Results {
		fmt.Println(v.Name)
	}
	return nil
}

func commandMapB(pointer *config) error {
	if pointer.prev == "" {
		return errors.New("First page")
	}
	res, err := http.Get(pointer.prev)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return err
	}
	locations := Locations{}
	err = json.Unmarshal(body, &locations)
	if err != nil {
		return err
	}
	pointer.next = locations.Next
	pointer.prev = locations.Previous
	for _, v := range locations.Results {
		fmt.Println(v.Name)
	}
	return nil
}
