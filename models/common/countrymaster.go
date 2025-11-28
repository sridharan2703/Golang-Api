// Package modelscommon contains structs and queries for Country.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 18-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all Country Value Master.

package modelscommon

import (
	"database/sql"
	"fmt"
)

// CountryMasterQuery - Query to fetch Country values
const CountryMasterQuery = `
SELECT id, countrycode, countryname, sequenceid, isactive 
FROM humanresources.countrymaster 
WHERE isactive = '1'
ORDER BY sequenceid, countryname
`

// CountryMaster struct to hold Country master data
type CountryMaster struct {
	ID          int    `json:"id"`
	CountryCode string `json:"countrycode"`
	CountryName string `json:"countryname"`
	SequenceID  string `json:"sequenceid"`
	IsActive    string `json:"isactive"`
}

// RetrieveCountryMaster scans Country master data from query results
func RetrieveCountryMaster(rows *sql.Rows) ([]CountryMaster, error) {
	var list []CountryMaster
	for rows.Next() {
		var rm CountryMaster
		err := rows.Scan(
			&rm.ID,
			&rm.CountryCode,
			&rm.CountryName,
			&rm.SequenceID,
			&rm.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning Country master row: %v", err)
		}
		list = append(list, rm)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}
	
	return list, nil
}