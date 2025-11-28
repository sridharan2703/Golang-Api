
// Package modelscommon contains structs and queries for Designation.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 17-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all Designation Value Master.

package modelscommon

import (
	"database/sql"
	"fmt"
)

// DesignationMasterQuery - Query to fetch Designation values
const DesignationMasterQuery = `
SELECT designationid, cadreid, designationname, designationdescription 
FROM humanresources.designationmaster
`

// DesignationMaster struct to hold Designation master data
type DesignationMaster struct {
	DesignationID          int    `json:"designationid"`
	CadreID               int    `json:"cadreid"`
	DesignationName       string `json:"designationname"`
	DesignationDescription string `json:"designationdescription"`
}

// RetrieveDesignationMaster scans Designation master data from query results
func RetrieveDesignationMaster(rows *sql.Rows) ([]DesignationMaster, error) {
	var list []DesignationMaster
	for rows.Next() {
		var rm DesignationMaster
		err := rows.Scan(
			&rm.DesignationID,
			&rm.CadreID,
			&rm.DesignationName,
			&rm.DesignationDescription,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning Designation master row: %v", err)
		}
		list = append(list, rm)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}
	
	return list, nil
}