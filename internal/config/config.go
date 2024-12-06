package config

import (
	"log/slog"

	"github.com/thrashdev/foodsearch/internal/database"
	"go.uber.org/zap"
)

type GlovoConfig struct {
	SearchURL  string
	FiltersURL string
	DishURL    string
}

type YandexConfig struct {
	SearchURL         string
	RestaurantMenuURL string
	Loc               YandexLocation
}

type YandexLocation struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type Config struct {
	Glovo           GlovoConfig
	Yandex          YandexConfig
	UpdateBatchSize int
	DB              database.Queries
	Logger          *zap.SugaredLogger
}
