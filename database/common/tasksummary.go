// Package databasecommon handles DB calls for Tasksummary API.
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On: 21-11-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 21-11-2025
package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelscommon "Hrmodule/models/common"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func GetTaskSummaryFromDB(decryptedData map[string]interface{}) ([]modelscommon.TaskSummaryStruct, int, error) {

	connectionString := credentials.Getdatabasemeivan()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, 0, fmt.Errorf("DB open error: %v", err)
	}
	defer db.Close()

	// Only employee_id required
	employeeID := fmt.Sprintf("%v", decryptedData["employee_id"])
	if employeeID == "" {
		return nil, 0, fmt.Errorf("employee_id is required")
	}

	var taskStatus interface{}
	var priorityVal interface{}

	// NULL handling
	if decryptedData["task_status_id"] == nil {
		taskStatus = nil
	} else {
		taskStatus = decryptedData["task_status_id"]
	}

	if decryptedData["priority"] == nil {
		priorityVal = nil
	} else {
		priorityVal = decryptedData["priority"]
	}

	// Query
	rows, err := db.Query(modelscommon.MyQueryTaskSummary, employeeID, taskStatus, priorityVal)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	result, err := modelscommon.RetrieveTaskSummary(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("retrieving result failed: %v", err)
	}

	return result, len(result), nil
}
