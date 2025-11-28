// Package modelssad contains structs and queries for Staff Additional details API.
//
// --- Creator's Info ---
// Creator: Vaishnavi
// Created On: 04-11-2025
// Last Modified By:  Rovita
// Last Modified Date: 12-1-2025
// This api is to feth the all the dependent details
package modelssad

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

// DatabaseDependent represents the raw structure from database
type DatabaseDependent struct {
	EmployeeID                  string  `json:"employeeid"`
	DependentName               string  `json:"dependent_name"`
	DependentRelationship       string  `json:"dependent_relationship"`
	DependentDOB                string  `json:"dependent_dob"`
	DependentAge                *int    `json:"dependent_age"`
	DependentMaritalStatus      *string `json:"dependent_marital_status"`
	DependentBloodGroup         *string `json:"dependent_blood_group"`
	DependentGender             *string `json:"dependent_gender"`
	DependentEmploymentStatus   *string `json:"dependent_employment_status"`
	DependentAadhaarNo          *string `json:"dependent_aadhaar_no"`
	IsTwins                     int     `json:"is_twins"` // Comes as 0/1 from database
	IsCurrentlyDependent        string  `json:"is_currently_dependent"`
	DependentMobileNo           *string `json:"dependent_mobile_no"`
	OptingForInsurance          int     `json:"opting_for_insurance"` // Comes as 0/1 from database
	OptingForLTC                int     `json:"opting_for_ltc"`       // Comes as 0/1 from database
	IsPersonDisabled            string  `json:"is_person_disabled"`
	DependentNatureOfDisability *string `json:"dependent_nature_of_disability"`
}

// ResponseDependent represents the clean structure for API response
type ResponseDependent struct {
	DependentName               string  `json:"dependent_name"`
	DependentRelationship       string  `json:"dependent_relationship"`
	DependentDOB                string  `json:"dependent_dob"`
	DependentAge                *int    `json:"dependent_age,omitempty"`
	DependentMaritalStatus      *string `json:"dependent_marital_status,omitempty"`
	DependentBloodGroup         *string `json:"dependent_blood_group,omitempty"`
	DependentGender             *string `json:"dependent_gender,omitempty"`
	DependentEmploymentStatus   *string `json:"dependent_employment_status,omitempty"`
	DependentAadhaarNo          *string `json:"dependent_aadhaar_no,omitempty"`
	IsTwins                     bool    `json:"is_twins"`
	DependentMobileNo           *string `json:"dependent_mobile_no,omitempty"`
	IsCurrentlyDependent        string  `json:"is_currently_dependent"`
	OptingForInsurance          bool    `json:"opting_for_insurance"`
	OptingForLTC                bool    `json:"opting_for_ltc"`
	IsPersonDisabled            string  `json:"is_person_disabled"`
	DependentNatureOfDisability *string `json:"dependent_nature_of_disability,omitempty"`
}

// DependentResponse is the final API response structure
type DependentResponse struct {
	EmployeeID string              `json:"employeeid"`
	Dependents []ResponseDependent `json:"dependents"`
}

// convertDatabaseToResponse converts database format to API response format and omits null/empty fields
func convertDatabaseToResponse(dbDependent DatabaseDependent) ResponseDependent {
	response := ResponseDependent{
		DependentName:         dbDependent.DependentName,
		DependentRelationship: dbDependent.DependentRelationship,
		DependentDOB:          dbDependent.DependentDOB,
		IsTwins:               dbDependent.IsTwins == 1,
		IsCurrentlyDependent:  dbDependent.IsCurrentlyDependent,
		OptingForInsurance:    dbDependent.OptingForInsurance == 1,
		OptingForLTC:          dbDependent.OptingForLTC == 1,
		IsPersonDisabled:      dbDependent.IsPersonDisabled,
	}

	// Only include fields that are not null/empty
	if dbDependent.DependentAge != nil {
		response.DependentAge = dbDependent.DependentAge
	}

	if dbDependent.DependentMaritalStatus != nil && *dbDependent.DependentMaritalStatus != "" {
		response.DependentMaritalStatus = dbDependent.DependentMaritalStatus
	}

	if dbDependent.DependentBloodGroup != nil && *dbDependent.DependentBloodGroup != "" {
		response.DependentBloodGroup = dbDependent.DependentBloodGroup
	}

	if dbDependent.DependentGender != nil && *dbDependent.DependentGender != "" {
		response.DependentGender = dbDependent.DependentGender
	}

	if dbDependent.DependentEmploymentStatus != nil && *dbDependent.DependentEmploymentStatus != "" {
		response.DependentEmploymentStatus = dbDependent.DependentEmploymentStatus
	}

	if dbDependent.DependentAadhaarNo != nil && *dbDependent.DependentAadhaarNo != "" {
		response.DependentAadhaarNo = dbDependent.DependentAadhaarNo
	}

	if dbDependent.DependentMobileNo != nil && *dbDependent.DependentMobileNo != "" {
		response.DependentMobileNo = dbDependent.DependentMobileNo
	}

	if dbDependent.DependentNatureOfDisability != nil && *dbDependent.DependentNatureOfDisability != "" {
		response.DependentNatureOfDisability = dbDependent.DependentNatureOfDisability
	}

	return response
}

// GenericDependentDetailsRetriever retrieves and processes dependent details
func GenericDependentDetailsRetriever(db *sql.DB, employeeID string) (interface{}, error) {
	log.Printf("Querying database for employee: %s", employeeID)

	// Query the database
	row := db.QueryRow(`SELECT humanresources.get_employee_dependent_details($1)`, employeeID)

	var jsonData []byte
	if err := row.Scan(&jsonData); err != nil {
		return nil, fmt.Errorf("error scanning JSON data: %v", err)
	}

	log.Printf("Raw JSON from database: %s", string(jsonData))

	// Try to unmarshal as array of database dependents
	var dbDependents []DatabaseDependent
	arrayErr := json.Unmarshal(jsonData, &dbDependents)

	if arrayErr == nil {
		// Successfully unmarshaled as array
		log.Printf("Successfully unmarshaled as array with %d elements", len(dbDependents))

		if len(dbDependents) == 0 {
			log.Println("No dependents found")
			return DependentResponse{
				EmployeeID: employeeID,
				Dependents: []ResponseDependent{},
			}, nil
		}

		// Convert all database records to response format
		responseDependents := make([]ResponseDependent, 0, len(dbDependents))
		for _, dbDep := range dbDependents {
			responseDependents = append(responseDependents, convertDatabaseToResponse(dbDep))
		}

		return DependentResponse{
			EmployeeID: employeeID,
			Dependents: responseDependents,
		}, nil
	}

	log.Printf("Array unmarshal failed: %v", arrayErr)

	// Try to unmarshal as single database dependent
	var singleDbDependent DatabaseDependent
	singleErr := json.Unmarshal(jsonData, &singleDbDependent)

	if singleErr == nil {
		// Successfully unmarshaled as single object
		log.Printf("Successfully unmarshaled as single object")

		responseDependent := convertDatabaseToResponse(singleDbDependent)

		return DependentResponse{
			EmployeeID: employeeID,
			Dependents: []ResponseDependent{responseDependent},
		}, nil
	}

	// Both attempts failed
	log.Printf("Single object unmarshal also failed: %v", singleErr)
	return nil, fmt.Errorf("JSON unmarshal error: expected array or object, got: %s", string(jsonData))
}
