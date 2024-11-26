package utils

import (
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
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

func DatabaseYandexRestaurantToModel(dbRest database.YandexRestaurant) models.YandexRestaurant {
	addr := ""
	if dbRest.Address.Valid {
		addr = dbRest.Address.String
	}

	phoneNumber := ""
	if dbRest.PhoneNumber.Valid {
		phoneNumber = dbRest.PhoneNumber.String
	}

	dbFloat, err := dbRest.DeliveryFee.Float64Value()
	deliveryFee := 0.0
	if err != nil {
		if dbFloat.Valid {
			deliveryFee = dbFloat.Float64
		}
	}

	return models.YandexRestaurant{
		Restaurant: models.Restaurant{
			ID:          dbRest.ID.Bytes,
			Name:        dbRest.Name,
			Address:     &addr,
			DeliveryFee: &deliveryFee,
			PhoneNumber: &phoneNumber,
			CreatedAt:   dbRest.CreatedAt.Time,
			UpdatedAt:   dbRest.UpdatedAt.Time,
		},
		YandexApiSlug: dbRest.YandexApiSlug}
}

func GoogleUUIDToPgtype(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func StringToPgtypeText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: true}
}

func TimeToPgtypeTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{Time: t, Valid: true}
}
