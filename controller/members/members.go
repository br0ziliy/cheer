// Package members
package members

import (
	"encoding/json"
	"net/http"

	"github.com/arapov/cheer/lib/flight"
	"github.com/arapov/cheer/middleware/acl"
	"github.com/arapov/cheer/model/group"
	"github.com/arapov/cheer/model/member"

	"github.com/blue-jay/core/router"
)

var (
	uri = "/members"
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

	items, _, err := member.All(c.DB)
	if err != nil {
		c.FlashError(err)
		items = []member.Item{}
	}

	v := c.View.New("member/index")
	v.Vars["items"] = items
	v.Render(w, r)
}

// Create displays the create form.
func Create(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	items, _, err := group.GetGroupIdName(c.DB)
	if err != nil {
		c.FlashError(err)
		items = []group.SanitizedItem{}
	}

	b, err := json.Marshal(items)
	if err != nil {
		c.FlashError(err)
		b = []byte(`{}`)
	}

	v := c.View.New("member/create")
	c.Repopulate(v.Vars, "ircnick")
	c.Repopulate(v.Vars, "fullname")
	c.Repopulate(v.Vars, "group_id")
	c.Repopulate(v.Vars, "groupname")
	v.Vars["json"] = b
	v.Render(w, r)
}

// Store handles the create form submission.
func Store(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	if !c.FormValid("ircnick") {
		Create(w, r)
		return
	}

	_, err := member.Create(c.DB, r.FormValue("ircnick"), r.FormValue("fullname"), r.FormValue("group_id"))
	if err != nil {
		c.FlashError(err)
		Create(w, r)
		return
	}

	c.FlashSuccess("Member added.")
	c.Redirect(uri)
}

// Show displays a single item.
func Show(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	item, _, err := member.ByID(c.DB, c.Param("id"))
	if err != nil {
		c.FlashError(err)
		c.Redirect(uri)
		return
	}

	v := c.View.New("member/show")
	v.Vars["item"] = item
	v.Render(w, r)
}

// Edit displays the edit form.
func Edit(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	// TODO: unify the code below with the same from Create()
	//       I don't know the langugage good enough even to create new func()
	//       And yes, I'm in a hurry. :)
	items, _, err := group.GetGroupIdName(c.DB)
	if err != nil {
		c.FlashError(err)
		items = []group.SanitizedItem{}
	}

	b, err := json.Marshal(items)
	if err != nil {
		c.FlashError(err)
		b = []byte(`{}`)
	}

	item, _, err := member.ByID(c.DB, c.Param("id"))
	if err != nil {
		c.FlashError(err)
		c.Redirect(uri)
		return
	}

	v := c.View.New("member/edit")
	c.Repopulate(v.Vars, "ircnick")
	c.Repopulate(v.Vars, "fullname")
	c.Repopulate(v.Vars, "group_id")
	c.Repopulate(v.Vars, "groupname")
	v.Vars["json"] = b
	v.Vars["item"] = item
	v.Render(w, r)
}

// Update handles the edit form submission.
func Update(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	if !c.FormValid("ircnick") {
		Edit(w, r)
		return
	}

	_, err := member.Update(c.DB, r.FormValue("ircnick"), r.FormValue("fullname"), r.FormValue("group_id"), c.Param("id"))
	if err != nil {
		c.FlashError(err)
		Edit(w, r)
		return
	}

	c.FlashSuccess("Member updated.")
	c.Redirect(uri)
}

// Destroy handles the delete form submission.
func Destroy(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	_, err := member.DeleteHard(c.DB, c.Param("id"))
	if err != nil {
		c.FlashError(err)
	} else {
		c.FlashNotice("Member deleted.")
	}

	c.Redirect(uri)
}
