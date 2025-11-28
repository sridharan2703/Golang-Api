// Package modelssad contains structs and queries for Staff Additional details API.
//
// --- Creator's Info ---
// Creator: Vaishnavi
// Created On: 04-11-2025
// Last Modified By:  Rovita
// Last Modified Date: 12-1-2025
// This api is to feth the all the education details
package modelssad

import (
	"database/sql"
	"fmt"
)

// EducationDetail represents a single education record
type EducationDetail struct {
	DegreeOrExam       string `json:"degree_or_exam"`
	BoardName          string `json:"board_name"`
	Institution        string `json:"institution"`
	UniversityName     string `json:"university_name"`
	EducationCountry   string `json:"education_country"`
	EducationState     string `json:"education_state"`
	MonthYearOfPassing string `json:"month_year_of_passing"`
	RegistrationNo     string `json:"registration_no"`
	Specialization     string `json:"specialization"`
	Mode               string `json:"mode"`
	PercentageOfMarks  string `json:"percentage_of_marks"`
	ObtainedMarks      string `json:"obtained_marks"`
	Class              string `json:"class"`
}

// GenericEducationDetailsRetriever retrieves and processes Education details
func GenericEducationDetailsRetriever(db *sql.DB, employeeID string) ([]EducationDetail, error) {
	rows, err := db.Query(`SELECT * FROM humanresources.get_employee_education_details($1)`, employeeID)
	if err != nil {
		return nil, fmt.Errorf("query execution error: %v", err)
	}
	defer rows.Close()

	var educationDetails []EducationDetail

	for rows.Next() {
		var detail EducationDetail
		var empID string // employeeid column, not used in output

		err := rows.Scan(
			&empID,
			&detail.DegreeOrExam,
			&detail.BoardName,
			&detail.Institution,
			&detail.UniversityName,
			&detail.EducationCountry,
			&detail.EducationState,
			&detail.MonthYearOfPassing,
			&detail.RegistrationNo,
			&detail.Specialization,
			&detail.Mode,
			&detail.PercentageOfMarks,
			&detail.ObtainedMarks,
			&detail.Class,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		educationDetails = append(educationDetails, detail)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return educationDetails, nil
}
