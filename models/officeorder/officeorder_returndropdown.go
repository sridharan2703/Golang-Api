// Package modelsofficeorder contains structs and queries for ReturnDropdown API.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 29-10-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 29-10-2025
package modelsofficeorder

import (
	"database/sql"
	"fmt"
)

// SQL Query
var MyQueryReturnDropdown = (`
SELECT DISTINCT
    m.task_id,
    n.user_id AS comment_by,
    e.firstname,
    n.user_role AS comment_role,
    n.updated_on AS comment_updated_on
FROM meivan.pcr_m m
JOIN meivan.pcr_a a ON m.task_id = a.task_id
JOIN meivan.pcr_n n ON m.task_id = n.task_id
JOIN humanresources.employeebasicinfo e ON e.employeeid = n.user_id
WHERE m.task_id = $1
ORDER BY n.updated_on DESC
`)

// Struct
type ReturnDropdownStruct struct {
	TaskID           *string `json:"task_id"`
	CommentBy        *string `json:"comment_by"`
	Firstname        *string `json:"firstname"`
	CommentRole      *string `json:"comment_role"`
	CommentUpdatedOn *string `json:"comment_updated_on"`
}

// Function to retrieve data
func RetrieveReturnDropdown(rows *sql.Rows) ([]ReturnDropdownStruct, error) {
	var list []ReturnDropdownStruct
	for rows.Next() {
		var s ReturnDropdownStruct
		var updatedOn sql.NullString

		err := rows.Scan(
			&s.TaskID,
			&s.CommentBy,
			&s.Firstname,
			&s.CommentRole,
			&updatedOn,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// ✅ Format updated_on to YYYY-MM-DD only
		if updatedOn.Valid {
			if len(updatedOn.String) >= 10 {
				formatted := updatedOn.String[:10]
				s.CommentUpdatedOn = &formatted
			} else {
				s.CommentUpdatedOn = &updatedOn.String
			}
		} else {
			s.CommentUpdatedOn = nil
		}

		list = append(list, s)
	}
	return list, nil
}
