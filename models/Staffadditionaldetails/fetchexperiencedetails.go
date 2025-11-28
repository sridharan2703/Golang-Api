// Package modelssad contains structs and queries for Staff Additional details API.
//
// --- Creator's Info ---
// Creator: Vaishnavi
// Created On: 04-11-2025
// Last Modified By:  Rovita
// Last Modified Date: 12-1-2025
// This api is to feth the all the experience details
package modelssad

import (
	"database/sql"
	"fmt"
)

// ExperienceDetail represents a single experience record
type ExperienceDetail struct {
	OrganizationName string `json:"organization_name"`
	Address1         string `json:"address1"`
	DesignationExp   string `json:"designation_experience"`
	FromDate         string `json:"from_date"`
	ToDate           string `json:"to_date"`
	TotalExperience  string `json:"total_experience"`
	PayScale         string `json:"pay_scale"`
	IsGovtEmployee   string `json:"is_govt_employee"`
	TypeOfEmployment string `json:"type_of_employment"`
}

// GenericExperienceDetailsRetriever retrieves and processes experience details
func GenericExperienceDetailsRetriever(db *sql.DB, employeeID string) ([]ExperienceDetail, error) {
	// Query all rows from the table-returning function
	rows, err := db.Query(`SELECT * FROM humanresources.get_employee_experience_details($1)`, employeeID)
	if err != nil {
		return nil, fmt.Errorf("query execution error: %v", err)
	}
	defer rows.Close()

	var experienceDetails []ExperienceDetail

	for rows.Next() {
		var detail ExperienceDetail
		var empID string // we'll read employeeid column but won't use it in the response

		err := rows.Scan(
			&empID,
			&detail.OrganizationName,
			&detail.Address1,
			&detail.DesignationExp,
			&detail.FromDate,
			&detail.ToDate,
			&detail.TotalExperience,
			&detail.PayScale,
			&detail.IsGovtEmployee,
			&detail.TypeOfEmployment,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		experienceDetails = append(experienceDetails, detail)
	}

	// Handle any iteration errors
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return experienceDetails, nil
}
