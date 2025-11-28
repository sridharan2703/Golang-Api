// Package modelscommon contains structs and queries for Bank.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 18-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all Bank Value Master.

package modelscommon

import (
	"database/sql"
	"fmt"
)

// BankMasterQuery - Query to fetch Bank values
const BankMasterQuery = `
SELECT id, bankname, sequenceid, isactive 
FROM humanresources.bankmaster 
WHERE isactive = '1'
ORDER BY sequenceid, bankname
`

// BankMaster struct to hold Bank master data
type BankMaster struct {
	ID         int    `json:"id"`
	BankName   string `json:"bankname"`
	SequenceID string `json:"sequenceid"`
	IsActive   string `json:"isactive"`
}

// RetrieveBankMaster scans Bank master data from query results
func RetrieveBankMaster(rows *sql.Rows) ([]BankMaster, error) {
	var list []BankMaster
	for rows.Next() {
		var rm BankMaster
		err := rows.Scan(
			&rm.ID,
			&rm.BankName,
			&rm.SequenceID,
			&rm.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning Bank master row: %v", err)
		}
		list = append(list, rm)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}
	
	return list, nil
}