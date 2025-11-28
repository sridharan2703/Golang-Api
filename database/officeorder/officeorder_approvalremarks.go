// Package databaseofficeorder handles database operations for OfficeOrderapproval_remarks.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 15-09-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 30-09-2025
package databaseofficeorder

import (
	credentials "Hrmodule/dbconfig"
	modelsofficeorder "Hrmodule/models/officeorder"
	"database/sql"
	"fmt"
	"net/http"
)

// GetOfficeCommentsFromDB fetches remarks using process_id + task_id
func GetOfficeCommentsFromDB(w http.ResponseWriter, r *http.Request, decryptedData map[string]interface{}) ([]modelsofficeorder.OfficeCommentStructure, int, error) {

	// 1. DB Connection
	connectionString := credentials.Getdatabasemeivan()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		http.Error(w, fmt.Sprintf("DB open error for comments: %v", err), http.StatusInternalServerError)
		return nil, 0, err
	}
	defer db.Close()

	// 2. Extract input parameters
	processID, ok1 := decryptedData["process_id"].(float64)
	taskID, ok2 := decryptedData["taskid"].(string)

	if !ok1 || int(processID) == 0 || !ok2 || taskID == "" {
		http.Error(w, "Missing required parameters process_id & taskid", http.StatusBadRequest)
		return nil, 0, fmt.Errorf("invalid input parameters")
	}

	// 3. Execute Query
	rows, err := db.Query(modelsofficeorder.QueryOfficeComments, int(processID), taskID)
	if err != nil {
		return nil, 0, fmt.Errorf("error executing getcomments query: %v", err)
	}
	defer rows.Close()

	// 4. Process Results
	data, err := modelsofficeorder.RetrieveOfficeComments(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("retrieving office comments failed: %v", err)
	}

	// 5. Return response
	return data, len(data), nil
}
