// Package modelscommon contains structs and queries for EmployeePresentScaleMaster.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 17-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all EmployeePresentScaleMaster Value Master.

package modelscommon

import (
	"database/sql"
	// "encoding/json"
	"fmt"
)

// EmployeePresentScaleMasterQuery - Query to fetch all EmployeePresentScaleMaster values
const EmployeePresentScaleMasterQuery = `
SELECT id, paylevel, presentscale, paybandname, paybandscale, grade, gradegroup, status, upperlimit 
FROM humanresources.employeepresentscalemaster
`

// EmployeePresentScaleMasterQueryWithFilter - Query to fetch EmployeePresentScaleMaster values filtered by gradegroup
const EmployeePresentScaleMasterQueryWithFilter = `
SELECT id, paylevel, presentscale, paybandname, paybandscale, grade, gradegroup, status, upperlimit 
FROM humanresources.employeepresentscalemaster 
WHERE gradegroup = $1
`

// EmployeePresentScaleMaster struct to hold EmployeePresentScaleMaster data
type EmployeePresentScaleMaster struct {
	ID           int     `json:"id"`
	PayLevel     string  `json:"paylevel,omitempty"`
	PresentScale string  `json:"presentscale,omitempty"`
	PayBandName  string  `json:"paybandname,omitempty"`
	PayBandScale string  `json:"paybandscale,omitempty"`
	Grade        string  `json:"grade,omitempty"`
	GradeGroup   string  `json:"gradegroup,omitempty"`
	Status       string  `json:"status,omitempty"`
	UpperLimit   float64 `json:"upperlimit,omitempty"`
}

// RetrieveEmployeePresentScaleMaster scans EmployeePresentScaleMaster data from query results
func RetrieveEmployeePresentScaleMaster(rows *sql.Rows) ([]EmployeePresentScaleMaster, error) {
	var list []EmployeePresentScaleMaster
	for rows.Next() {
		var epsm EmployeePresentScaleMaster
		
		// Use nullable types for scanning
		var (
			payLevel     sql.NullString
			presentScale sql.NullString
			payBandName  sql.NullString
			payBandScale sql.NullString
			grade        sql.NullString
			gradeGroup   sql.NullString
			status       sql.NullString
			upperLimit   sql.NullFloat64
		)
		
		err := rows.Scan(
			&epsm.ID,
			&payLevel,
			&presentScale,
			&payBandName,
			&payBandScale,
			&grade,
			&gradeGroup,
			&status,
			&upperLimit,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning EmployeePresentScaleMaster row: %v", err)
		}
		
		// Convert nullable types to regular types, handling NULL values
		if payLevel.Valid {
			epsm.PayLevel = payLevel.String
		}
		if presentScale.Valid {
			epsm.PresentScale = presentScale.String
		}
		if payBandName.Valid {
			epsm.PayBandName = payBandName.String
		}
		if payBandScale.Valid {
			epsm.PayBandScale = payBandScale.String
		}
		if grade.Valid {
			epsm.Grade = grade.String
		}
		if gradeGroup.Valid {
			epsm.GradeGroup = gradeGroup.String
		}
		if status.Valid {
			epsm.Status = status.String
		}
		if upperLimit.Valid {
			epsm.UpperLimit = upperLimit.Float64
		}
		
		list = append(list, epsm)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}
	
	return list, nil
}