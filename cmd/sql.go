package cmd

import (
	"fmt"

	sql "github.com/aodin/aspect"

	db "github.com/aodin/listofthings/db"
)

// TODO iteration order cannot be guaranteed
var applications = map[string][]*sql.TableElem{
	"auth": {
		db.Users,
		db.Sessions,
	},
	"things": {
		db.Things,
	},
}

func SQL(all bool, apps ...string) {
	if all {
		AllApps()
		return
	}
	if len(apps) == 0 {
		ListApps()
		return
	}
	var invalid bool
	for _, app := range apps {
		invalid = invalid || !HasApp(app)
	}
	if invalid {
		return
	}
	for _, app := range apps {
		PrintApp(app, applications[app])
	}
}

func AllApps() {
	for name, tables := range applications {
		PrintApp(name, tables)
	}
}

func PrintApp(app string, tables []*sql.TableElem) {
	fmt.Printf("-- %s\n", app)
	for _, table := range tables {
		fmt.Println(table.Create())
	}
}

func ListApps() {
	fmt.Println("Available applications:")
	for key := range applications {
		fmt.Printf(" * %s\n", key)
	}
}

func HasApp(app string) bool {
	_, ok := applications[app]
	if !ok {
		fmt.Printf(" * '%s' is not a valid application name\n", app)
		return false
	}
	return true
}
