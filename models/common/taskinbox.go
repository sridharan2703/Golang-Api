// Package modelscommon contains data structures and DB scan logic for InboxTasksRole API.
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
//
// Path:Task Inbox  Page
package modelscommon

import (
	"database/sql"
	"fmt"
)

const MyQueryInboxTasksRole = `
SELECT *
FROM getinboxtasks_role($1, $2)
`

// InboxTasksRole defines the structure for getinboxtasks_role output
type InboxTasksRole struct {
	order_No       *string `json:"order_no"`
	TaskID         *string `json:"taskid"`
	EmployeeID     *string `json:"employeeid"`
	Name           *string `json:"Name"`
	Task_updatedon *string `json:"task_updatedon"`
	UpdatedBy      *string `json:"updatedby"`
	ActivitySeqNo  *int    `json:"activityseqno"`
	Remarks        *string `json:"remarks"`
	ProcessName    *string `json:"processname"`
	ProcessKey     *string `json:"processkeyword"`
	Path           *string `json:"path"`
	Component      *string `json:"component"`
	CoverPageNo    *string `json:"coverpageno"`
	ReferenceNo    *string `json:"referenceno"`
	ProcessID      *int    `json:"processid"`
	Badge          *string `json:"badge"`
	Priority       *string `json:"priority"`
	Starred        *string `json:"starred"`
}

// RetrieveInboxTasksRole scans rows into []InboxTasksRole
func RetrieveInboxTasksRole(rows *sql.Rows) ([]InboxTasksRole, error) {
	var result []InboxTasksRole

	for rows.Next() {
		var t InboxTasksRole
		err := rows.Scan(
			&t.order_No,
			&t.TaskID,
			&t.EmployeeID,
			&t.Name,
			&t.Task_updatedon,
			&t.UpdatedBy,
			&t.ActivitySeqNo,
			&t.Remarks,
			&t.ProcessName,
			&t.ProcessKey,
			&t.Path,
			&t.Component,
			&t.CoverPageNo,
			&t.ReferenceNo,
			&t.ProcessID,
			&t.Badge,
			&t.Priority,
			&t.Starred,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		result = append(result, t)
	}
	return result, nil
}
