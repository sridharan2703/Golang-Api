// Package databasecommon handles database connections and queries related to DefaultRoleName data.
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On:30-07-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 30-07-2025

package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelscommon "Hrmodule/models/common"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func DefaultRoleNamedatabase(decryptedData map[string]interface{}) ([]modelscommon.DefaultRoleNamestructure, int, error) {
	// Connection string for Postgres
	connectionString := credentials.Getdatabasemeivan()

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, 0, fmt.Errorf("DB open error: %v", err)
	}
	defer db.Close()

	// Extract order_type_id from decrypted data
	UserName, ok := decryptedData["UserName"].(string)
	if !ok || UserName == "" {
		return nil, 0, fmt.Errorf("missing 'UserName' in request data")
	}

	// Execute the query
	rows, err := db.Query(modelscommon.MyQueryDefaultRoleName, UserName)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	// Map results
	DefaultRoleNameapi, err := modelscommon.RetrieveDefaultRoleName(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving data: %v", err)
	}

	return DefaultRoleNameapi, len(DefaultRoleNameapi), nil
}
