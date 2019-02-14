package main

import (
	"github.com/nielsvanm/firewatch/views"

	"github.com/nielsvanm/firewatch/internal/database"
	"github.com/nielsvanm/firewatch/internal/middleware"
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
	protectedRouter.AddEndpoint("/account/", views.AccountView)
	protectedRouter.AddEndpoint("/account/all-device-logout/", views.LogOutAllDevicesView)
	protectedRouter.AddEndpoint("/account/change-password/", views.ChangePasswordView)
	protectedRouter.AddEndpoint("/account/delete/", views.DeleteAccountView)

	s.Start()
}
