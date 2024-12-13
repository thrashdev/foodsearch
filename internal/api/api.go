package api

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/thrashdev/foodsearch/internal/database"
	"github.com/thrashdev/foodsearch/internal/models"
	"github.com/thrashdev/foodsearch/internal/utils"
	"go.uber.org/zap"
)

type ApiConfig struct {
	DB     database.Queries
	Logger *zap.SugaredLogger
}

type Restaurant struct {
	ID         uuid.UUID
	GlovoRest  *models.GlovoRestaurant
	YandexRest *models.YandexRestaurant
}

func (cfg *ApiConfig) handlerGetRestaurants(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	restaurantRows, err := cfg.DB.GetAllRestaurants(ctx)
	if err != nil {
		cfg.Logger.Panicw("Error fetching restaurants from DB", "error", err)
	}
}

func serializeRestaurants(rows []database.GetAllRestaurantsRow) []Restaurant {
	result := []Restaurant{}
	for _, row := range rows {
		var glovoRest *models.GlovoRestaurant
		glovoRest = nil
		glovoDeliveryFee, err := row.DeliveryFee.Float64Value()
		if err != nil {
			log.Fatalf("Error converting pgtype.Numeric to float64, err:%v", err)
		}
		if row.ID_2.Valid {
			glovoRest = &models.GlovoRestaurant{
				Restaurant: models.Restaurant{
					ID:          row.ID_2.Bytes,
					Name:        row.Name.String,
					Address:     &row.Address.String,
					DeliveryFee: &glovoDeliveryFee.Float64,
					PhoneNumber: &row.PhoneNumber.String,
				},
			}
		}
		rest := Restaurant{
			ID:        row.ID.Bytes,
			GlovoRest: glovoRest,
		}
		result = append(result)
	}

}
