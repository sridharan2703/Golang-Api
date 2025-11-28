// Package modelssad contains structs and queries for Staff Additional details API.
//
// --- Creator's Info ---
// Creator: Vaishnavi
// Created On: 04-11-2025
// Last Modified By:  Rovita
// Last Modified Date: 12-1-2025
// This api is to feth the all the contact details
package modelssad

import (
	"database/sql"
	"fmt"
)

// ContactDetails represents the raw contact data structure from database
type ContactDetails struct {
	EmployeeID         string
	PresentAddress1    string
	PresentAddress2    string
	PresentCountry     string
	PresentState       string
	PresentDistrict    string
	PresentCity        string
	PresentPincode     string
	PresentCountryCode string
	PresentAreaCode    string
	PermanentAddress1  string
	PermanentAddress2  string
	PermanentCountry   string
	PermanentState     string
	PermanentDistrict  string
	PermanentCity      string
	PermanentPincode   string
	IsQuarters         string
	QuartersNo         string
}

// ContactDetailsResponses represents the structured contact details response
type ContactDetailsResponses struct {
	EmployeeID       string                   `json:"employeeid"`
	CurrentAddress   CurrentAddressSections   `json:"current_address"`
	PermanentAddress PermanentAddressSections `json:"permanent_address"`
}

type CurrentAddressSections struct {
	Address1    string `json:"address1"`
	Address2    string `json:"address2"`
	Country     string `json:"country"`
	State       string `json:"state"`
	District    string `json:"district"`
	City        string `json:"city"`
	Pincode     string `json:"pincode"`
	CountryCode string `json:"country_code"`
	AreaCode    string `json:"area_code"`
}

type PermanentAddressSections struct {
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	Country  string `json:"country"`
	State    string `json:"state"`
	District string `json:"district"`
	City     string `json:"city"`
	Pincode  string `json:"pincode"`
}

// GenericContactDetailsRetriever retrieves and processes contact details from the table function
func GenericContactDetailsRetriever(db *sql.DB, employeeID string) (interface{}, error) {
	query := `
		SELECT 
			employeeid,
			present_address1,
			present_address2,
			present_country,
			present_state,
			present_district,
			present_city,
			present_pincode,
			present_countrycode,
			present_areacode,
			permanent_address1,
			permanent_address2,
			permanent_country,
			permanent_state,
			permanent_district,
			permanent_city,
			permanent_pincode,
			is_quarters,
			quarters_no
		FROM humanresources.get_employee_contact_details($1)
	`

	row := db.QueryRow(query, employeeID)

	var data ContactDetails
	err := row.Scan(
		&data.EmployeeID,
		&data.PresentAddress1,
		&data.PresentAddress2,
		&data.PresentCountry,
		&data.PresentState,
		&data.PresentDistrict,
		&data.PresentCity,
		&data.PresentPincode,
		&data.PresentCountryCode,
		&data.PresentAreaCode,
		&data.PermanentAddress1,
		&data.PermanentAddress2,
		&data.PermanentCountry,
		&data.PermanentState,
		&data.PermanentDistrict,
		&data.PermanentCity,
		&data.PermanentPincode,
		&data.IsQuarters,
		&data.QuartersNo,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return createEmptyContactDetails(employeeID), nil
		}
		return nil, fmt.Errorf("error scanning contact details: %v", err)
	}

	return buildContactDetailsResponse(data), nil
}

// buildContactDetailsResponse transforms raw ContactDetails into structured response
func buildContactDetailsResponse(data ContactDetails) ContactDetailsResponses {
	return ContactDetailsResponses{
		EmployeeID: data.EmployeeID,
		CurrentAddress: CurrentAddressSections{
			Address1:    data.PresentAddress1,
			Address2:    data.PresentAddress2,
			Country:     data.PresentCountry,
			State:       data.PresentState,
			District:    data.PresentDistrict,
			City:        data.PresentCity,
			Pincode:     data.PresentPincode,
			CountryCode: data.PresentCountryCode,
			AreaCode:    data.PresentAreaCode,
		},
		PermanentAddress: PermanentAddressSections{
			Address1: data.PermanentAddress1,
			Address2: data.PermanentAddress2,
			Country:  data.PermanentCountry,
			State:    data.PermanentState,
			District: data.PermanentDistrict,
			City:     data.PermanentCity,
			Pincode:  data.PermanentPincode,
		},
	}
}

// createEmptyContactDetails creates an empty contact details structure
func createEmptyContactDetails(employeeID string) ContactDetailsResponses {
	return ContactDetailsResponses{
		EmployeeID:       employeeID,
		CurrentAddress:   CurrentAddressSections{},
		PermanentAddress: PermanentAddressSections{},
	}
}
