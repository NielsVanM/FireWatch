package main

import (
	"github.com/nielsvanm/firewatch/web/views"

	"github.com/nielsvanm/firewatch/internal/database"
	"github.com/nielsvanm/firewatch/internal/middlware"
	"github.com/nielsvanm/firewatch/internal/server"
)

func main() {
	var s = server.NewServer(8000)
	database.NewDB("postgres", "Password8", "firewatch", "localhost", 5432)

	// Create server and endpoints
	unprotectedRouter := s.AddRouter("UnprotectedRouter", "/")
	unprotectedRouter.AddEndpoint("/auth/login/", views.LoginView)

	protectedRouter := s.AddRouter("ProtectedRouter", "/")
	protectedRouter.AddMiddlewware(middlware.AuthorizationMiddleware)
	protectedRouter.AddEndpoint("/", views.DashboardView)

	s.Start()
}
