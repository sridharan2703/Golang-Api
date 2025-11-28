// Package modelscommon contains data structures and DB scan logic for Tasksummary API.
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On:21-11-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 21-11-2025
//
// Path:Task Inbox  Page
package modelscommon

import (
	"database/sql"
)

type TaskSummaryStruct struct {
	ProcessID         int64  `json:"process_id"`
	ProcessName       string `json:"process_name"`
	TaskID            string `json:"task_id"`
	EmployeeID        string `json:"employee_id"`
	EmployeeName      string `json:"employee_name"`
	Pending           string `json:"pending"`
	TaskStatusID      int64  `json:"task_status_id"`
	TaskStatusDesc    string `json:"task_status_desc"`
	Priority          int64  `json:"priority"`
	PriorityDesc      string `json:"priority_desc"`
	InitiatedBy       string `json:"initiated_by"`
	InitiatedName     string `json:"initiated_name"`
	DesignationName   string `json:"designation_name"`
	DateOfAppointment string `json:"date_of_appointment"`
	InitiatedOn       string `json:"initiated_on"`
	UpdatedOn         string `json:"updated_on"`
	Elapsed_days      int64  `json:"elapsed_days"`
	SourceTable       string `json:"source_table"`
}

var MyQueryTaskSummary = `
SELECT 
    process_id,
    process_name,
    task_id,
    employee_id,
    employee_name,
    pending,
    task_status_id,
    task_status_desc,
    priority,
    priority_desc,
    initiated_by,
    initiated_name,
	designation_name,
	date_of_appointment,
    initiated_on,
    updated_on,
	elapsed_days,
    source_table
FROM meivan.task_summary($1, $2, $3);
`

func RetrieveTaskSummary(rows *sql.Rows) ([]TaskSummaryStruct, error) {

	// ============================
	// Helper functions INSIDE
	// ============================
	getString := func(v sql.NullString) string {
		if v.Valid {
			return v.String
		}
		return ""
	}

	getInt := func(v sql.NullInt64) int64 {
		if v.Valid {
			return v.Int64
		}
		return 0
	}

	formatDate := func(t sql.NullTime) string {
		if t.Valid {
			return t.Time.Format("02-01-2006") // dd-mm-yyyy
		}
		return ""
	}
	// ============================

	var results []TaskSummaryStruct

	for rows.Next() {

		var rec TaskSummaryStruct

		var (
			processID       sql.NullInt64
			processName     sql.NullString
			taskID          sql.NullString
			employeeID      sql.NullString
			employeeName    sql.NullString
			pending         sql.NullString
			taskStatusID    sql.NullInt64
			taskStatusDesc  sql.NullString
			priority        sql.NullInt64
			priorityDesc    sql.NullString
			initiatedBy     sql.NullString
			initiatedName   sql.NullString
			designationName sql.NullString
			dateOfAppt      sql.NullTime
			initiatedOn     sql.NullTime
			updatedOn       sql.NullTime
			elapsed_days    sql.NullInt64
			sourceTable     sql.NullString
		)

		err := rows.Scan(
			&processID,
			&processName,
			&taskID,
			&employeeID,
			&employeeName,
			&pending,
			&taskStatusID,
			&taskStatusDesc,
			&priority,
			&priorityDesc,
			&initiatedBy,
			&initiatedName,
			&designationName,
			&dateOfAppt,
			&initiatedOn,
			&updatedOn,
			&elapsed_days,
			&sourceTable,
		)

		if err != nil {
			return nil, err
		}

		// Assign values without showing Valid:true
		rec.ProcessID = getInt(processID)
		rec.ProcessName = getString(processName)
		rec.TaskID = getString(taskID)
		rec.EmployeeID = getString(employeeID)
		rec.EmployeeName = getString(employeeName)
		rec.Pending = getString(pending)
		rec.TaskStatusID = getInt(taskStatusID)
		rec.TaskStatusDesc = getString(taskStatusDesc)
		rec.Priority = getInt(priority)
		rec.PriorityDesc = getString(priorityDesc)
		rec.InitiatedBy = getString(initiatedBy)
		rec.InitiatedName = getString(initiatedName)
		rec.DesignationName = getString(designationName)
		rec.DateOfAppointment = formatDate(dateOfAppt)
		rec.InitiatedOn = formatDate(initiatedOn)
		rec.UpdatedOn = formatDate(updatedOn)
		rec.Elapsed_days = getInt(elapsed_days)
		rec.SourceTable = getString(sourceTable)

		results = append(results, rec)
	}

	return results, nil
}
