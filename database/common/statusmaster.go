// Package databasecommon handles DB calls for StatusMaster API.
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On:26-08-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 26-08-2025
package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelscommon "Hrmodule/models/common"
	"database/sql"
	"fmt"
)

// StatusMasterDatabase executes the query
func StatusMasterDatabase(decryptedData map[string]interface{}) ([]modelscommon.StatusMaster, int, error) {
	connectionString := credentials.Getdatabasemeivan()

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, 0, fmt.Errorf("DB open error: %v", err)
	}
	defer db.Close()

	// Extract order_type_id from decrypted data
	StatusName, ok := decryptedData["statusname"].(string)
	if !ok || StatusName == "" {
		return nil, 0, fmt.Errorf("missing 'statusname' in request data")
	}

	// Execute query
	rows, err := db.Query(modelscommon.MyQueryStatusMaster, StatusName)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying DB: %v", err)
	}
	defer rows.Close()

	// Map results
	data, err := modelscommon.RetrieveStatusMaster(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving data: %v", err)
	}

	return data, len(data), nil
}
