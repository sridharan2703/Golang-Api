// Package modelscommon contains structs and queries for castecagetory.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 15-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all castecagetory Master.

package modelscommon

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// CasteCategoryMasterQuery - Query to fetch caste category values
const CasteCategoryMasterQuery = `
SELECT id, name, description, isactive 
FROM humanresources.castecategory 
WHERE isactive='1'
ORDER BY id
`

// CasteCategoryMaster struct to hold caste category master data
type CasteCategoryMaster struct {
	ID          int            `json:"id"`
	Name        sql.NullString `json:"name"`
	Description sql.NullString `json:"description"`
	IsActive    string         `json:"isactive"`
}

// MarshalJSON custom marshaller to handle NullString properly
func (cc CasteCategoryMaster) MarshalJSON() ([]byte, error) {
	type Alias CasteCategoryMaster
	return json.Marshal(&struct {
		Name        interface{} `json:"name"`
		Description interface{} `json:"description"`
		*Alias
	}{
		Name:        getStringValues(cc.Name),
		Description: getStringValues(cc.Description),
		Alias:       (*Alias)(&cc),
	})
}

// Helper function to get string value from NullString
func getStringValues(ns sql.NullString) interface{} {
	if ns.Valid {
		return ns.String
	}
	return nil
}

// RetrieveCasteCategoryMaster scans caste category master data from query results
func RetrieveCasteCategoryMaster(rows *sql.Rows) ([]CasteCategoryMaster, error) {
	var list []CasteCategoryMaster
	for rows.Next() {
		var cm CasteCategoryMaster
		err := rows.Scan(
			&cm.ID,
			&cm.Name,
			&cm.Description,
			&cm.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning caste category master row: %v", err)
		}
		list = append(list, cm)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}
	
	return list, nil
}