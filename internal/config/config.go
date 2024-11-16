package config

import "github.com/thrashdev/foodsearch/internal/database"

type GlovoConfig struct {
	SearchURL  string
	FiltersURL string
	DishURL    string
}

type Config struct {
	Glovo GlovoConfig
	DB    database.Queries
}
