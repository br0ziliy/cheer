// Package flight_test
package flight_test

import (
	"log"
	"net/http/httptest"
	"testing"

	"github.com/arapov/cheer/lib/env"
	"github.com/arapov/cheer/lib/flight"
)

// TestRace tests for race conditions.
func TestRace(t *testing.T) {
	for i := 0; i < 100; i++ {
		go func() {
			// Load the configuration file
			config, err := env.LoadConfig("../../env.json.example")
			if err != nil {
				t.Fatal(err)
			}

			// Set up the session cookie store
			config.Session.SetupConfig()

			// Set up the views
			config.View.SetTemplates(config.Template.Root, config.Template.Children)

			// Store the view in flight
			flight.StoreConfig(*config)

			// Test the context retrieval
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "http://localhost/foo", nil)
			c := flight.Context(w, r)

			c.Config.Asset.Folder = "foo"
			log.Println(c.Config.Asset.Folder)

			c.View.BaseURI = "bar"
			log.Println(c.View.BaseURI)
		}()
	}
}
