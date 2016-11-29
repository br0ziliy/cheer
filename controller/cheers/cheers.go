// Package cheers
package cheers

import (
	"encoding/json"
	"net/http"

	"github.com/arapov/cheer/lib/flight"
	"github.com/arapov/cheer/model/cheer"
	"github.com/arapov/cheer/model/member"

	"github.com/blue-jay/core/router"
)

var (
	uri = "/cheers"
)

// Load the routes.
func Load() {
	c := router.Chain()
	router.Post(uri, Index, c...)
	router.Post(uri+"/submit", Store, c...)
	router.Get(uri+"/submit", Submit, c...)
}

// Index displays the items.
func Index(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)
	ircnick := r.PostFormValue("ircnick") // todo: lib/flight must do it
	referer := r.RemoteAddr

	group, _, err := member.GetMyGroup(c.DB, ircnick)
	if err != nil {
		c.FlashNotice("The " + ircnick + " nick does not exist.")
		//c.FlashError(err)
		group = member.Item{}
		c.Redirect("/")
	}

	items, _, err := member.GetMates(c.DB, ircnick)
	if err != nil {
		c.FlashError(err)
		items = []member.Item{}
	}

	nicks, _, err := member.GetNicks(c.DB)
	if err != nil {
		c.FlashError(err)
		nicks = []member.SanitizedItem{}
	}

	b, err := json.Marshal(nicks)
	if err != nil {
		c.FlashError(err)
		b = []byte(`{}`)
	}

	v := c.View.New("cheer/index")
	v.Vars["ircnick"] = ircnick
	v.Vars["referer"] = referer
	v.Vars["items"] = items
	v.Vars["group"] = group
	v.Vars["json"] = b
	v.Render(w, r)
}

func Submit(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	v := c.View.New("cheer/submit")
	v.Render(w, r)
}

// Store handles the create form submission.
func Store(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)
	referer := r.RemoteAddr

	// TODO: check how stupid is the stuff below
	x := make(map[string][]string)
	for key, values := range r.Form {
		for _, value := range values {
			x[key] = append(x[key], value)
		}
	}
	for i := 0; i < len(x["to_id"]); i++ {
		_, err := cheer.Create(c.DB, r.PostFormValue("from_id"), x["to_id"][i], x["points"][i], x["message"][i], referer)
		if err != nil {
			c.FlashError(err)
			c.Redirect(uri)
			return
		}
	}

	c.FlashSuccess("Cheers added.")
	c.Redirect(uri + "/submit")
}
