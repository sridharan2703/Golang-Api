// Package databasecommon contains structs and queries for languagemaster.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 15-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all languagemaster.

package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelssad "Hrmodule/models/common"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// GetLanguageMaster fetches the list of language master values from Postgres
func GetLanguageMaster() ([]modelssad.LanguageMaster, error) {
	var result []modelssad.LanguageMaster

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

	rows, err := db.Query(modelssad.LanguageMasterQuery)
	if err != nil {
		return result, fmt.Errorf("error querying Postgres: %v", err)
	}
	defer rows.Close()

	records, err := modelssad.RetrieveLanguageMaster(rows)
	if err != nil {
		return result, fmt.Errorf("error retrieving language master data: %v", err)
	}

	result = records
	return result, nil
}