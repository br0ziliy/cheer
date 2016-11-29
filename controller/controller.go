// Package controller loads the routes for each of the controllers.
package controller

import (
	"github.com/arapov/cheer/controller/about"
	"github.com/arapov/cheer/controller/debug"
	"github.com/arapov/cheer/controller/home"
	"github.com/arapov/cheer/controller/login"
	"github.com/arapov/cheer/controller/register"
	"github.com/arapov/cheer/controller/static"
	"github.com/arapov/cheer/controller/status"
)

// LoadRoutes loads the routes for each of the controllers.
func LoadRoutes() {
	about.Load()
	debug.Load()
	register.Load()
	login.Load()
	home.Load()
	static.Load()
	status.Load()
}
