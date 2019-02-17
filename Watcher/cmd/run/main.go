package main

import (
	"github.com/nielsvanm/firewatch/core/config"
	"github.com/nielsvanm/firewatch/core/database"
	"github.com/nielsvanm/firewatch/core/middleware"
	"github.com/nielsvanm/firewatch/core/server"
	"github.com/nielsvanm/firewatch/routes"
)

func main() {

	cfg := config.LoadConfig("./configs/config.json")

	// Connect to DB
	database.NewDB(cfg.Database.Username, cfg.Database.Password, cfg.Database.Name, cfg.Database.Host, cfg.Database.Port)

	// Create server and endpoints
	var s = server.NewServer(cfg.Server.Port)

	unprotectedRouter := s.AddRouter("UnprotectedRouter", "/")
	unprotectedRouter.AddMiddlewware(middleware.HTTPLogMiddleware)
	unprotectedRouter.ParseRouteMap(routes.UnprotectedRoutes)

	protectedRouter := s.AddRouter("ProtectedRouter", "/")
	protectedRouter.AddMiddlewware(middleware.HTTPLogMiddleware)
	protectedRouter.AddMiddlewware(middleware.AuthorizationMiddleware)
	protectedRouter.ParseRouteMap(routes.ProtectedRoutes)

	s.Start()
}
