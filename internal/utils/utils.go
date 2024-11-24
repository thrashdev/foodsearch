package utils

import (
	"log"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/thrashdev/foodsearch/internal/database"
	"github.com/thrashdev/foodsearch/internal/models"
)

func FloatToNumeric(number float64) (value pgtype.Numeric) {
	parsed := strconv.FormatFloat(number, 'f', -1, 64)
	if err := value.Scan(parsed); err != nil {
		log.Fatal("Error scanning numeric: %v", err)
	}
	return value
}

func DatabaseGlovoRestaurantToModel(dbRest database.GlovoRestaurant) models.GlovoRestaurant {
	deliveryFee, err := dbRest.DeliveryFee.Float64Value()
	if err != nil {
		log.Fatalf("Error converting pgtype.Numeric to float64, err:%v", err)
	}
	return models.GlovoRestaurant{
		Restaurant: models.Restaurant{
			ID:          dbRest.ID.Bytes,
			Name:        dbRest.Name,
			Address:     &dbRest.Address,
			DeliveryFee: &deliveryFee.Float64,
			PhoneNumber: &dbRest.PhoneNumber.String,
			CreatedAt:   dbRest.CreatedAt.Time,
			UpdatedAt:   dbRest.UpdatedAt.Time,
		},
		GlovoApiStoreID:   int(dbRest.GlovoApiStoreID),
		GlovoApiAddressID: int(dbRest.GlovoApiAddressID),
		GlovoApiSlug:      dbRest.GlovoApiSlug,
	}
}

func PrintErrors(errCh chan error) {
	for {
		err := <-errCh
		if err == nil {
			break
		}
		log.Println(err)
	}
}
