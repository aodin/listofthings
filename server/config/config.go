package config

import "flag"

type Config struct {
	Port      int    `json:"port"`
	Templates string `json:"templates"`
	Static    string `json:"static"`
}

// Read command line variables and generate a config file
func Parse() Config {
	// TODO Parse a json config file contained in the flag -config
	var config Config
	flag.IntVar(&config.Port, "port", 9001, "port for the HTTP server")
	flag.StringVar(
		&config.Templates,
		"templates",
		"./templates",
		"template directory",
	)
	flag.StringVar(
		&config.Static,
		"static",
		"./dist",
		"static directory",
	)
	flag.Parse()
	return config
}
