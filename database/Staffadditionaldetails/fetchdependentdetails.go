// Package databasesad contains structs and queries for Staff Additional details API.
//
// --- Creator's Info ---
// Creator: Vaishnavi
// Created On: 04-11-2025
// Last Modified By:  Rovita
// Last Modified Date: 12-1-2025
// This api is to feth the all the dependent details
package databasesad

import (
	credentials "Hrmodule/dbconfig"
	modelssad "Hrmodule/models/Staffadditionaldetails"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func FetchEmployeeDependentDetails(employeeID string) (interface{}, error) {
	connectionString := credentials.Getdatabasehr()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("DB open error: %v", err)
	}
	defer db.Close()

	// Use the generic retriever function
	data, err := modelssad.GenericDependentDetailsRetriever(db, employeeID)
	if err != nil {
		return nil, fmt.Errorf("retrieving dependent details failed: %v", err)
	}

	return data, nil
}
