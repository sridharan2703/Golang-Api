// Package modelscommon contains structs and queries for BloodGroup.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 15-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all BloodGroup Master.

package modelscommon

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// BloodGroupMasterQuery - Query to fetch BloodGroup values
const BloodGroupMasterQuery = `
SELECT id, bloodgroupname, description, isactive 
FROM humanresources.bloodgroupmaster 
WHERE isactive='1'
ORDER BY id
`

// BloodGroupMaster struct to hold BloodGroup master data
type BloodGroupMaster struct {
	ID             int            `json:"id"`
	BloodGroupName sql.NullString `json:"bloodgroupname"`
	Description    sql.NullString `json:"description"`
	IsActive       string         `json:"isactive"`
}

// MarshalJSON custom marshaller to handle NullString properly
func (bg BloodGroupMaster) MarshalJSON() ([]byte, error) {
	type Alias BloodGroupMaster
	return json.Marshal(&struct {
		BloodGroupName interface{} `json:"bloodgroupname"`
		Description    interface{} `json:"description"`
		*Alias
	}{
		BloodGroupName: getStringValue(bg.BloodGroupName),
		Description:    getStringValue(bg.Description),
		Alias:          (*Alias)(&bg),
	})
}

// Helper function to get string value from NullString
func getStringValue(ns sql.NullString) interface{} {
	if ns.Valid {
		return ns.String
	}
	return nil
}

// RetrieveBloodGroupMaster scans BloodGroup master data from query results
func RetrieveBloodGroupMaster(rows *sql.Rows) ([]BloodGroupMaster, error) {
	var list []BloodGroupMaster
	for rows.Next() {
		var bm BloodGroupMaster
		err := rows.Scan(
			&bm.ID,
			&bm.BloodGroupName,
			&bm.Description,
			&bm.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning BloodGroup master row: %v", err)
		}
		list = append(list, bm)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}
	
	return list, nil
}