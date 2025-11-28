// Package databasesad contains structs and queries for Staff Additional details API.
//
// --- Creator's Info ---
// Creator: Vaishnavi
// Created On: 04-11-2025
// Last Modified By:  Rovita
// Last Modified Date: 12-1-2025
// This api is to feth the all the file details
package databasesad

import (
	credentials "Hrmodule/dbconfig"
	modelssad "Hrmodule/models/Staffadditionaldetails"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// GetRoleBasedModules fetches the list of modules based on role_name
func GetRoleBasedModules(roleName string) ([]modelssad.RoleBasedModule, error) {
	var modules []modelssad.RoleBasedModule

	// Validate input
	if roleName == "" {
		return modules, fmt.Errorf("role_name cannot be empty")
	}

	connectionString := credentials.Getdatabasehr()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return modules, fmt.Errorf("database connection error: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		return modules, fmt.Errorf("database ping failed: %v", err)
	}

	log.Printf("Fetching modules for role: %s", roleName)

	// Execute parameterized query
	rows, err := db.Query(modelssad.RoleBasedModuleQuery, roleName)
	if err != nil {
		return modules, fmt.Errorf("database query error: %v", err)
	}
	defer rows.Close()

	// Retrieve and process results
	modules, err = modelssad.RetrieveRoleBasedModules(rows)
	if err != nil {
		return modules, fmt.Errorf("error retrieving modules: %v", err)
	}

	log.Printf("Successfully retrieved %d modules for role: %s", len(modules), roleName)
	return modules, nil
}
