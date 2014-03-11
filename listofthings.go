package main

import (
	"github.com/aodin/listofthings/server"
	"log"
)

func main() {
	// Create and parse the server configuration
	config := server.Parse()
	log.Println("Config:", config)

	// Create a memory store with a limited number of items
	memory := server.NewMemoryStore(25)

	// Add a few things
	names := []string{"Bass-o-matic", "Swill", "Jam Hawkers"}

	var err error
	for _, name := range names {
		_, err = memory.Create(&server.Thing{Name: name})
		if err != nil {
			panic(err)
		}
	}
	log.Println(memory.List())

	// Create a new server
	s, err := server.New(config, memory)
	if err != nil {
		panic(err)
	}

	log.Println("Starting server")
	err = s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
