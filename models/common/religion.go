
// Package modelscommon contains structs and queries for Religion.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 15-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all Religion Value Master.

package modelscommon

import (
	"database/sql"
	"fmt"
)

// ReligionMasterQuery - Query to fetch Religion values
const ReligionMasterQuery = `
SELECT id, name, description, isactive 
FROM humanresources.religion 
WHERE isactive='1'
ORDER BY id
`

// ReligionMaster struct to hold Religion master data
type ReligionMaster struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    string `json:"isactive"`
}

// RetrieveReligionMaster scans Religion master data from query results
func RetrieveReligionMaster(rows *sql.Rows) ([]ReligionMaster, error) {
	var list []ReligionMaster
	for rows.Next() {
		var rm ReligionMaster
		err := rows.Scan(
			&rm.ID,
			&rm.Name,
			&rm.Description,
			&rm.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning Religion master row: %v", err)
		}
		list = append(list, rm)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}
	
	return list, nil
}