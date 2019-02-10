package main

import (
	"github.com/nielsvanm/firewatch/views"

	"github.com/nielsvanm/firewatch/internal/database"
	middleware "github.com/nielsvanm/firewatch/internal/middlware"
	"github.com/nielsvanm/firewatch/internal/server"
)

func main() {
	var s = server.NewServer(8000)
	database.NewDB("postgres", "Password8", "firewatch", "localhost", 5432)

	// Create server and endpoints
	unprotectedRouter := s.AddRouter("UnprotectedRouter", "/")
	unprotectedRouter.AddMiddlewware(middleware.HTTPLogMiddleware)
	unprotectedRouter.AddEndpoint("/auth/login/", views.LoginView)
	unprotectedRouter.AddEndpoint("/auth/logout/", views.LogoutView)

	protectedRouter := s.AddRouter("ProtectedRouter", "/")
	protectedRouter.AddMiddlewware(middleware.HTTPLogMiddleware)
	protectedRouter.AddMiddlewware(middleware.AuthorizationMiddleware)
	protectedRouter.AddEndpoint("/", views.DashboardView)

	s.Start()
}
