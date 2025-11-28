// Package modelssad contains structs and queries for Staff Additional details API.
//
// --- Creator's Info ---
// Creator: Vaishnavi
// Created On: 04-11-2025
// Last Modified By:  Rovita
// Last Modified Date: 12-1-2025
// This api is to feth the all active employees.
package modelssad

import (
	"database/sql"
	"fmt"
)

// -------------------- PostgreSQL --------------------

// Query to fetch employee IDs in ascending order
var MyQueryEmployeeDropdown = `
SELECT concat(ebi.employeeid,'-',ebi.displayname) as employeeid 
FROM humanresources.employeebasicinfo ebi
JOIN humanresources.employeeappointmentdetails ead
ON ebi.employeeid = ead.employeeid
WHERE ead.retirementdate >= NOW()
ORDER BY employeeid ASC
`

// Struct to hold a single employee ID
type EmployeeDropdown struct {
	EmployeeID string `json:"employeeid"`
}

// RetrieveEmployeeDropdown scans employee IDs from query results
func RetrieveEmployeeDropdown(rows *sql.Rows) ([]EmployeeDropdown, error) {
	var list []EmployeeDropdown
	for rows.Next() {
		var e EmployeeDropdown
		err := rows.Scan(&e.EmployeeID)
		if err != nil {
			return nil, fmt.Errorf("error scanning employee ID row: %v", err)
		}
		list = append(list, e)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}

	return list, nil
}
