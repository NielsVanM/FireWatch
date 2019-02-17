package main

import (
	"github.com/nielsvanm/firewatch/views"

	"github.com/nielsvanm/firewatch/core/config"
	"github.com/nielsvanm/firewatch/core/database"
	"github.com/nielsvanm/firewatch/core/middleware"
	"github.com/nielsvanm/firewatch/core/server"
)

func main() {

	cfg := config.LoadConfig("./configs/config.json")
	serverCfg := cfg.Server
	dbCfg := cfg.Database

	var s = server.NewServer(serverCfg.Port)
	database.NewDB(dbCfg.Username, dbCfg.Password, dbCfg.Name, dbCfg.Host, dbCfg.Port)

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
