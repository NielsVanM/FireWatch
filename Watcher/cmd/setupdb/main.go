package main

import (
	"github.com/nielsvanm/firewatch/core/database"
	"github.com/nielsvanm/firewatch/core/models"
)

func main() {
	database.NewDB("postgres", "Password8", "firewatch", "localhost", 5432)

	for _, query := range models.SetupQueries {
		database.DB.Query(query)
	}
}
