// Package modelssad contains structs and queries for Staff Additional details API.
//
// --- Creator's Info ---
// Creator: Vaishnavi
// Created On: 04-11-2025
// Last Modified By:  Rovita
// Last Modified Date: 12-1-2025
// This api is to feth the all the language details
package modelssad

import (
	"database/sql"
	"fmt"
)

// LanguageDetail represents a single language proficiency record
type LanguageDetail struct {
	Language string `json:"language"`
	Read     string `json:"read"`
	Write    string `json:"write"`
	Speak    string `json:"speak"`
}

// HindiProficiency represents Hindi language proficiency details
type HindiProficiency struct {
	HindiLevelOfKnowledge string `json:"hindi_level_of_knowledge"`
	HindiWorkingKnowledge string `json:"hindi_working_knowledge"`
	HindiProficiency      string `json:"hindi_proficiency"`
}

// EmployeeLanguageResponse represents the overall structured response
type EmployeeLanguageResponse struct {
	EmployeeID       string           `json:"employeeid"`
	Languages        []LanguageDetail `json:"languages"`
	HindiProficiency HindiProficiency `json:"hindi_proficiency"`
}

// GenericLanguageDetailsRetriever retrieves and processes language details from table-returning function
func GenericLanguageDetailsRetriever(db *sql.DB, employeeID string) (interface{}, error) {
	query := `SELECT employeeid, language, reads, writes, speaks,
	                 hindi_level_of_knowledge, hindi_working_knowledge, hindi_proficiency
	          FROM humanresources.get_employee_language_details($1)`

	rows, err := db.Query(query, employeeID)
	if err != nil {
		return nil, fmt.Errorf("query execution error: %v", err)
	}
	defer rows.Close()

	var languages []LanguageDetail
	var hindi HindiProficiency
	var empID string

	for rows.Next() {
		var lang, read, write, speak, hindiKnow, hindiWork, hindiProf sql.NullString
		if err := rows.Scan(&empID, &lang, &read, &write, &speak, &hindiKnow, &hindiWork, &hindiProf); err != nil {
			return nil, fmt.Errorf("row scan error: %v", err)
		}

		// Collect normal language info
		if lang.Valid && lang.String != "" {
			languages = append(languages, LanguageDetail{
				Language: lang.String,
				Read:     read.String,
				Write:    write.String,
				Speak:    speak.String,
			})
		}

		// Fill Hindi proficiency once
		if hindi.HindiLevelOfKnowledge == "" && (hindiKnow.Valid || hindiWork.Valid || hindiProf.Valid) {
			hindi = HindiProficiency{
				HindiLevelOfKnowledge: hindiKnow.String,
				HindiWorkingKnowledge: hindiWork.String,
				HindiProficiency:      hindiProf.String,
			}
		}
	}

	// Handle case where no rows were found
	if empID == "" {
		return EmployeeLanguageResponse{
			EmployeeID:       employeeID,
			Languages:        []LanguageDetail{},
			HindiProficiency: HindiProficiency{},
		}, nil
	}

	return EmployeeLanguageResponse{
		EmployeeID:       empID,
		Languages:        languages,
		HindiProficiency: hindi,
	}, nil
}
