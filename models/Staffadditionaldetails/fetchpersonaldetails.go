// Package modelssad contains structs and queries for Staff Additional details API.
//
// --- Creator's Info ---
// Creator: Vaishnavi
// Created On: 04-11-2025
// Last Modified By:  Rovita
// Last Modified Date: 12-1-2025
// This api is to feth the all the personal details
package modelssad

import (
	"database/sql"
	"fmt"
)

// PersonalDetails represents the complete personal details structure
type PersonalDetails struct {
	EmployeeID             string `json:"employeeid"`
	Firstname              string `json:"firstname"`
	Middlename             string `json:"middlename"`
	Lastname               string `json:"lastname"`
	Gender                 string `json:"gender"`
	MaritalStatus          string `json:"marital_status"`
	DOB                    string `json:"dob"`
	Age                    string `json:"age"`
	Nationality            string `json:"nationality"`
	BirthState             string `json:"birth_state"`
	BirthDistrict          string `json:"birth_district"`
	BirthPlace             string `json:"birth_place"`
	Hometown               string `json:"hometown"`
	Religion               string `json:"religion"`
	CasteCategory          string `json:"caste_category"`
	EmergencyContactNo     string `json:"emergency_contact_no"`
	MobileNo               string `json:"mobile_no"`
	IsPhysicallyChallenged string `json:"is_physically_challenged"`
	PercentageOfDisability string `json:"percentage_of_disability"`
	NatureOfDisability     string `json:"nature_of_disability"`
	PersonalEmail          string `json:"personal_email"`
	AadhaarNo              string `json:"aadhaar_no"`
	MotherTongue           string `json:"mother_tongue"`
	PanCardNo              string `json:"pan_card_no"`
	FatherName             string `json:"father_name"`
	MotherName             string `json:"mother_name"`
	SpouseName             string `json:"spouse_name"`
	BankName               string `json:"bank_name"`
	IFSCCode               string `json:"ifsc_code"`
	BankAccountNo          string `json:"bank_account_no"`
}

// Structured response (your same layout)
type PersonalDetailsStructuredResponse struct {
	EmployeeID      string             `json:"employeeid"`
	PersonalDetails PersonalDetailsSec `json:"personal_details"`
	FamilyDetails   FamilyDetailsSec   `json:"family_details"`
	BankDetails     BankDetailsSec     `json:"bank_details"`
}

type PersonalDetailsSec struct {
	FirstName              string `json:"first_name"`
	MiddleName             string `json:"middle_name"`
	LastName               string `json:"last_name"`
	Gender                 string `json:"gender"`
	MaritalStatus          string `json:"marital_status"`
	DOB                    string `json:"dob"`
	Age                    string `json:"age"`
	Nationality            string `json:"nationality"`
	BirthState             string `json:"birth_state"`
	BirthDistrict          string `json:"birth_district"`
	BirthPlace             string `json:"birth_place"`
	Hometown               string `json:"hometown"`
	Religion               string `json:"religion"`
	CasteCategory          string `json:"caste_category"`
	EmergencyContactNo     string `json:"emergency_contact_no"`
	MobileNo               string `json:"mobile_no"`
	IsPhysicallyChallenged string `json:"is_physically_challenged"`
	PercentageOfDisability string `json:"percentage_of_disability"`
	NatureOfDisability     string `json:"nature_of_disability"`
	PersonalEmail          string `json:"personal_email"`
	AadhaarNo              string `json:"aadhaar_no"`
	MotherTongue           string `json:"mother_tongue"`
	PanCardNo              string `json:"pan_card_no"`
}

type FamilyDetailsSec struct {
	FatherName string `json:"father_name"`
	MotherName string `json:"mother_name"`
	SpouseName string `json:"spouse_name"`
}

type BankDetailsSec struct {
	BankName      string `json:"bank_name"`
	IFSCCode      string `json:"ifsc_code"`
	BankAccountNo string `json:"bank_account_no"`
}

// Transform flat DB row into structured JSON
func BuildPersonalDetailsResponse(data PersonalDetails) PersonalDetailsStructuredResponse {
	return PersonalDetailsStructuredResponse{
		EmployeeID: data.EmployeeID,
		PersonalDetails: PersonalDetailsSec{
			FirstName:              data.Firstname,
			MiddleName:             data.Middlename,
			LastName:               data.Lastname,
			Gender:                 data.Gender,
			MaritalStatus:          data.MaritalStatus,
			DOB:                    data.DOB,
			Age:                    data.Age,
			Nationality:            data.Nationality,
			BirthState:             data.BirthState,
			BirthDistrict:          data.BirthDistrict,
			BirthPlace:             data.BirthPlace,
			Hometown:               data.Hometown,
			Religion:               data.Religion,
			CasteCategory:          data.CasteCategory,
			EmergencyContactNo:     data.EmergencyContactNo,
			MobileNo:               data.MobileNo,
			IsPhysicallyChallenged: data.IsPhysicallyChallenged,
			PercentageOfDisability: data.PercentageOfDisability,
			NatureOfDisability:     data.NatureOfDisability,
			PersonalEmail:          data.PersonalEmail,
			AadhaarNo:              data.AadhaarNo,
			MotherTongue:           data.MotherTongue,
			PanCardNo:              data.PanCardNo,
		},
		FamilyDetails: FamilyDetailsSec{
			FatherName: data.FatherName,
			MotherName: data.MotherName,
			SpouseName: data.SpouseName,
		},
		BankDetails: BankDetailsSec{
			BankName:      data.BankName,
			IFSCCode:      data.IFSCCode,
			BankAccountNo: data.BankAccountNo,
		},
	}
}

// ✅ Updated retriever (for TABLE-returning SQL function)
func GenericPersonalDetailsRetriever(db *sql.DB, employeeID string) (interface{}, error) {
	row := db.QueryRow(`SELECT * FROM humanresources.get_employee_personal_details($1)`, employeeID)

	var data PersonalDetails

	err := row.Scan(
		&data.EmployeeID,
		&data.Firstname,
		&data.Middlename,
		&data.Lastname,
		&data.Gender,
		&data.MaritalStatus,
		&data.FatherName,
		&data.MotherName,
		&data.SpouseName,
		&data.DOB,
		&data.Age,
		&data.Nationality,
		&data.BirthState,
		&data.BirthDistrict,
		&data.BirthPlace,
		&data.Hometown,
		&data.Religion,
		&data.CasteCategory,
		&data.EmergencyContactNo,
		&data.MobileNo,
		&data.IsPhysicallyChallenged,
		&data.PercentageOfDisability,
		&data.NatureOfDisability,
		&data.PersonalEmail,
		&data.AadhaarNo,
		&data.MotherTongue,
		&data.BankName,
		&data.IFSCCode,
		&data.BankAccountNo,
		&data.PanCardNo,
	)
	if err != nil {
		return nil, fmt.Errorf("error scanning personal details: %v", err)
	}

	return BuildPersonalDetailsResponse(data), nil
}
