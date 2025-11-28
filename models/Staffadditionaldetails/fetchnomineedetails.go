// Package modelssad contains structs and queries for Staff Additional details API.
//
// --- Creator's Info ---
// Creator: Vaishnavi
// Created On: 04-11-2025
// Last Modified By:  Rovita
// Last Modified Date: 12-1-2025
// This api is to feth the all the nominee details
package modelssad

import (
	"database/sql"
	"fmt"
)

// NomineeInfo represents nominee information for each scheme
type NomineeInfo struct {
	Name            string `json:"name"`
	Relationship    string `json:"relationship"`
	DOB             string `json:"dob"`
	SharePercentage string `json:"share_percentage"`
}

// NomineeDetails represents all nominee details for an employee
type NomineeDetails struct {
	EmployeeID      string      `json:"employeeid"`
	GratuityNominee NomineeInfo `json:"gratuity_nominee"`
	GTISNominee     NomineeInfo `json:"gtis_nominee"`
	NPSNominee      NomineeInfo `json:"nps_nominee"`
	GPFNominee      NomineeInfo `json:"gpf_nominee"`
}

// GenericNomineeDetailsRetriever retrieves and processes nominee details safely handling NULLs
func GenericNomineeDetailsRetriever(db *sql.DB, employeeID string) (NomineeDetails, error) {
	query := `
		SELECT
			employeeid,
			gratuity_nominee_name,
			gratuity_nominee_relationship,
			gratuity_nominee_dob,
			gratuity_nominee_share_percentage,
			gtis_nominee_name,
			gtis_nominee_relationship,
			gtis_nominee_dob,
			gtis_nominee_share_percentage,
			nps_nominee_name,
			nps_nominee_relationship,
			nps_nominee_dob,
			nps_nominee_share_percentage,
			gpf_nominee_name,
			gpf_nominee_relationship,
			gpf_nominee_dob,
			gpf_nominee_share_percentage
		FROM humanresources.get_employee_nominee_details($1);
	`

	// use NullString for safe scanning
	var (
		empID                                                 sql.NullString
		gratuityName, gratuityRel, gratuityDOB, gratuityShare sql.NullString
		gtisName, gtisRel, gtisDOB, gtisShare                 sql.NullString
		npsName, npsRel, npsDOB, npsShare                     sql.NullString
		gpfName, gpfRel, gpfDOB, gpfShare                     sql.NullString
	)

	err := db.QueryRow(query, employeeID).Scan(
		&empID,
		&gratuityName, &gratuityRel, &gratuityDOB, &gratuityShare,
		&gtisName, &gtisRel, &gtisDOB, &gtisShare,
		&npsName, &npsRel, &npsDOB, &npsShare,
		&gpfName, &gpfRel, &gpfDOB, &gpfShare,
	)

	if err == sql.ErrNoRows {
		return createEmptyNomineeDetails(employeeID), nil
	} else if err != nil {
		return createEmptyNomineeDetails(employeeID), fmt.Errorf("error fetching nominee details: %v", err)
	}

	nominee := NomineeDetails{
		EmployeeID:      nullToString(empID),
		GratuityNominee: NomineeInfo{nullToString(gratuityName), nullToString(gratuityRel), nullToString(gratuityDOB), nullToString(gratuityShare)},
		GTISNominee:     NomineeInfo{nullToString(gtisName), nullToString(gtisRel), nullToString(gtisDOB), nullToString(gtisShare)},
		NPSNominee:      NomineeInfo{nullToString(npsName), nullToString(npsRel), nullToString(npsDOB), nullToString(npsShare)},
		GPFNominee:      NomineeInfo{nullToString(gpfName), nullToString(gpfRel), nullToString(gpfDOB), nullToString(gpfShare)},
	}

	return nominee, nil
}

// helper: safely convert sql.NullString to string
func nullToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// helper: creates an empty nominee details structure
func createEmptyNomineeDetails(employeeID string) NomineeDetails {
	return NomineeDetails{
		EmployeeID:      employeeID,
		GratuityNominee: NomineeInfo{},
		GTISNominee:     NomineeInfo{},
		NPSNominee:      NomineeInfo{},
		GPFNominee:      NomineeInfo{},
	}
}
