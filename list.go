package main

import (
	"log"
	"os"
	"path/filepath"

	sql "github.com/aodin/aspect"
	"github.com/aodin/volta/config"
	"github.com/codegangsta/cli"

	"github.com/aodin/listofthings/cmd"
	"github.com/aodin/listofthings/server"
)

func main() {
	app := cli.NewApp()
	app.Name = "quilt"
	app.Usage = "Start the Quilt server"
	app.Action = startServer
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log, l",
			Value: "",
			Usage: "Sets the log output file path",
		},
		cli.StringFlag{
			Name:  "config, c",
			Value: "./settings.json",
			Usage: "Sets the configuration file",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "sql",
			Usage: "output SQL for the given applications",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "all, a",
					Usage: "output the sql of all applications",
				},
			},
			Action: func(c *cli.Context) {
				cmd.SQL(c.Bool("all"), c.Args()...)
			},
		},
	}
	app.Run(os.Args)
}

func setUp(file string) (*sql.DB, config.Config) {
	// Parse the given configuration file
	conf, err := config.ParseFile(file)
	if err != nil {
		log.Panicf("quilt: could not parse configuration: %s", err)
	}

	// Connect to the database
	conn, err := sql.Connect(conf.Database.Driver, conf.Database.Credentials())
	if err != nil {
		log.Panicf("quilt: could not connect to the database: %s", err)
	}
	return conn, conf
}

func startServer(c *cli.Context) {
	logF := c.String("log")
	file := c.String("config")
	// Set the log output - if no path given, use stdout
	// TODO log rotation?
	if logF != "" {
		if err := os.MkdirAll(filepath.Dir(logF), 0776); err != nil {
			log.Panic(err)
		}
		l, err := os.OpenFile(logF, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Panic(err)
		}
		defer l.Close()
		log.SetOutput(l)
	}
	conn, conf := setUp(file)
	defer conn.Close()
	log.Println("Starting server on", conf.Address())
	log.Panic(server.New(conf, conn).ListenAndServe())
}
