// Package modelsofficeorder contains structs and queries for Taskdetails in tasksummary  API.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 21-11-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 21-11-2025
package modelsofficeorder

import (
	"database/sql"
)

type PCRTaskDetailsStruct struct {
	TaskID          string `json:"task_id"`
	ProcessID       int64  `json:"process_id"`
	CoverPageNo     string `json:"cover_page_no"`
	EmployeeID      string `json:"employee_id"`
	EmployeeName    string `json:"employee_name"`
	Department      string `json:"department"`
	Designation     string `json:"designation"`
	VisitFrom       string `json:"visit_from"`
	VisitTo         string `json:"visit_to"`
	Duration        int64  `json:"duration"`
	NatureOfVisit   string `json:"nature_of_visit"`
	ClaimType       string `json:"claim_type"`
	CityTown        string `json:"city_town"`
	Country         string `json:"country"`
	OrderNo         string `json:"order_no"`
	OrderDate       string `json:"order_date"`
	OriginalOrderNo string `json:"original_order_no"`
	OrderType       string `json:"order_type"`
}

var MyQueryPCRTaskDetails = `
SELECT 
    task_id,
    process_id,
    cover_page_no,
    employee_id,
    employee_name,
    department,
    designation,
    visit_from,
    visit_to,
    duration,
    nature_of_visit,
    claim_type,
    city_town,
    country,
    order_no,
    order_date,
    original_order_no,
    order_type
FROM meivan.pcr_task_details_tasksummary($1, $2);
`

func RetrievePCRTaskDetails(rows *sql.Rows) ([]PCRTaskDetailsStruct, error) {

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

	var results []PCRTaskDetailsStruct

	for rows.Next() {

		var rec PCRTaskDetailsStruct

		var (
			taskID          sql.NullString
			processID       sql.NullInt64
			coverPageNo     sql.NullString
			employeeID      sql.NullString
			employeeName    sql.NullString
			department      sql.NullString
			designation     sql.NullString
			visitFrom       sql.NullTime
			visitTo         sql.NullTime
			duration        sql.NullInt64
			natureOfVisit   sql.NullString
			claimType       sql.NullString
			cityTown        sql.NullString
			country         sql.NullString
			orderNo         sql.NullString
			orderDate       sql.NullTime
			originalOrderNo sql.NullString
			orderType       sql.NullString
		)

		err := rows.Scan(
			&taskID,
			&processID,
			&coverPageNo,
			&employeeID,
			&employeeName,
			&department,
			&designation,
			&visitFrom,
			&visitTo,
			&duration,
			&natureOfVisit,
			&claimType,
			&cityTown,
			&country,
			&orderNo,
			&orderDate,
			&originalOrderNo,
			&orderType,
		)

		if err != nil {
			return nil, err
		}

		rec.TaskID = getString(taskID)
		rec.ProcessID = getInt(processID)
		rec.CoverPageNo = getString(coverPageNo)
		rec.EmployeeID = getString(employeeID)
		rec.EmployeeName = getString(employeeName)
		rec.Department = getString(department)
		rec.Designation = getString(designation)
		rec.VisitFrom = formatDate(visitFrom)
		rec.VisitTo = formatDate(visitTo)
		rec.Duration = getInt(duration)
		rec.NatureOfVisit = getString(natureOfVisit)
		rec.ClaimType = getString(claimType)
		rec.CityTown = getString(cityTown)
		rec.Country = getString(country)
		rec.OrderNo = getString(orderNo)
		rec.OrderDate = formatDate(orderDate)
		rec.OriginalOrderNo = getString(originalOrderNo)
		rec.OrderType = getString(orderType)

		results = append(results, rec)
	}

	return results, nil
}
