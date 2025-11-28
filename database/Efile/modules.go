// Package databaseefile contains structs and queries for ALLModules.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 24-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all ALLModules Value Master.

package databaseefile

import (
	credentials "Hrmodule/dbconfig"
	"Hrmodule/models/Efile"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// ALLModulesMaster struct to hold ALLModules master data
type ALLModulesMaster struct {
	ModuleName string `json:"module_name"`
}

// GetALLModulesMaster fetches the list of ALLModules values from Postgres based on role_name
func GetALLModulesMaster(roleName string) ([]ALLModulesMaster, error) {
	var result []ALLModulesMaster

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

	// Execute query with role_name parameter
	rows, err := db.Query(modelsefile.ALLModulesMasterQuery, roleName)
	if err != nil {
		return result, fmt.Errorf("error querying Postgres: %v", err)
	}
	defer rows.Close()

	records, err := RetrieveALLModulesMaster(rows)
	if err != nil {
		return result, fmt.Errorf("error retrieving ALLModules master data: %v", err)
	}

	result = records
	return result, nil
}

// RetrieveALLModulesMaster scans ALLModules master data from query results
func RetrieveALLModulesMaster(rows *sql.Rows) ([]ALLModulesMaster, error) {
	var list []ALLModulesMaster
	for rows.Next() {
		var rm ALLModulesMaster
		err := rows.Scan(
			&rm.ModuleName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning ALLModules master row: %v", err)
		}
		list = append(list, rm)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}
	
	return list, nil
}