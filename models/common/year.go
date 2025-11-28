// Package modelscommon contains structs and queries for YearMaster.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 17-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all YearMaster Value Master.

package modelscommon

import (
	"database/sql"
	"fmt"
)

// YearMasterQuery - Query to fetch YearMaster values
const YearMasterQuery = `
SELECT generate_series(
    EXTRACT(YEAR FROM CURRENT_DATE)::integer - 50,
    EXTRACT(YEAR FROM CURRENT_DATE)::integer
) AS year;
`

// YearMaster struct to hold YearMaster data
type YearMaster struct {
	Year int `json:"year"`
}

// RetrieveYearMaster scans YearMaster data from query results
func RetrieveYearMaster(rows *sql.Rows) ([]YearMaster, error) {
	var list []YearMaster
	for rows.Next() {
		var ym YearMaster
		err := rows.Scan(
			&ym.Year,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning YearMaster row: %v", err)
		}
		list = append(list, ym)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}
	
	return list, nil
}