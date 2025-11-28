// Package modelsofficeorder contains structs and queries for OfficeOrder_visitdetails API.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 15-09-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 15-09-2025
package modelsofficeorder

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// ----------------------
// Queries
// ----------------------

// Query to fetch all visit details from postgresql
var MyQueryVisitDetails = `SELECT * FROM wf_integration.WF_officeorder WHERE status in (0,1)`

// Query to fetch visit details by EmployeeID + CoverPageNo
var MyQueryVisitDetailsemployeecoverpage = `SELECT * FROM wf_integration.WF_officeorder WHERE status in (0,1)AND EmployeeID = $1 AND CoverPageNo = $2`

// ----------------------
// Structs
// ----------------------

type LeaveDetails struct {
	LeaveTypeIDValue   string  `json:"leavetype"`
	LProposedVisitFrom string  `json:"startdate"`
	LProposedVisitTo   string  `json:"enddate"`
	DurationOfVisit    float64 `json:"duration"`
}

type VisitDetails struct {
	EmployeeID                 *string        `json:"employeeid"`
	FacultyName                *string        `json:"facultyname"`
	Department                 *string        `json:"department"`
	Designation                *string        `json:"designation"`
	FacultyDetails             *string        `json:"facultydetails"`
	VisitFrom                  *string        `json:"visitfrom"`
	VisitTo                    *string        `json:"visitto"`
	NatureOfParticipation      *string        `json:"natureofparticipation"`
	NatureOfParticipationValue *string        `json:"natureofparticipation_value"`
	Claimtype                  *string        `json:"claimtype"`
	Country                    *string        `json:"country"`
	CoverPageNo                *string        `json:"coverpageno"`
	InitiatedOn                *string        `json:"initiatedon"`
	CityTown                   *string        `json:"citytown"`
	LeaveDetails               []LeaveDetails `json:"leavedetails"`
}

// ----------------------
// Retrieval Functions
// ----------------------

// RetrieveVisitDetails retrieves all visit details from the rows
func RetrieveVisitDetails(rows *sql.Rows) ([]VisitDetails, error) {
	var list []VisitDetails
	for rows.Next() {
		var v VisitDetails
		var leaveDetailsJSON sql.NullString // changed here
		var status string
		var updatedOn string
		var updatedBy string

		err := rows.Scan(
			&v.EmployeeID,
			&v.FacultyName,
			&v.Department,
			&v.Designation,
			&v.FacultyDetails,
			&v.VisitFrom,
			&v.VisitTo,
			&v.Country,
			&v.CityTown,
			&v.CoverPageNo,
			&v.NatureOfParticipation,
			&v.NatureOfParticipationValue,
			&v.InitiatedOn,
			&v.Claimtype,
			&leaveDetailsJSON,
			&status,
			&updatedOn,
			&updatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		if leaveDetailsJSON.Valid && leaveDetailsJSON.String != "" {
			err = json.Unmarshal([]byte(leaveDetailsJSON.String), &v.LeaveDetails)
			if err != nil {
				return nil, fmt.Errorf("error unmarshaling leave details JSON: %v", err)
			}
		} else {
			v.LeaveDetails = nil // or empty slice: []LeaveDetails{}
		}

		list = append(list, v)
	}
	return list, nil
}

// RetrieveVisitDetailsByCover retrieves visit details by EmployeeID and CoverPageNo from the rows
func RetrieveVisitDetailsByCover(rows *sql.Rows) ([]VisitDetails, error) {
	var list []VisitDetails
	for rows.Next() {
		var v VisitDetails
		var leaveDetailsJSON sql.NullString // changed here too
		var status string
		var updatedOn string
		var updatedBy string

		err := rows.Scan(
			&v.EmployeeID,
			&v.FacultyName,
			&v.Department,
			&v.Designation,
			&v.FacultyDetails,
			&v.VisitFrom,
			&v.VisitTo,
			&v.Country,
			&v.CityTown,
			&v.CoverPageNo,
			&v.NatureOfParticipation,
			&v.NatureOfParticipationValue,
			&v.InitiatedOn,
			&v.Claimtype,
			&leaveDetailsJSON,
			&status,
			&updatedOn,
			&updatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		if leaveDetailsJSON.Valid && leaveDetailsJSON.String != "" {
			err = json.Unmarshal([]byte(leaveDetailsJSON.String), &v.LeaveDetails)
			if err != nil {
				return nil, fmt.Errorf("error unmarshaling leave details JSON: %v", err)
			}
		} else {
			v.LeaveDetails = nil // or []LeaveDetails{}
		}

		list = append(list, v)
	}
	return list, nil
}
