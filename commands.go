package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
)

func commandExit(pointer *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(pointer *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, v := range getCommands() {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}
	return nil
}

func commandMap(pointer *config, args []string) error {
	var body []byte
	if val, ok := pointer.cache.Get(pointer.next); ok {
		body = val
	} else {
		res, err := http.Get(pointer.next)
		if err != nil {
			return err
		}
		body, err = io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return err
		}
		pointer.cache.Add(pointer.next, body)
	}
	locations := Locations{}
	err := json.Unmarshal(body, &locations)
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

func commandMapB(pointer *config, args []string) error {
	if pointer.prev == "" {
		return errors.New("First page")
	}
	var body []byte
	if val, ok := pointer.cache.Get(pointer.prev); ok {
		body = val
	} else {
		res, err := http.Get(pointer.prev)
		if err != nil {
			return err
		}
		body, err = io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return err
		}
	}
	locations := Locations{}
	err := json.Unmarshal(body, &locations)
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

func commandExplore(pointer *config, args []string) error {
	if len(args) == 0 {
		return errors.New("Provide a location name")
	}
	var body []byte
	if val, ok := pointer.cache.Get(args[0]); ok {
		body = val
	} else {
		res, err := http.Get("https://pokeapi.co/api/v2/location-area/" + args[0] + "/")
		if err != nil {
			return err
		}
		body, err = io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return err
		}
		pointer.cache.Add(args[0], body)
	}
	pokemon := PokemonArea{}
	err := json.Unmarshal(body, &pokemon)
	if err != nil {
		return err
	}
	fmt.Println("Exploring" + args[0] + "\nFound:")
	for _, v := range pokemon.PokemonEncounters {
		fmt.Println(v.Pokemon.Name)
	}
	return nil
}

func commandCatch(pointer *config, args []string) error {
	if len(args) == 0 {
		return errors.New("Provide a pokemon")
	}
	var body []byte
	if val, ok := pointer.cache.Get(args[0]); ok {
		body = val
	} else {
		res, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + args[0] + "/")
		if err != nil {
			return err
		}
		body, err = io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return err
		}
		pointer.cache.Add(args[0], body)
	}
	pokemon := Pokemon{}
	err := json.Unmarshal(body, &pokemon)
	if err != nil {
		return err
	}
	catch := rand.Intn(pokemon.BaseExperience)
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	if catch > 40 {
		fmt.Printf("%s escaped!\n", pokemon.Name)
		return nil
	}
	fmt.Printf("%s was caught!\n", pokemon.Name)
	pointer.caught[pokemon.Name] = pokemon
	return nil
}

func commandInspect(pointer *config, args []string) error {
	if pokemon, ok := pointer.caught[args[0]]; ok {
		fmt.Printf("Name: %s\nHeight: %v\nWeight: %v\n", pokemon.Name, pokemon.Height, pokemon.Weight)
		fmt.Println("Stats:")
		for _, v := range pokemon.Stats {
			fmt.Printf("\t-%s: %v\n", v.Stat.Name, v.BaseStat)
		}
		fmt.Println("Types:")
		for _, v := range pokemon.Types {
			fmt.Printf("\t- %s\n", v.Type.Name)
		}
		return nil
	}
	fmt.Println("you have not caught that pokemon")
	return nil
}
