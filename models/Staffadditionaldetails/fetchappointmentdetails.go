// Package modelssad contains structs and queries for Staff Additional details API.
//
// --- Creator's Info ---
// Creator: Vaishnavi
// Created On: 04-11-2025
// Last Modified By:  Rovita
// Last Modified Date: 12-1-2025
// This api is to feth the all the appointment details
package modelssad

import (
	"database/sql"
	"fmt"
)

// EmployeeAppointmentDetails represents the complete appointment details structure from database
type EmployeeAppointmentDetails struct {
	EmployeeID         string `json:"employeeid"`
	EmployeeName       string `json:"employee_name"`
	EmployeeType       string `json:"employee_type"`
	Department         string `json:"department"`
	Designation        string `json:"designation"`
	Section            string `json:"section"`
	RouteTo            string `json:"route_to"`
	Grade              string `json:"grade"`
	EmployeeGroup      string `json:"employee_group"`
	EmployeeStatus     string `json:"employee_status"`
	IsActive           string `json:"is_active"`
	PayInfo            string `json:"pay_info"`
	BasicPay           string `json:"basic_pay"`
	NonPracticePay     string `json:"non_practice_pay"`
	NameOfPayBand      string `json:"name_of_pay_band"`
	DateOfJoining      string `json:"date_of_joining"`
	DateOfConfirmation string `json:"date_of_confirmation"`
	DateOfRetirement   string `json:"date_of_retirement"`
	OfficeRoomNo       string `json:"office_room_no"`
	OfficeExtensionNo  string `json:"office_extension_no"`
}

// AppointmentDetailsResponses represents the structured response
type AppointmentDetailsResponses struct {
	EmployeeID            string                        `json:"employeeid"`
	EmploymentInformation EmploymentInformationSections `json:"employment_information"`
	EmploymentDetails     EmploymentDetailsSections     `json:"employment_details"`
}

type EmploymentInformationSections struct {
	EmployeeName   string `json:"employee_name"`
	EmployeeType   string `json:"employee_type"`
	Department     string `json:"department"`
	Designation    string `json:"designation"`
	Section        string `json:"section"`
	RouteTo        string `json:"route_to"`
	Grade          string `json:"grade"`
	EmployeeGroup  string `json:"employee_group"`
	EmployeeStatus string `json:"employee_status"`
	IsActive       string `json:"is_active"`
}

type EmploymentDetailsSections struct {
	PayInfo            string `json:"pay_info"`
	BasicPay           string `json:"basic_pay"`
	NonPracticePay     string `json:"non_practice_pay"`
	NameOfPayBand      string `json:"name_of_pay_band"`
	DateOfJoining      string `json:"date_of_joining"`
	DateOfConfirmation string `json:"date_of_confirmation"`
	DateOfRetirement   string `json:"date_of_retirement"`
	OfficeRoomNo       string `json:"office_room_no"`
	OfficeExtensionNo  string `json:"office_extension_no"`
}

// BuildAppointmentDetailsResponse transforms flat EmployeeAppointmentDetails into structured response
func BuildAppointmentDetailsResponse(data EmployeeAppointmentDetails) AppointmentDetailsResponses {
	return AppointmentDetailsResponses{
		EmployeeID: data.EmployeeID,
		EmploymentInformation: EmploymentInformationSections{
			EmployeeName:   data.EmployeeName,
			EmployeeType:   data.EmployeeType,
			Department:     data.Department,
			Designation:    data.Designation,
			Section:        data.Section,
			RouteTo:        data.RouteTo,
			Grade:          data.Grade,
			EmployeeGroup:  data.EmployeeGroup,
			EmployeeStatus: data.EmployeeStatus,
			IsActive:       data.IsActive,
		},
		EmploymentDetails: EmploymentDetailsSections{
			PayInfo:            data.PayInfo,
			BasicPay:           data.BasicPay,
			NonPracticePay:     data.NonPracticePay,
			NameOfPayBand:      data.NameOfPayBand,
			DateOfJoining:      data.DateOfJoining,
			DateOfConfirmation: data.DateOfConfirmation,
			DateOfRetirement:   data.DateOfRetirement,
			OfficeRoomNo:       data.OfficeRoomNo,
			OfficeExtensionNo:  data.OfficeExtensionNo,
		},
	}
}

// GenericAppointmentDetailsRetriever retrieves and processes appointment details
func GenericAppointmentDetailsRetriever(db *sql.DB, employeeID string) (interface{}, error) {
	query := `SELECT * FROM humanresources.get_employee_appointment_details($1)`
	row := db.QueryRow(query, employeeID)

	var data EmployeeAppointmentDetails
	err := row.Scan(
		&data.EmployeeID,
		&data.EmployeeName,
		&data.EmployeeType,
		&data.Department,
		&data.Designation,
		&data.Section,
		&data.RouteTo,
		&data.Grade,
		&data.EmployeeGroup,
		&data.PayInfo,
		&data.BasicPay,
		&data.NonPracticePay,
		&data.NameOfPayBand,
		&data.DateOfJoining,
		&data.DateOfConfirmation,
		&data.OfficeRoomNo,
		&data.OfficeExtensionNo,
		&data.IsActive,
		&data.EmployeeStatus,
		&data.DateOfRetirement,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no appointment details found for employee ID: %s", employeeID)
		}
		return nil, fmt.Errorf("error scanning appointment details: %v", err)
	}

	return BuildAppointmentDetailsResponse(data), nil
}
