// Package modelscommon contains structs and queries for State.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 18-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all State Value Master.

package modelscommon

import (
	"database/sql"
	"fmt"
	
)

// StateMasterQuery - Query to fetch State values
const StateMasterQuery = `
SELECT id, countrycode, statename, sequenceid, isactive 
FROM humanresources.statemaster 
WHERE isactive = '1' AND countrycode = $1
ORDER BY sequenceid, statename
`

// StateMaster struct to hold State master data
type StateMaster struct {
	ID          int    `json:"id"`
	CountryCode string `json:"countrycode"`
	StateName   string `json:"statename"`
	SequenceID  string `json:"sequenceid"`
	IsActive    string `json:"isactive"`
}

// RetrieveStateMaster scans State master data from query results
func RetrieveStateMaster(rows *sql.Rows) ([]StateMaster, error) {
	var list []StateMaster
	for rows.Next() {
		var rm StateMaster
		err := rows.Scan(
			&rm.ID,
			&rm.CountryCode,
			&rm.StateName,
			&rm.SequenceID,
			&rm.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning State master row: %v", err)
		}
		list = append(list, rm)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}
	
	return list, nil
}