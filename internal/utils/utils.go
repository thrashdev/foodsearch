package utils

import (
	"log"
	"strconv"
	"strings"
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

func GlovoRestDBtoModel(dbRest database.GlovoRestaurant) models.GlovoRestaurant {
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

func YandexRestDBtoModel(dbRest database.YandexRestaurant) models.YandexRestaurant {
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

func YandexRestModelToDB(rest models.YandexRestaurant) database.BatchCreateYandexRestaurantsParams {
	addr := pgtype.Text{String: "", Valid: false}
	if rest.Address != nil {
		addr.String = *rest.Address
		addr.Valid = true
	}

	deliveryFee := pgtype.Numeric{Valid: false}
	if rest.DeliveryFee != nil {
		deliveryFee = FloatToNumeric(*rest.DeliveryFee)
	}

	phoneNumber := pgtype.Text{String: "", Valid: false}
	if rest.Address != nil {
		phoneNumber.String = *rest.PhoneNumber
		phoneNumber.Valid = true
	}
	arg := database.BatchCreateYandexRestaurantsParams{
		ID:            pgtype.UUID{Bytes: rest.ID, Valid: true},
		Name:          strings.TrimSpace(rest.Name),
		Address:       addr,
		DeliveryFee:   deliveryFee,
		PhoneNumber:   phoneNumber,
		YandexApiSlug: rest.YandexApiSlug,
		CreatedAt:     pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
		UpdatedAt:     pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
	}

	return arg
}

func YandexDishModelToDB(d models.YandexDish) database.BatchCreateYandexDishesParams {
	arg := database.BatchCreateYandexDishesParams{
		ID:                 GoogleUUIDToPgtype(d.ID),
		Name:               strings.TrimSpace(d.Name),
		Price:              FloatToNumeric(d.Price),
		DiscountedPrice:    FloatToNumeric(d.DiscountedPrice),
		Description:        StringToPgtypeText(d.Description),
		YandexRestaurantID: GoogleUUIDToPgtype(d.YandexRestaurantID),
		CreatedAt:          TimeToPgtypeTimestamp(d.CreatedAt),
		UpdatedAt:          TimeToPgtypeTimestamp(d.UpdatedAt),
		YandexApiID:        int32(d.YandexApiID),
	}
	return arg
}

func GlovoRestModelToDB(rest models.GlovoRestaurant) database.BatchCreateGlovoRestaurantsParams {
	arg := database.BatchCreateGlovoRestaurantsParams{
		ID:                pgtype.UUID{Bytes: uuid.New(), Valid: true},
		Name:              strings.TrimSpace(rest.Name),
		Address:           *rest.Address,
		DeliveryFee:       FloatToNumeric(*rest.DeliveryFee),
		PhoneNumber:       pgtype.Text{String: *rest.PhoneNumber, Valid: true},
		GlovoApiStoreID:   int32(rest.GlovoApiStoreID),
		GlovoApiAddressID: int32(rest.GlovoApiAddressID),
		GlovoApiSlug:      rest.GlovoApiSlug,
		CreatedAt:         pgtype.Timestamp{Time: time.Now().UTC(), InfinityModifier: 0, Valid: true},
		UpdatedAt:         pgtype.Timestamp{Time: time.Now().UTC(), InfinityModifier: 0, Valid: true},
	}

	return arg
}

func GlovoDishModelToDB(dish models.GlovoDish) database.BatchCreateGlovoDishesParams {
	arg := database.BatchCreateGlovoDishesParams{
		ID:                pgtype.UUID{Bytes: uuid.New(), Valid: true},
		Name:              strings.TrimSpace(dish.Name),
		Description:       dish.Description,
		Price:             FloatToNumeric(dish.Price),
		DiscountedPrice:   FloatToNumeric(dish.DiscountedPrice),
		GlovoApiDishID:    int32(dish.GlovoAPIDishID),
		GlovoRestaurantID: pgtype.UUID{Bytes: dish.GlovoRestaurantID, Valid: true},
		CreatedAt:         pgtype.Timestamp{Time: time.Now().UTC(), InfinityModifier: 0, Valid: true},
		UpdatedAt:         pgtype.Timestamp{Time: time.Now().UTC(), InfinityModifier: 0, Valid: true},
	}

	return arg
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
