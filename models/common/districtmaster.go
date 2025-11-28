// Package modelscommon contains structs and queries for District.
//
// --- Creator's Info ---
// Creator: Rovita
// Created On: 18-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all District Value Master.

package modelscommon

import (
	"database/sql"
	"fmt"
)

// DistrictMasterQuery - Query to fetch District values
const DistrictMasterQuery = `
SELECT id, countrycode, stateid, districtname, sequenceid, isactive 
FROM humanresources.districtmaster 
WHERE isactive = '1' AND countrycode = $1 AND stateid = $2
ORDER BY sequenceid, districtname
`

// DistrictMaster struct to hold District master data
type DistrictMaster struct {
	ID           int    `json:"id"`
	CountryCode  string `json:"countrycode"`
	StateID      string `json:"stateid"`
	DistrictName string `json:"districtname"`
	SequenceID   string `json:"sequenceid"`
	IsActive     string `json:"isactive"`
}

// RetrieveDistrictMaster scans District master data from query results
func RetrieveDistrictMaster(rows *sql.Rows) ([]DistrictMaster, error) {
	var list []DistrictMaster
	for rows.Next() {
		var rm DistrictMaster
		err := rows.Scan(
			&rm.ID,
			&rm.CountryCode,
			&rm.StateID,
			&rm.DistrictName,
			&rm.SequenceID,
			&rm.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning District master row: %v", err)
		}
		list = append(list, rm)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}
	
	return list, nil
}