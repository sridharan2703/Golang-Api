// Package databaseefile contains structs and queries for ALLMenus.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 24-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all ALLMenus Value Master.

package databaseefile

import (
	credentials "Hrmodule/dbconfig"
	"Hrmodule/models/Efile"
	"database/sql"
	"fmt"
	// "strings"
	_ "github.com/lib/pq"
)

// ALLMenusMaster struct to hold ALLMenus master data
type ALLMenusMaster struct {
	ModuleName   string  `json:"module_name"`
	ModuleStatus string  `json:"module_status"`
	ModuleID     *string `json:"module_id"`
}

// RoleModules struct to hold modules grouped by role
type RoleModules struct {
	RoleName string          `json:"role_name"`
	Modules  []ModuleDetail `json:"modules"`
}

// ModuleDetail struct for individual module details
type ModuleDetail struct {
	ModuleID     string `json:"module_id"`
	ModuleName   string `json:"module_name"`
	ModuleStatus string `json:"module_status"`
}

// GetALLMenusMaster fetches the list of ALLMenus values from Postgres based on role_name
func GetALLMenusMaster(roleName string) ([]ALLMenusMaster, error) {
	var result []ALLMenusMaster

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
	rows, err := db.Query(modelsefile.ALLMenusMasterQuery, roleName)
	if err != nil {
		return result, fmt.Errorf("error querying Postgres: %v", err)
	}
	defer rows.Close()

	records, err := RetrieveALLMenusMaster(rows)
	if err != nil {
		return result, fmt.Errorf("error retrieving ALLMenus master data: %v", err)
	}

	result = records
	return result, nil
}

// GetALLMenusMasterGroupedByRole fetches all modules grouped by role_name
func GetALLMenusMasterGroupedByRole() ([]RoleModules, error) {
	var result []RoleModules

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

	// Execute query to get all roles and their modules
	rows, err := db.Query(modelsefile.ALLMenusMasterGroupedByRoleQuery)
	if err != nil {
		return result, fmt.Errorf("error querying Postgres: %v", err)
	}
	defer rows.Close()

	// Group modules by role
	roleModulesMap := make(map[string][]ModuleDetail)
	
	for rows.Next() {
		var roleName sql.NullString
		var moduleID, moduleName, moduleStatus string
		
		err := rows.Scan(&roleName, &moduleID, &moduleName, &moduleStatus)
		if err != nil {
			return result, fmt.Errorf("error scanning grouped data row: %v", err)
		}
		
		// Handle NULL role_name - use "Unassigned" for modules not assigned to any role
		actualRoleName := "Unassigned"
		if roleName.Valid {
			actualRoleName = roleName.String
		}
		
		moduleDetail := ModuleDetail{
			ModuleID:     moduleID,
			ModuleName:   moduleName,
			ModuleStatus: moduleStatus,
		}
		
		roleModulesMap[actualRoleName] = append(roleModulesMap[actualRoleName], moduleDetail)
	}
	
	if err := rows.Err(); err != nil {
		return result, fmt.Errorf("error during rows iteration: %v", err)
	}

	// Convert map to slice
	for roleName, modules := range roleModulesMap {
		roleModule := RoleModules{
			RoleName: roleName,
			Modules:  modules,
		}
		result = append(result, roleModule)
	}

	return result, nil
}

// RetrieveALLMenusMaster scans ALLMenus master data from query results
func RetrieveALLMenusMaster(rows *sql.Rows) ([]ALLMenusMaster, error) {
	var list []ALLMenusMaster
	for rows.Next() {
		var rm ALLMenusMaster
		err := rows.Scan(
			&rm.ModuleName,
			&rm.ModuleStatus,
			&rm.ModuleID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning ALLMenus master row: %v", err)
		}
		list = append(list, rm)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}
	
	return list, nil
}