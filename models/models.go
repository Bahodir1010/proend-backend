package models

import (
	"log"
	"os"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop/v6"
)

// DB is a connection to your database to be used
// throughout your application.
var DB *pop.Connection

func init() {
	var err error
	env := envy.Get("GO_ENV", "development")
	
	// ===================================================================
	// === THIS IS THE CORRECTED LOGIC ===
	// ===================================================================
	// First, try to get the database URL from the environment variable set by Docker Compose.
	dbUrl := os.Getenv("DATABASE_URL")

	if dbUrl != "" {
		// If DATABASE_URL is set, use it directly. This is for Docker.
		log.Printf("Connecting to database using DATABASE_URL environment variable.")
		DB, err = pop.NewConnection(&pop.ConnectionDetails{
			URL: dbUrl,
		})
	} else {
		// If DATABASE_URL is NOT set, fall back to the old method. This is for running outside of Docker.
		log.Printf("DATABASE_URL not found. Connecting using pop.Connect with env: %s", env)
		DB, err = pop.Connect(env)
	}
	// ===================================================================

	if err != nil {
		log.Fatal(err)
	}
	pop.Debug = env == "development"
}