
// Package modelssad contains structs and queries for Staff Additional details API.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 15-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all Combo Value Master.
package modelscommon

import (
	"database/sql"
	"fmt"
)

// ComboMasterQuery - Query to fetch combo values based on comboname
const ComboMasterQuery = `
SELECT id, comboname, displayseq, combovalue, isactive 
FROM humanresources.combovaluesmaster 
WHERE isactive='1' AND comboname=$1
ORDER BY displayseq, id
`

// ComboMaster struct to hold combo master data
type ComboMaster struct {
	ID         int    `json:"id"`
	ComboName  string `json:"comboname"`
	DisplaySeq string `json:"displayseq"`
	ComboValue string `json:"combovalue"`
	IsActive   string `json:"isactive"`
}

// RetrieveComboMaster scans combo master data from query results
func RetrieveComboMaster(rows *sql.Rows) ([]ComboMaster, error) {
	var list []ComboMaster
	for rows.Next() {
		var cm ComboMaster
		err := rows.Scan(
			&cm.ID,
			&cm.ComboName,
			&cm.DisplaySeq,
			&cm.ComboValue,
			&cm.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning combo master row: %v", err)
		}
		list = append(list, cm)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}
	
	return list, nil
}