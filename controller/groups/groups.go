// Package groups
package groups

import (
	"net/http"

	"github.com/arapov/cheer/lib/flight"
	"github.com/arapov/cheer/middleware/acl"
	"github.com/arapov/cheer/model/group"

	"github.com/blue-jay/core/router"
)

var (
	uri = "/groups"
)

// Load the routes.
func Load() {
	c := router.Chain(acl.DisallowAnon)
	router.Get(uri, Index, c...)
	router.Get(uri+"/create", Create, c...)
	router.Post(uri+"/create", Store, c...)
	router.Get(uri+"/view/:id", Show, c...)
	router.Get(uri+"/edit/:id", Edit, c...)
	router.Patch(uri+"/edit/:id", Update, c...)
	router.Delete(uri+"/:id", Destroy, c...)
}

// Index displays the items.
func Index(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	items, _, err := group.All(c.DB)
	if err != nil {
		c.FlashError(err)
		items = []group.Item{}
	}

	v := c.View.New("group/index")
	v.Vars["items"] = items
	v.Render(w, r)
}

// Create displays the create form.
func Create(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	v := c.View.New("group/create")
	c.Repopulate(v.Vars, "name")
	c.Repopulate(v.Vars, "description")
	v.Render(w, r)
}

// Store handles the create form submission.
func Store(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	if !c.FormValid("name") {
		Create(w, r)
		return
	}

	_, err := group.Create(c.DB, r.FormValue("name"), r.FormValue("description"))
	if err != nil {
		c.FlashError(err)
		Create(w, r)
		return
	}

	c.FlashSuccess("Group added.")
	c.Redirect(uri)
}

// Show displays a single item.
func Show(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	item, _, err := group.ByID(c.DB, c.Param("id"))
	if err != nil {
		c.FlashError(err)
		c.Redirect(uri)
		return
	}

	v := c.View.New("group/show")
	v.Vars["item"] = item
	v.Render(w, r)
}

// Edit displays the edit form.
func Edit(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	item, _, err := group.ByID(c.DB, c.Param("id"))
	if err != nil {
		c.FlashError(err)
		c.Redirect(uri)
		return
	}

	v := c.View.New("group/edit")
	c.Repopulate(v.Vars, "name")
	c.Repopulate(v.Vars, "description")
	v.Vars["item"] = item
	v.Render(w, r)
}

// Update handles the edit form submission.
func Update(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	if !c.FormValid("name") {
		Edit(w, r)
		return
	}

	_, err := group.Update(c.DB, r.FormValue("name"), r.FormValue("description"), c.Param("id"))
	if err != nil {
		c.FlashError(err)
		Edit(w, r)
		return
	}

	c.FlashSuccess("Group updated.")
	c.Redirect(uri)
}

// Destroy handles the delete form submission.
func Destroy(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	_, err := group.DeleteHard(c.DB, c.Param("id"))
	if err != nil {
		c.FlashError(err)
	} else {
		c.FlashNotice("Item deleted.")
	}

	c.Redirect(uri)
}
