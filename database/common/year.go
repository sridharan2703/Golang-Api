// Package databasecommon contains structs and queries for YearMaster.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 17-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all YearMaster Value Master.

package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelssad "Hrmodule/models/common"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// GetYearMaster fetches the list of YearMaster values from Postgres
func GetYearMaster() ([]modelssad.YearMaster, error) {
	var result []modelssad.YearMaster

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

	rows, err := db.Query(modelssad.YearMasterQuery)
	if err != nil {
		return result, fmt.Errorf("error querying Postgres: %v", err)
	}
	defer rows.Close()

	records, err := modelssad.RetrieveYearMaster(rows)
	if err != nil {
		return result, fmt.Errorf("error retrieving YearMaster data: %v", err)
	}

	result = records
	return result, nil
}