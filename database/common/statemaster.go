// Package databasecommon contains structs and queries for State.
//
// --- Creator's Info ---
// Creator: Rovita
// Created On: 18-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all State Value Master.

package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelssad "Hrmodule/models/common"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// GetStateMaster fetches the list of State values from Postgres
func GetStateMaster(countryCode string) ([]modelssad.StateMaster, error) {
	var result []modelssad.StateMaster

	connectionString := credentials.Getdatabasehr()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return result, fmt.Errorf("DB open error: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		return result, fmt.Errorf("DB connection failed: %v", err)
	}

	rows, err := db.Query(modelssad.StateMasterQuery, countryCode)
	if err != nil {
		return result, fmt.Errorf("error querying Postgres: %v", err)
	}
	defer rows.Close()

	records, err := modelssad.RetrieveStateMaster(rows)
	if err != nil {
		return result, fmt.Errorf("error retrieving State master data: %v", err)
	}

	result = records
	return result, nil
}