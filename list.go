package main

import (
	"log"
	"time"

	"github.com/aodin/listofthings/db"
	"github.com/aodin/listofthings/server"
	"github.com/aodin/volta/config"
)

var conf = config.Config{
	Port:        9001,
	ProxyPort:   9000,
	TemplateDir: "./templates",
	StaticDir:   "./dist",
	StaticURL:   "/dist/",
	Cookie: config.CookieConfig{
		Age:      365 * 24 * time.Hour,
		Domain:   "",
		HttpOnly: false,
		Name:     "listuserid",
		Path:     "/",
		Secure:   false,
	},
}

func main() {
	// Create a memory store with a limited number of items
	store := server.NewMemoryStore(25)

	// Add a few things
	names := []string{"Bass-o-matic", "Swill", "Jam Hawkers"}

	for _, name := range names {
		thing := db.NewThing(name)
		if _, err := store.Create(&thing); err != nil {
			log.Panic(err)
		}
	}

	// Create a new server
	log.Println("Starting server on", conf.Address())
	log.Panic(server.New(conf, store).ListenAndServe())
}
