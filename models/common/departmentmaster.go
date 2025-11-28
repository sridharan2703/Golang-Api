// Package modelscommon contains structs and queries for DepartmentMaster.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 17-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all DepartmentMaster Value Master.

package modelscommon

import (
	"database/sql"
	"fmt"
)

// DepartmentMasterQuery - Query to fetch DepartmentMaster values
const DepartmentMasterQuery = `
SELECT id, campusid, departmentcode, departmentname, type 
FROM humanresources.departmentmaster
`

// DepartmentMaster struct to hold DepartmentMaster data
type DepartmentMaster struct {
	ID             int    `json:"id"`
	CampusID       int    `json:"campusid"`
	DepartmentCode string `json:"departmentcode"`
	DepartmentName string `json:"departmentname"`
	Type           string `json:"type"`
}

// RetrieveDepartmentMaster scans DepartmentMaster data from query results
func RetrieveDepartmentMaster(rows *sql.Rows) ([]DepartmentMaster, error) {
	var list []DepartmentMaster
	for rows.Next() {
		var dm DepartmentMaster
		err := rows.Scan(
			&dm.ID,
			&dm.CampusID,
			&dm.DepartmentCode,
			&dm.DepartmentName,
			&dm.Type,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning DepartmentMaster row: %v", err)
		}
		list = append(list, dm)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}
	
	return list, nil
}