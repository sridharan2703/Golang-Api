// Package databasesad contains structs and queries for Staff Additional details API.
//
// --- Creator's Info ---
// Creator: Vaishnavi
// Created On: 04-11-2025
// Last Modified By:  Rovita
// Last Modified Date: 12-1-2025
// This api is to feth the all the education details
package databasesad

import (
	credentials "Hrmodule/dbconfig"
	modelssad "Hrmodule/models/Staffadditionaldetails"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// FetchEmployeeEducationDetails connects to DB and retrieves employee education details
func FetchEmployeeEducationDetails(employeeID string) (interface{}, error) {
	connectionString := credentials.Getdatabasehr()

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("DB open error: %v", err)
	}
	defer db.Close()

	// Verify DB connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("DB connection failed: %v", err)
	}

	// Use the generic retriever function
	data, err := modelssad.GenericEducationDetailsRetriever(db, employeeID)
	if err != nil {
		return nil, fmt.Errorf("retrieving Education details failed: %v", err)
	}

	return data, nil
}
