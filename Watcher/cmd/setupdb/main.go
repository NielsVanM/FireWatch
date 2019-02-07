package main

import (
	"github.com/nielsvanm/firewatch/internal/database"
	"github.com/nielsvanm/firewatch/internal/models"
)

func main() {
	database.NewDB("postgres", "Password8", "firewatch", "localhost", 5432)

	for _, query := range models.SetupQueries {
		database.Database.Query(query)
	}
}
