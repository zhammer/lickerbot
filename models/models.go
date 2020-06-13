package models

import (
	"log"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop/v5"
)

// DB is a connection to your database to be used
// throughout your application.
var DB *pop.Connection

func init() {
	var err error
	originalEnv := envy.Get("GO_ENV", "development")
	// to allow override database when running `buffalo dev`
	env := envy.Get("GO_ENV_OVERRIDE", originalEnv)
	DB, err = pop.Connect(env)
	if err != nil {
		log.Fatal(err)
	}
	pop.Debug = env == "development" || env == "test"
}
