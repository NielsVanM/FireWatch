package main

import (
	"github.com/nielsvanm/firewatch/web/views"

	"github.com/nielsvanm/firewatch/internal/database"
	"github.com/nielsvanm/firewatch/internal/server"
)

func main() {
	var s = server.NewServer(8000)
	database.NewDB("postgres", "Password8", "firewatch", "localhost", 5432)

	// Create server and endpoints
	unprotectedRouter := s.AddRouter("UnprotectedRouter", "/", nil)

	unprotectedRouter.AddEndpoint("/", views.DashboardView)

	s.Start()
}
