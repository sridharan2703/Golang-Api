// Package modelscommon contains structs and queries for Staff Additional details API.

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

// OfficialLanguageQuery - Query to fetch official language values based on questiontype
const OfficialLanguageQuery = `
SELECT id, questiontype, optiontext, isactive 
FROM humanresources.officiallanguage 
WHERE isactive='1' AND questiontype=$1
ORDER BY id
`

// OfficialLanguage struct to hold official language data
type OfficialLanguage struct {
	ID          int    `json:"id"`
	QuestionType string `json:"questiontype"`
	OptionText  string `json:"optiontext"`
	IsActive    string `json:"isactive"`
}

// RetrieveOfficialLanguage scans official language data from query results
func RetrieveOfficialLanguage(rows *sql.Rows) ([]OfficialLanguage, error) {
	var list []OfficialLanguage
	for rows.Next() {
		var ol OfficialLanguage
		err := rows.Scan(
			&ol.ID,
			&ol.QuestionType,
			&ol.OptionText,
			&ol.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning official language row: %v", err)
		}
		list = append(list, ol)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}
	
	return list, nil
}