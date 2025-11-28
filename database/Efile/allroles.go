// Package databaseefile contains structs and queries for ALLRoles.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 24-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all ALLRoles Value Master.

package databaseefile

import (
	credentials "Hrmodule/dbconfig"
	"Hrmodule/models/Efile"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// ALLRolesMaster struct to hold ALLRoles master data
type ALLRolesMaster struct {
	Rolename string `json:"rolename"`
	// Status   string `json:"status"`
}

// GetALLRolesMaster fetches the list of ALLRoles values from Postgres
func GetALLRolesMaster() ([]ALLRolesMaster, error) {
	var result []ALLRolesMaster

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

	rows, err := db.Query(modelsefile.ALLRolesMasterQuery)
	if err != nil {
		return result, fmt.Errorf("error querying Postgres: %v", err)
	}
	defer rows.Close()

	records, err := RetrieveALLRolesMaster(rows)
	if err != nil {
		return result, fmt.Errorf("error retrieving ALLRoles master data: %v", err)
	}

	result = records
	return result, nil
}

// RetrieveALLRolesMaster scans ALLRoles master data from query results
func RetrieveALLRolesMaster(rows *sql.Rows) ([]ALLRolesMaster, error) {
	var list []ALLRolesMaster
	for rows.Next() {
		var rm ALLRolesMaster
		err := rows.Scan(
			&rm.Rolename,
			// &rm.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning ALLRoles master row: %v", err)
		}
		list = append(list, rm)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}
	
	return list, nil
}