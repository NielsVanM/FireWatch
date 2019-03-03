package main

import (
	"flag"
	"os"

	"github.com/nielsvanm/firewatch/core/config"
	"github.com/nielsvanm/firewatch/core/database"
	"github.com/nielsvanm/firewatch/core/middleware"
	"github.com/nielsvanm/firewatch/core/models"
	"github.com/nielsvanm/firewatch/core/server"
	"github.com/nielsvanm/firewatch/routes"

	log "github.com/sirupsen/logrus"
)

func main() {

	cfg := config.LoadConfig("./configs/config.json")

	// Connect to DB
	database.NewDB(cfg.Database.Username, cfg.Database.Password, cfg.Database.Name, cfg.Database.Host, cfg.Database.Port)

	Flags()

	// Create server and endpoints
	var s = server.NewServer(cfg.Server.Port)
	s.SetStaticDir(cfg.Server.StaticDir)

	// Add logging middleware
	s.MasterRouter.Use(middleware.HTTPLogMiddleware)

	unprotectedRouter := s.AddRouter("UnprotectedRouter", "/")
	unprotectedRouter.ParseRouteMap(routes.UnprotectedRoutes)

	protectedRouter := s.AddRouter("ProtectedRouter", "/")
	protectedRouter.AddMiddleware(middleware.AuthorizationMiddleware)
	protectedRouter.ParseRouteMap(routes.ProtectedRoutes)

	s.Start()
}

// Flags parses necessarry command line flags
func Flags() {
	createAdmin := flag.String("CreateAdmin", "", "Creates a admin account with the provided password")

	flag.Parse()

	if *createAdmin != "" {
		a := models.NewAccount("admin", *createAdmin)
		a.Save()

		log.Info("Created administator user with password " + *createAdmin)

		os.Exit(0)
	}
}
