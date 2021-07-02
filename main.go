package main

import (
	"github.com/Project-Planner/backend/config"
	"github.com/Project-Planner/backend/web"
	"github.com/Project-Planner/backend/xmldb"
	"log"
)

// main is the entry point of the application and basically an "Avengers assemble!"
func main() {
	log.Println("Starting up ...")

	// Load config
	c, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Create database implementation
	db, err := xmldb.NewDatabase(c.DBConfig)
	if err != nil {
		log.Fatal(err)
	}
	db.AddUser("Lukas", "bla")

	// Start web server
	web.ListenAndServe(db, c.ServerConfig)
}
