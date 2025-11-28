// Package databasecommon contains structs and queries for District.
//
// --- Creator's Info ---
// Creator: Rovita
// Created On: 18-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all District Value Master.

package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelssad "Hrmodule/models/common"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// GetDistrictMaster fetches the list of District values from Postgres
func GetDistrictMaster(countryCode, stateID string) ([]modelssad.DistrictMaster, error) {
	var result []modelssad.DistrictMaster

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

	rows, err := db.Query(modelssad.DistrictMasterQuery, countryCode, stateID)
	if err != nil {
		return result, fmt.Errorf("error querying Postgres: %v", err)
	}
	defer rows.Close()

	records, err := modelssad.RetrieveDistrictMaster(rows)
	if err != nil {
		return result, fmt.Errorf("error retrieving District master data: %v", err)
	}

	result = records
	return result, nil
}