// Package modelssad defines models and SQL queries for Staff Additional Details.
//
// --- Creator's Info ---
// Creator: Rovita
// Created On: 11-11-2025
// Description: Structs and queries for Insert/Update of employee basic details.
package modelssad

import (
	"encoding/json"
)

// MasterSadPayload defines the complete structure for SAD using master_sad procedure
type MasterSadPayload struct {
	// Task and Process fields
	TaskID        string `json:"task_id"`
	ProcessID     int    `json:"process_id"`
	EmployeeID    string `json:"employee_id"`
	EmployeeName  string `json:"employee_name"`
	AssignTo      string `json:"assign_to"`
	AssignedRole  string `json:"assigned_role"`
	TaskStatusID  int    `json:"task_status_id"`
	ActivitySeqNo int    `json:"activity_seq_no"`
	InitiatedBy   string `json:"initiated_by"`
	UpdatedBy     string `json:"updated_by"`

	// Integer flag parameters
	IsTaskReturn   int    `json:"is_task_return"`
	IsTaskApproved int    `json:"is_task_approved"`
	EmailFlag      int    `json:"email_flag"`
	TemplateID     int    `json:"template_id"`
	RejectFlag     int    `json:"reject_flag"`
	RejectRole     string `json:"reject_role"`
	Badge          int    `json:"badge"`
	Priority       int    `json:"priority"`
	Starred        int    `json:"starred"`

	// Personal information
	FirstName              string   `json:"first_name"`
	MiddleName             *string  `json:"middle_name"`
	LastName               string   `json:"last_name"`
	Gender                 string   `json:"gender"`
	MaritalStatus          string   `json:"marital_status"`
	FatherName             string   `json:"father_name"`
	MotherName             string   `json:"mother_name"`
	SpouseName             *string  `json:"spouse_name"`
	DOB                    string   `json:"dob"`
	Age                    int      `json:"age"`
	Nationality            string   `json:"nationality"`
	Religion               string   `json:"religion"`
	CasteCategory          string   `json:"caste_category"`
	EmergencyContactNo     string   `json:"emergency_contact_no"`
	MobileNo               string   `json:"mobile_no"`
	IsPhysicallyChallenged int      `json:"is_physically_challenged"`
	PercentageOfDisability *float64 `json:"percentage_of_disability"`
	NatureOfDisability     *string  `json:"nature_of_disability"`
	PersonalEmail          string   `json:"personal_email"`
	AadhaarNo              string   `json:"aadhaar_no"`
	MotherTongue           string   `json:"mother_tongue"`
	BankName               string   `json:"bank_name"`
	IFSCCode               string   `json:"ifsc_code"`
	BankAcctNo             string   `json:"bank_acct_no"`
	IdentificationMarks    string   `json:"identification_marks"`
	PANCardNo              string   `json:"pan_card_no"`
	EmployeeType           string   `json:"employee_type"`
	Department             string   `json:"department"`
	Designation            string   `json:"designation"`
	Section                string   `json:"section"`
	RouteTo                string   `json:"route_to"`
	Grade                  string   `json:"grade"`
	EmpGroup               string   `json:"emp_group"`
	PayInfo                string   `json:"pay_info"`
	BasicPay               float64  `json:"basic_pay"`
	NonPracticePay         float64  `json:"non_practice_pay"`
	NameOfPayBand          string   `json:"name_of_pay_band"`
	DateOfJoining          string   `json:"date_of_joining"`
	DateOfConfirmation     string   `json:"date_of_confirmation"`
	OfficeRoomNo           string   `json:"office_room_no"`
	OfficeExtensionNo      string   `json:"office_extension_no"`
	IsActive               int      `json:"is_active"`
	EmployeeStatus         string   `json:"employee_status"`
	DateOfRetirement       string   `json:"date_of_retirement"`
	Comments               string   `json:"comments"`
	UserRole               string   `json:"user_role"`

	// Additional data sections (JSONB fields)
	ContactData    json.RawMessage `json:"contact_data"`
	DependentsData json.RawMessage `json:"dependents_data"`
	EducationData  json.RawMessage `json:"education_data"`
	ExperienceData json.RawMessage `json:"experience_data"`
	LanguageData   json.RawMessage `json:"language_data"`
	DocumentsData  json.RawMessage `json:"documents_data"`
}

// EducationData defines the structure for education information
type EducationData struct {
	DegreeOrExam       string `json:"degree_or_exam"`
	NameOfBoard        string `json:"name_of_board"`
	Institution        string `json:"institution"`
	UniversityName     string `json:"university_name"`
	Country            string `json:"country"`
	State              string `json:"state"`
	MonthYearOfPassing string `json:"month_year_of_passing"`
	RegNo              string `json:"reg_no"`
	Specialization     string `json:"specialization"`
	Mode               string `json:"mode"`
	PercentageOfMarks  string `json:"percentage_of_marks"`
	ObtainedMarks      string `json:"obtained_marks"`
	Class              string `json:"class"`
	InitiatedBy        string `json:"initiated_by"`
}

// ContactData defines the structure for contact information
type ContactData struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// DependentsData defines the structure for dependents information
type DependentsData struct {
	Name     string `json:"name"`
	Relation string `json:"relation"`
	DOB      string `json:"dob"`
}

// ExperienceData defines the structure for experience information
type ExperienceData struct {
	OrganizationName string `json:"organization_name"`
	From             string `json:"from"`
	To               string `json:"to"`
}

// LanguageData defines the structure for language information
type LanguageData struct {
	Language    string `json:"language"`
	Proficiency string `json:"proficiency"`
}

// DocumentsData defines the structure for documents information
type DocumentsData struct {
	DocType      string `json:"doc_type"`
	DocumentName string `json:"document_name"`
	DocumentPath string `json:"document_path"`
}
