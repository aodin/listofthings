package main

import (
	"log"

	"github.com/aodin/listofthings/db"
	"github.com/aodin/listofthings/server"
	"github.com/aodin/listofthings/server/config"
)

func main() {
	// Create a memory store with a limited number of items
	store := server.NewMemoryStore(25)

	// Create and parse the server configuration
	config := config.Parse()

	// Add a few things
	names := []string{"Bass-o-matic", "Swill", "Jam Hawkers"}

	for _, name := range names {
		thing := db.NewThing(name)
		if _, err := store.Create(&thing); err != nil {
			log.Panic(err)
		}
	}

	// Create a new server
	srv, err := server.New(config, store)
	if err != nil {
		log.Panic(err)
	}

	log.Println("Starting server on port:", config.Port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
