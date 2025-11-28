// Package modelssad contains structs and queries for Staff Additional details API.
//
// --- Creator's Info ---
// Creator: Vaishnavi
// Created On: 04-11-2025
// Last Modified By:  Rovita
// Last Modified Date: 12-1-2025
// This api is to feth the all the document details
package modelssad

import (
	"database/sql"
	"fmt"
)

// DocumentDetail represents a single document record
type DocumentDetail struct {
	EmployeeID   string `json:"employeeid"`
	DocumentName string `json:"document_name"`
}

// GenericDocumentDetailsRetriever retrieves and processes document details
func GenericDocumentDetailsRetriever(db *sql.DB, employeeID string) (interface{}, error) {
	rows, err := db.Query(`SELECT * FROM humanresources.get_employee_document_details($1)`, employeeID)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var documents []DocumentDetail

	for rows.Next() {
		var doc DocumentDetail
		if err := rows.Scan(&doc.EmployeeID, &doc.DocumentName); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		documents = append(documents, doc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return documents, nil
}
