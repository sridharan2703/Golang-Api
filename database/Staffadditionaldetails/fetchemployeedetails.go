// Package databasesad contains structs and queries for Staff Additional details API.
// --- Creator's Info ---
// Creator: Vaishnavi
// Created On: 04-11-2025
// Last Modified By:
// Last Modified Date:
// This package handles database operations for employee E-File details
package databasesad

import (
	credentials "Hrmodule/dbconfig"
	modelssad "Hrmodule/models/Staffadditionaldetails"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// GetEmployeeEFile fetches complete employee details based on employee ID and category
// This is a generic function that works for all categories
func GetEmployeeEFile(employeeID, category string) ([]modelssad.CategoryResponse, error) {
	var result []modelssad.CategoryResponse

	// Get database connection string
	connectionString := credentials.Getdatabasehr()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return result, fmt.Errorf("DB open error: %v", err)
	}
	defer db.Close()

	var rows *sql.Rows

	// Execute query based on whether employee ID is provided
	if employeeID != "" {
		// Query with employee ID filter - uses complete query for all categories
		rows, err = db.Query(modelssad.GetEFileQueryByCategory(category), employeeID)
	} else {
		// Query all employees - uses complete query for all categories
		rows, err = db.Query(modelssad.GetEFileQueryAllByCategory(category))
	}

	if err != nil {
		return result, fmt.Errorf("error querying Postgres: %v", err)
	}
	defer rows.Close()

	// Retrieve and process the employee data
	records, err := modelssad.RetrieveEmployeeEFile(rows, category)
	if err != nil {
		return result, fmt.Errorf("error retrieving employee E-File data: %v", err)
	}

	result = records
	return result, nil
}
