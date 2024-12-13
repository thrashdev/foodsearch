package fetcher

import (
	"fmt"

	"github.com/thrashdev/foodsearch/internal/config"
)

type DBActionResult struct {
	records []DBActionResultRecord
}

type DBActionResultRecord struct {
	format       string
	rowsAffected int64
}

func makeDBActionResultRecord(format string, rowsAffected int64) DBActionResultRecord {
	return DBActionResultRecord{format: format + "\n", rowsAffected: rowsAffected}
}

type startupCommand func(*config.ServiceConfig) (DBActionResult, error)

func (dbAR *DBActionResult) print() {
	for _, i := range dbAR.records {
		fmt.Printf(i.format, i.rowsAffected)
	}
}
