// Package databaseofficeorder handles database operations for OfficeOrdertaskvisitdetails.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 15-09-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 15-11-2025
package databaseofficeorder

import (
	credentials "Hrmodule/dbconfig"
	modelsofficeorder "Hrmodule/models/officeorder"
	"database/sql"
	"fmt"
)

func GetTaskVisitDetailsFromDB(decryptedData map[string]interface{}) ([]modelsofficeorder.TaskDetails, int, error) {
	connectionString := credentials.Getdatabasemeivan()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, 0, fmt.Errorf("DB open error: %v", err)
	}
	defer db.Close()

	// Extract and validate parameters
	status, _ := decryptedData["status"].(string)
	emp, _ := decryptedData["employeeid"].(string)
	cp, _ := decryptedData["coverpageno"].(string)

	// Ensure at least one valid combination of parameters exists
	if status == "" && (emp == "" || cp == "") {
		return nil, 0, fmt.Errorf("missing required parameters: either 'status' or both 'employeeid' and 'coverpageno' must be provided")
	}

	var rows *sql.Rows
	var isCompletedQuery bool

	// Special Case: status = "complete" - use dedicated completed query
	if status == "complete" {
		rows, err = db.Query(modelsofficeorder.MyQueryTaskDetailsBycompleted)
		if err != nil {
			return nil, 0, fmt.Errorf("query error (completed tasks): %v", err)
		}
		defer rows.Close()
		isCompletedQuery = true
	}

	// Case 1: status + employee + coverpageno (non-complete status)
	if rows == nil && status != "" && emp != "" && cp != "" {
		rows, err = db.Query(modelsofficeorder.MyQueryTaskDetailsByEmployeeCoverPage, status, emp, cp)
		if err != nil {
			return nil, 0, fmt.Errorf("query error (by emp/cover/status): %v", err)
		}
		defer rows.Close()
	}

	// Case 2: status only (non-complete status)
	if rows == nil && status != "" {
		rows, err = db.Query(modelsofficeorder.MyQueryTaskDetails, status)
		if err != nil {
			return nil, 0, fmt.Errorf("query error (by status only): %v", err)
		}
		defer rows.Close()
	}

	// Case 3: employee + coverpageno (no status)
	if rows == nil && emp != "" && cp != "" {
		rows, err = db.Query(modelsofficeorder.MyQueryTaskDetailsByEmployeeCoverPageOnly, emp, cp)
		if err != nil {
			return nil, 0, fmt.Errorf("query error (by emp/cover only): %v", err)
		}
		defer rows.Close()
	}

	if rows == nil {
		return nil, 0, fmt.Errorf("no valid query condition met")
	}

	// Retrieve data using appropriate function
	var data []modelsofficeorder.TaskDetails
	if isCompletedQuery {
		data, err = modelsofficeorder.RetrieveCompletedTaskDetails(rows)
	} else {
		data, err = modelsofficeorder.RetrieveTaskDetails(rows)
	}

	if err != nil {
		return nil, 0, fmt.Errorf("data retrieval error: %v", err)
	}

	return data, len(data), nil
}
