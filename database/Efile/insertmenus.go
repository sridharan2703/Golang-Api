// Package databaseefile contains structs and queries for InsertModules.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 24-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to insert data into category_role_map.

package databaseefile

import (
	credentials "Hrmodule/dbconfig"
	"Hrmodule/models/Efile"
	"database/sql"
	"fmt"
	"strings"
	_ "github.com/lib/pq"
)

// InsertMultipleCategoryRoleMap inserts multiple records into category_role_map table with duplicate check
func InsertMultipleCategoryRoleMap(moduleIDs []string, roleName, status string) (int64, error) {
	var totalRowsAffected int64

	connectionString := credentials.Getdatabasehr()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return totalRowsAffected, fmt.Errorf("DB open error: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		return totalRowsAffected, fmt.Errorf("DB connection failed: %v", err)
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return totalRowsAffected, fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// Check for existing records to avoid duplicates
	existingModules := make(map[string]bool)
	
	// Build placeholders for the IN clause
	placeholders := make([]string, len(moduleIDs))
	args := make([]interface{}, len(moduleIDs))
	for i, moduleID := range moduleIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = moduleID
	}

	checkQuery := fmt.Sprintf(
		`SELECT module_id FROM meivan.category_role_map 
		 WHERE module_id IN (%s) AND role_name = $%d`,
		strings.Join(placeholders, ","),
		len(moduleIDs)+1,
	)

	// Add roleName as the last argument
	args = append(args, roleName)

	rows, err := tx.Query(checkQuery, args...)
	if err != nil {
		return totalRowsAffected, fmt.Errorf("error checking existing data: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var existingModuleID string
		if err := rows.Scan(&existingModuleID); err != nil {
			return totalRowsAffected, fmt.Errorf("error scanning existing data: %v", err)
		}
		existingModules[existingModuleID] = true
	}

	if err := rows.Err(); err != nil {
		return totalRowsAffected, fmt.Errorf("error during existing data check: %v", err)
	}

	// Prepare insert statement
	stmt, err := tx.Prepare(modelsefile.InsertCategoryRoleMapQuery)
	if err != nil {
		return totalRowsAffected, fmt.Errorf("error preparing insert statement: %v", err)
	}
	defer stmt.Close()

	// Insert only non-existing records
	for _, moduleID := range moduleIDs {
		if existingModules[moduleID] {
			continue // Skip existing records
		}

		result, err := stmt.Exec(moduleID, roleName, status)
		if err != nil {
			return totalRowsAffected, fmt.Errorf("error inserting module %s: %v", moduleID, err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return totalRowsAffected, fmt.Errorf("error getting rows affected for module %s: %v", moduleID, err)
		}

		totalRowsAffected += rowsAffected
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return totalRowsAffected, fmt.Errorf("error committing transaction: %v", err)
	}

	// If no rows were inserted because all were duplicates
	if totalRowsAffected == 0 && len(existingModules) > 0 {
		return totalRowsAffected, fmt.Errorf("all modules already exist for this role")
	}

	return totalRowsAffected, nil
}

// InsertCategoryRoleMap - Single insert (kept for backward compatibility)
func InsertCategoryRoleMap(moduleID, roleName, status string) (int64, error) {
	return InsertMultipleCategoryRoleMap([]string{moduleID}, roleName, status)
}