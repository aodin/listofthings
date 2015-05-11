package main

import (
	"log"
	"time"

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
	// Create a new server
	log.Println("Starting server on", conf.Address())
	log.Panic(server.New(conf).ListenAndServe())
}
