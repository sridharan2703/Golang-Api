// package modelsofficeorder handles DB access for Approval page
//
// --- Creator's Info ---
// Creator: Ramya
//
// Created On: 09-10-2025
// Last Modified By:
// Last Modified Date:
package modelsofficeorder

import (
	"database/sql"
	"fmt"
	"time"
)

// MyBaseQueryOfficeOrderMaster is the base SQL query to fetch all columns.
// The WHERE clause will be dynamically constructed.
const MyBaseQueryOfficeOrderMaster = `
	SELECT *
	FROM meivan.officeorder_master
`

// OfficeOrderMasterStructure (remains unchanged)
type OfficeOrderMasterStructure struct {
	OfficeOrderID         *string    `json:"OfficeOrderID"`
	CoverPageNo           *string    `json:"CoverPageNo"`
	EmployeeID            *string    `json:"EmployeeID"`
	FacultyName           *string    `json:"FacultyName"`
	Department            *string    `json:"Department"`
	Designation           *string    `json:"Designation"`
	VisitFrom             *time.Time `json:"VisitFrom"`
	VisitTo               *time.Time `json:"VisitTo"`
	Duration              *int       `json:"Duration"`
	LeaveType             *string    `json:"LeaveType"`
	NatureOfParticipation *string    `json:"NatureOfParticipation"`
	ClaimType             *string    `json:"ClaimType"`
	Country               *string    `json:"Country"`
	CityTown              *string    `json:"CityTown"`
	SigningAuthority      *string    `json:"SigningAuthority"`
	ReceipientTo          *string    `json:"ReceipientTo"`
	AssignTo              *string    `json:"AssignTo"`
	AssignedRole          *string    `json:"AssignedRole"`
	TaskStatusID          *int       `json:"TaskStatusID"`
	ActivitySeqNo         *int       `json:"ActivitySeqNo"`
	IsTaskReturn          *bool      `json:"IsTaskReturn"`
	IsTaskApproved        *bool      `json:"IsTaskApproved"`
	InitiatedBy           *string    `json:"InitiatedBy"`
	InitiatedOn           *time.Time `json:"InitiatedOn"`
	UpdatedBy             *string    `json:"UpdatedBy"`
	UpdatedOn             *time.Time `json:"UpdatedOn"`
	Template              *string    `json:"Template"`
	Body                  *string    `json:"Body"`
	Header                *string    `json:"Header"`
	Footer                *string    `json:"Footer"`
	ReferenceNumber       *string    `json:"ReferenceNumber"`
	Subject               *string    `json:"Subject"`
	Ref                   *string    `json:"Ref"`
	Date                  *string    `json:"Date"`
	Remarks               *string    `json:"Remarks"`
}

// RetrieveOfficeOrderMaster (remains unchanged)
func RetrieveOfficeOrderMaster(rows *sql.Rows) ([]OfficeOrderMasterStructure, error) {
	var officeOrderList []OfficeOrderMasterStructure

	for rows.Next() {
		var order OfficeOrderMasterStructure

		err := rows.Scan(
			&order.OfficeOrderID, &order.CoverPageNo, &order.EmployeeID, &order.FacultyName,
			&order.Department, &order.Designation, &order.VisitFrom, &order.VisitTo,
			&order.Duration, &order.LeaveType, &order.NatureOfParticipation, &order.ClaimType,
			&order.Country, &order.CityTown, &order.SigningAuthority, &order.ReceipientTo,
			&order.AssignTo, &order.AssignedRole, &order.TaskStatusID, &order.ActivitySeqNo,
			&order.IsTaskReturn, &order.IsTaskApproved, &order.InitiatedBy, &order.InitiatedOn,
			&order.UpdatedBy, &order.UpdatedOn, &order.Template, &order.Body,
			&order.Header, &order.Footer, &order.ReferenceNumber, &order.Subject,
			&order.Ref, &order.Date, &order.Remarks,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning OfficeOrderMasterStructure row: %v", err)
		}
		officeOrderList = append(officeOrderList, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %v", err)
	}
	return officeOrderList, nil
}

func GetOfficeOrderMasters(db *sql.DB, taskStatusID int, employeeID, coverPageNo string) ([]OfficeOrderMasterStructure, error) {
	// Start the query with the mandatory task status condition
	// The taskStatusID is the first parameter ($1)
	query := MyBaseQueryOfficeOrderMaster + " WHERE taskstatusid = $1"

	var args []interface{} = []interface{}{taskStatusID}

	placeholderCount := 2

	if employeeID != "" {
		query += fmt.Sprintf(" AND employeeid = $%d", placeholderCount)
		args = append(args, employeeID)
		placeholderCount++
	}

	if coverPageNo != "" {
		query += fmt.Sprintf(" AND coverpageno = $%d", placeholderCount)
		args = append(args, coverPageNo)
		placeholderCount++
	}

	// Execute the query
	rows, err := db.Query(query, args...)
	if err != nil {

		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	// Use the existing retrieval function
	return RetrieveOfficeOrderMaster(rows)
}
