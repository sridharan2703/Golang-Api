// Package modelscommon contains structs and queries for languagemaster.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 15-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all languagemaster.

package modelscommon

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// LanguageMasterQuery - Query to fetch language master values
const LanguageMasterQuery = `
SELECT id, langcode, langname, countrycode, isactive 
FROM humanresources.languagemaster 
WHERE isactive='1'
ORDER BY id
`

// LanguageMaster struct to hold language master data
type LanguageMaster struct {
	ID          int            `json:"id"`
	LangCode    string         `json:"langcode"`
	LangName    string         `json:"langname"`
	CountryCode sql.NullString `json:"countrycode"`
	IsActive    string         `json:"isactive"`
}

// MarshalJSON custom marshaller to handle NullString properly
func (lm LanguageMaster) MarshalJSON() ([]byte, error) {
	type Alias LanguageMaster
	return json.Marshal(&struct {
		CountryCode interface{} `json:"countrycode"`
		*Alias
	}{
		CountryCode: getStringValuelang(lm.CountryCode),
		Alias:       (*Alias)(&lm),
	})
}

// Helper function to get string value from NullString
func getStringValuelang(ns sql.NullString) interface{} {
	if ns.Valid {
		return ns.String
	}
	return nil
}

// RetrieveLanguageMaster scans language master data from query results
func RetrieveLanguageMaster(rows *sql.Rows) ([]LanguageMaster, error) {
	var list []LanguageMaster
	for rows.Next() {
		var lm LanguageMaster
		err := rows.Scan(
			&lm.ID,
			&lm.LangCode,
			&lm.LangName,
			&lm.CountryCode,
			&lm.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning language master row: %v", err)
		}
		list = append(list, lm)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}
	
	return list, nil
}