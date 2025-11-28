// Package databaseofficeorder contains structs and queries for Taskdetails in tasksummary  API.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 21-11-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 21-11-2025
package databaseofficeorder

import (
	credentials "Hrmodule/dbconfig"
	modelsofficeorder "Hrmodule/models/officeorder"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func GetPCRTaskDetailsFromDB(decryptedData map[string]interface{}) ([]modelsofficeorder.PCRTaskDetailsStruct, int, error) {

	connectionString := credentials.Getdatabasemeivan()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, 0, fmt.Errorf("DB open error: %v", err)
	}
	defer db.Close()

	processID := fmt.Sprintf("%v", decryptedData["process_id"])
	taskID := fmt.Sprintf("%v", decryptedData["task_id"])

	if processID == "" || taskID == "" {
		return nil, 0, fmt.Errorf("process_id and task_id are required")
	}

	rows, err := db.Query(modelsofficeorder.MyQueryPCRTaskDetails, processID, taskID)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	result, err := modelsofficeorder.RetrievePCRTaskDetails(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("retrieving result failed: %v", err)
	}

	return result, len(result), nil
}
