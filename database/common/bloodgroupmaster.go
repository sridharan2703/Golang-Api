// Package databasecommon contains structs and queries for BloodGroup.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 15-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all BloodGroup Value Master.

package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelssad "Hrmodule/models/common"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// GetBloodGroupMaster fetches the list of BloodGroup values from Postgres
func GetBloodGroupMaster() ([]modelssad.BloodGroupMaster, error) {
	var result []modelssad.BloodGroupMaster

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

	rows, err := db.Query(modelssad.BloodGroupMasterQuery)
	if err != nil {
		return result, fmt.Errorf("error querying Postgres: %v", err)
	}
	defer rows.Close()

	records, err := modelssad.RetrieveBloodGroupMaster(rows)
	if err != nil {
		return result, fmt.Errorf("error retrieving BloodGroup master data: %v", err)
	}

	result = records
	return result, nil
}