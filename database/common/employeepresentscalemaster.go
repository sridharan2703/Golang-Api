// Package databasecommon contains structs and queries for EmployeePresentScaleMaster.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 17-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all EmployeePresentScaleMaster Value Master.

package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelscommon "Hrmodule/models/common"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// GetEmployeePresentScaleMaster fetches the list of EmployeePresentScaleMaster values from Postgres with gradegroup filter
func GetEmployeePresentScaleMaster(gradeGroup string) ([]modelscommon.EmployeePresentScaleMaster, error) {
	var result []modelscommon.EmployeePresentScaleMaster

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

	var rows *sql.Rows
	if gradeGroup != "" {
		// Query with gradegroup filter
		rows, err = db.Query(modelscommon.EmployeePresentScaleMasterQueryWithFilter, gradeGroup)
	} else {
		// Query without filter (get all records)
		rows, err = db.Query(modelscommon.EmployeePresentScaleMasterQuery)
	}
	
	if err != nil {
		return result, fmt.Errorf("error querying Postgres: %v", err)
	}
	defer rows.Close()

	records, err := modelscommon.RetrieveEmployeePresentScaleMaster(rows)
	if err != nil {
		return result, fmt.Errorf("error retrieving EmployeePresentScaleMaster data: %v", err)
	}

	result = records
	return result, nil
}