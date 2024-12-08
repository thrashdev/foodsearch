package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/thrashdev/foodsearch/internal/database"
	"github.com/thrashdev/foodsearch/internal/models"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	resp, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Error while writing response :%w", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
	return nil
}

func RespondWithError(w http.ResponseWriter, code int, msg string) error {
	type errResponse struct {
		Error string `json:"error"`
	}

	resp := errResponse{msg}
	err := RespondWithJSON(w, code, resp)
	if err != nil {
		return err
	}
	return nil
}

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
		ID:                 PgtypeID(d.ID),
		Name:               strings.TrimSpace(d.Name),
		Price:              FloatToNumeric(d.Price),
		DiscountedPrice:    FloatToNumeric(d.DiscountedPrice),
		Description:        StringToPgtypeText(d.Description),
		YandexRestaurantID: PgtypeID(d.YandexRestaurantID),
		CreatedAt:          TimeToPgtypeTimestamp(d.CreatedAt),
		UpdatedAt:          TimeToPgtypeTimestamp(d.UpdatedAt),
		YandexApiID:        int32(d.YandexApiID),
	}
	return arg
}

func YandexDishDBtoModel(dish database.YandexDish) models.YandexDish {
	price, err := dish.Price.Float64Value()
	if err != nil {
		log.Fatal(err)
	}
	discPrice, err := dish.DiscountedPrice.Float64Value()
	if err != nil {
		log.Fatal(err)
	}
	return models.YandexDish{
		Dish: models.Dish{
			ID:              dish.ID.Bytes,
			Name:            dish.Name,
			Description:     dish.Description.String,
			Price:           price.Float64,
			DiscountedPrice: discPrice.Float64,
			CreatedAt:       dish.CreatedAt.Time,
			UpdatedAt:       dish.UpdatedAt.Time,
		},
		YandexRestaurantID: dish.YandexRestaurantID.Bytes,
		YandexApiID:        int(dish.YandexApiID),
	}
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

func GlovoDishDBtoModel(dish database.GlovoDish) models.GlovoDish {
	price, err := dish.Price.Float64Value()
	if err != nil {
		log.Fatal(err)
	}
	discPrice, err := dish.DiscountedPrice.Float64Value()
	if err != nil {
		log.Fatal(err)
	}

	return models.GlovoDish{
		Dish: models.Dish{
			ID:              dish.ID.Bytes,
			Name:            dish.Name,
			Description:     dish.Description,
			Price:           price.Float64,
			DiscountedPrice: discPrice.Float64,
			CreatedAt:       dish.CreatedAt.Time,
			UpdatedAt:       dish.UpdatedAt.Time,
		},
		GlovoAPIDishID:    int(dish.GlovoApiDishID),
		GlovoRestaurantID: dish.GlovoRestaurantID.Bytes,
	}
}

func RestaurantBindingModeltoDB(b models.RestaurantBinding) database.BatchCreateRestaurantBindingParams {
	// glovoID := pgtype.UUID{}
	// if (b.GlovoRestaurantID == uuid.UUID{}) {
	// 	glovoID.Valid = false
	// } else {
	// 	glovoID = pgtypeID(b.GlovoRestaurantID)
	// }
	//
	// yandexID := pgtype.UUID{}
	// if (b.YandexRestaurantID == uuid.UUID{}) {
	// 	yandexID.Valid = false
	// } else {
	// 	yandexID = pgtypeID(b.GlovoRestaurantID)
	// }

	arg := database.BatchCreateRestaurantBindingParams{
		ID:                 PgtypeID(b.ID),
		GlovoRestaurantID:  PgtypeID(b.GlovoRestaurantID),
		YandexRestaurantID: PgtypeID(b.YandexRestaurantID),
	}
	return arg
}

func DishBindingsModelToDB(b models.DishBinding) database.BatchCreateDishBindingsParams {
	arg := database.BatchCreateDishBindingsParams{
		ID:                  PgtypeID(b.ID),
		RestaurantBindingID: PgtypeID(b.RestaurantBindingID),
		GlovoDishID:         PgtypeID(b.GlovoDishID),
		YandexDishID:        PgtypeID(b.YandexDishID),
	}
	return arg
}

func PgtypeID(id uuid.UUID) pgtype.UUID {
	if id == uuid.Nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: id, Valid: true}
}

func StringToPgtypeText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: true}
}

func TimeToPgtypeTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{Time: t, Valid: true}
}

func SliceContains[S []K, K comparable](collection S, element K) bool {
	for _, el := range collection {
		if el == element {
			return true
		}
	}

	return false
}
