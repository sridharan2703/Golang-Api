// Package databaseofficeorder handles DB access for ReturnDropdown API.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 29-10-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 29-10-2025
package databaseofficeorder

import (
	credentials "Hrmodule/dbconfig"
	modelsofficeorder "Hrmodule/models/officeorder"
	"database/sql"
	"fmt"
)

type ReturnDropdownRequest struct {
	TaskID string `json:"task_id"`
}

// GetReturnDropdownFromDB fetches dropdown data from DB
func GetReturnDropdownFromDB(decryptedData map[string]interface{}) ([]modelsofficeorder.ReturnDropdownStruct, int, error) {
	connectionString := credentials.Getdatabasemeivan()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, 0, fmt.Errorf("DB open error: %v", err)

	}
	defer db.Close()

	TaskID, _ := decryptedData["task_id"].(string)

	rows, err := db.Query(modelsofficeorder.MyQueryReturnDropdown, TaskID)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	data, err := modelsofficeorder.RetrieveReturnDropdown(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("retrieving result failed: %v", err)
	}

	return data, len(data), nil
}
