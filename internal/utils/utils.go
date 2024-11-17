package utils

import (
	"log"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

func FloatToNumeric(number float64) (value pgtype.Numeric) {
	parsed := strconv.FormatFloat(number, 'f', -1, 64)
	if err := value.Scan(parsed); err != nil {
		log.Fatal("Error scanning numeric: %v", err)
	}
	return value
}
