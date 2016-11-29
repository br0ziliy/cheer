// Package home displays the Home page.
package home

import (
	"encoding/json"
	"net/http"

	"github.com/arapov/cheer/lib/flight"
	"github.com/arapov/cheer/model/member"

	"github.com/blue-jay/core/router"
)

// Load the routes.
func Load() {
	router.Get("/", Index)
}

// Index displays the home page.
func Index(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	items, _, err := member.GetNicks(c.DB)
	if err != nil {
		c.FlashError(err)
		items = []member.SanitizedItem{}
	}

	b, err := json.Marshal(items)
	if err != nil {
		c.FlashError(err)
		b = []byte(`{}`)
	}

	v := c.View.New("home/index")
	c.Repopulate(v.Vars, "ircnick")
	v.Vars["json"] = b
	v.Render(w, r)
}
