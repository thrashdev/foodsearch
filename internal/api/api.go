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
	ID         uuid.UUID                `json:"id"`
	GlovoRest  *models.GlovoRestaurant  `json:"glovo_restaurant"`
	YandexRest *models.YandexRestaurant `json:"yandex_restaurant"`
}

func (cfg *ApiConfig) HandlerGetRestaurants(w http.ResponseWriter, r *http.Request) {
	cfg.Logger.Info("Serving all restaurants on /v1/restaurants")
	ctx := context.Background()
	restaurantRows, err := cfg.DB.GetAllRestaurants(ctx)
	if err != nil {
		cfg.Logger.Panicw("Error fetching restaurants from DB", "error", err)
	}
	restaurants := []Restaurant{}
	for _, row := range restaurantRows {
		restaurants = append(restaurants, serializeRestaurant(row))
	}

	utils.RespondWithJSON(w, 200, restaurants)
}

func serializeRestaurant(row database.GetAllRestaurantsRow) Restaurant {
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
			GlovoApiStoreID:   int(row.GlovoApiStoreID.Int32),
			GlovoApiAddressID: int(row.GlovoApiAddressID.Int32),
			GlovoApiSlug:      row.GlovoApiSlug.String,
		}
	}

	yandexAddress := ""
	if row.Address_2.Valid {
		yandexAddress = row.Address_2.String
	}

	yandexDeliveryFee := 0.0
	floatValue, err := row.DeliveryFee_2.Float64Value()
	if err != nil {
		if floatValue.Valid {
			yandexDeliveryFee = floatValue.Float64
		}
	}

	yandexPhoneNumber := ""
	if row.PhoneNumber_2.Valid {
		yandexPhoneNumber = row.PhoneNumber.String
	}

	var yandexRest *models.YandexRestaurant
	yandexRest = nil

	if row.ID_3.Valid {
		yandexRest = &models.YandexRestaurant{
			Restaurant: models.Restaurant{
				ID:          row.ID_3.Bytes,
				Name:        row.Name_2.String,
				Address:     &yandexAddress,
				DeliveryFee: &yandexDeliveryFee,
				PhoneNumber: &yandexPhoneNumber,
			},
			YandexApiSlug: row.YandexApiSlug.String,
		}
	}

	rest := Restaurant{
		ID:         row.ID.Bytes,
		GlovoRest:  glovoRest,
		YandexRest: yandexRest,
	}

	return rest

}
