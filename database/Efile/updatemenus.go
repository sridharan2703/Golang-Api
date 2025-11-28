// Package databaseefile contains structs and queries for UpdateModules.
//
// --- Creator's Info ---
// Creator: Rovita
// Created On: 24-11-2025
// Last Modified By:
// Last Modified Date:
// This api is to update status in category_role_map.

package databaseefile

import (
	credentials "Hrmodule/dbconfig"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

// UpdateMultipleCategoryRoleMapStatus updates status for multiple module IDs in category_role_map table
func UpdateMultipleCategoryRoleMapStatus(moduleIDs []string, roleName, status string) (int64, error) {
	var rowsAffected int64

	connectionString := credentials.Getdatabasehr()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return rowsAffected, fmt.Errorf("DB open error: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		return rowsAffected, fmt.Errorf("DB connection failed: %v", err)
	}

	// Build placeholders for the IN clause
	placeholders := make([]string, len(moduleIDs))
	args := make([]interface{}, len(moduleIDs)+2) // +2 for status and role_name

	// First argument is status
	args[0] = status

	// Add module IDs to arguments and placeholders
	for i, moduleID := range moduleIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		args[i+1] = moduleID
	}

	// Last argument is role_name
	args[len(args)-1] = roleName

	// Build the query with dynamic IN clause
	query := fmt.Sprintf(
		`UPDATE meivan.category_role_map 
		SET status = $1, updated_at = NOW() 
		WHERE module_id IN (%s) AND role_name = $%d`,
		strings.Join(placeholders, ","),
		len(args),
	)

	// Execute UPDATE query
	result, err := db.Exec(query, args...)
	if err != nil {
		return rowsAffected, fmt.Errorf("error updating Postgres: %v", err)
	}

	// Get the number of rows affected
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return rowsAffected, fmt.Errorf("error getting rows affected: %v", err)
	}

	return rowsAffected, nil
}

// UpdateCategoryRoleMapStatus - Single update (kept for backward compatibility)
func UpdateCategoryRoleMapStatus(moduleID, roleName, status string) (int64, error) {
	return UpdateMultipleCategoryRoleMapStatus([]string{moduleID}, roleName, status)
}
