package fetcher

import "fmt"

type DBActionResult struct {
	format       string
	rowsAffected int64
}

func makeDBActionResult(format string, rowsAffected int64) DBActionResult {
	return DBActionResult{format: format + "\n", rowsAffected: rowsAffected}
}

func (dbAR *DBActionResult) print() {
	fmt.Printf(dbAR.format, dbAR.rowsAffected)
}
