// Package modelscommon contains structs and queries for City.
//
// --- Creator's Info ---
// Creator: Rovita
// Created On: 18-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all City Value Master.

package modelscommon

import (
	"database/sql"
	"fmt"
)

// CityMasterQuery - Query to fetch City values
const CityMasterQuery = `
SELECT id, countrycode, stateid, cityname, sequenceid, isactive 
FROM humanresources.citymaster 
WHERE isactive = '1' AND countrycode = $1 AND stateid = $2
ORDER BY sequenceid, cityname
`

// CityMaster struct to hold City master data
type CityMaster struct {
	ID          int    `json:"id"`
	CountryCode string `json:"countrycode"`
	StateID     string `json:"stateid"`
	CityName    string `json:"cityname"`
	SequenceID  string `json:"sequenceid"`
	IsActive    string `json:"isactive"`
}

// RetrieveCityMaster scans City master data from query results
func RetrieveCityMaster(rows *sql.Rows) ([]CityMaster, error) {
	var list []CityMaster
	for rows.Next() {
		var rm CityMaster
		err := rows.Scan(
			&rm.ID,
			&rm.CountryCode,
			&rm.StateID,
			&rm.CityName,
			&rm.SequenceID,
			&rm.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning City master row: %v", err)
		}
		list = append(list, rm)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}
	
	return list, nil
}