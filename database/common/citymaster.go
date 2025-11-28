// Package databasecommon contains structs and queries for City.
//
// --- Creator's Info ---
// Creator: Rovita
// Created On: 18-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all City Value Master.

package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelssad "Hrmodule/models/common"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// GetCityMaster fetches the list of City values from Postgres
func GetCityMaster(countryCode, stateID string) ([]modelssad.CityMaster, error) {
	var result []modelssad.CityMaster

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

	rows, err := db.Query(modelssad.CityMasterQuery, countryCode, stateID)
	if err != nil {
		return result, fmt.Errorf("error querying Postgres: %v", err)
	}
	defer rows.Close()

	records, err := modelssad.RetrieveCityMaster(rows)
	if err != nil {
		return result, fmt.Errorf("error retrieving City master data: %v", err)
	}

	result = records
	return result, nil
}