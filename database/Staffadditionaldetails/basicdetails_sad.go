// // Package databasesad interacts with the Staff Additional Details DB.
// //
// // --- Creator's Info ---
// // Creator: Rovita
// // Created On: 11-11-2025
// // Last Modified By:
// // Last Modified Date:
// // Description: Insert and Update operations for employee basic details.
package databasesad

import (
	credentials "Hrmodule/dbconfig"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
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

// CallMasterSADStoredProcedure calls the master_sad stored procedure for complete SAD operations
func CallMasterSADStoredProcedure(action string, data MasterSadPayload) (int64, error) {
	db, err := sql.Open("postgres", credentials.Getdatabasehr())
	if err != nil {
		return 0, fmt.Errorf("DB connection failed: %v", err)
	}
	defer db.Close()

	// Build the complete SQL query with DO block
	query := buildMasterSADQuery(action, data)

	log.Printf("Executing master_sad for employee: %s", data.EmployeeID)

	// Execute the DO block
	_, err = db.Exec(query)
	if err != nil {
		log.Printf("Error calling master_sad stored procedure: %v", err)
		return 0, fmt.Errorf("error executing master_sad stored procedure: %v", err)
	}

	// Since the stored procedure works asynchronously with the DO block,
	// we need to query for the SAD_ID separately
	// You can either:
	// 1. Query by task_id to get the sad_id
	// 2. Use RETURNING clause if the procedure supports it
	// 3. Return a placeholder value

	var sadID int64
	taskID := data.TaskID
	if taskID == "" {
		// If task_id was generated, we can't easily retrieve the sad_id
		log.Printf("Successfully executed master_sad procedure for action: %s, employee: %s",
			action, data.EmployeeID)
		return 0, nil
	}

	// Query to get the SAD_ID by task_id
	queryID := `SELECT sad_id FROM meivan.staff_additional_details WHERE task_id = $1 ORDER BY sad_id DESC LIMIT 1`
	err = db.QueryRow(queryID, taskID).Scan(&sadID)
	if err != nil {
		log.Printf("Warning: Could not retrieve SAD_ID for task_id %s: %v", taskID, err)
		// Don't fail the request, as the procedure executed successfully
		return 0, nil
	}

	log.Printf("Successfully executed master_sad procedure for action: %s, employee: %s, SAD ID: %d",
		action, data.EmployeeID, sadID)

	return sadID, nil
}

// buildMasterSADQuery builds the complete SQL query with DO block and named parameters
func buildMasterSADQuery(action string, data MasterSadPayload) string {
	// Handle task_id - generate new UUID if empty
	taskIDValue := "gen_random_uuid()"
	if data.TaskID != "" {
		taskIDValue = fmt.Sprintf("'%s'::uuid", data.TaskID)
	}

	// Handle NULL values for string pointers
	middleNameValue := "NULL::varchar"
	if data.MiddleName != nil && *data.MiddleName != "" {
		middleNameValue = fmt.Sprintf("'%s'::varchar", escapeSQLString(*data.MiddleName))
	}

	spouseNameValue := "NULL::varchar"
	if data.SpouseName != nil && *data.SpouseName != "" {
		spouseNameValue = fmt.Sprintf("'%s'::varchar", escapeSQLString(*data.SpouseName))
	}

	natureOfDisabilityValue := "NULL::varchar"
	if data.NatureOfDisability != nil && *data.NatureOfDisability != "" {
		natureOfDisabilityValue = fmt.Sprintf("'%s'::varchar", escapeSQLString(*data.NatureOfDisability))
	}

	rejectRoleValue := "NULL::varchar"
	if data.RejectRole != "" {
		rejectRoleValue = fmt.Sprintf("'%s'::varchar", escapeSQLString(data.RejectRole))
	}

	// Handle percentage of disability
	percentageOfDisabilityValue := "0::numeric"
	if data.PercentageOfDisability != nil {
		percentageOfDisabilityValue = fmt.Sprintf("%f::numeric", *data.PercentageOfDisability)
	}

	// Handle JSON data - validate and format properly
	contactDataValue := "NULL::jsonb"
	if isValidJSON(data.ContactData) {
		contactDataValue = fmt.Sprintf("'%s'::jsonb", escapeSQLString(string(data.ContactData)))
	}

	dependentsDataValue := "NULL::jsonb"
	if isValidJSON(data.DependentsData) {
		dependentsDataValue = fmt.Sprintf("'%s'::jsonb", escapeSQLString(string(data.DependentsData)))
	}

	educationDataValue := "NULL::jsonb"
	if isValidJSON(data.EducationData) {
		educationDataValue = fmt.Sprintf("'%s'::jsonb", escapeSQLString(string(data.EducationData)))
	}

	experienceDataValue := "NULL::jsonb"
	if isValidJSON(data.ExperienceData) {
		experienceDataValue = fmt.Sprintf("'%s'::jsonb", escapeSQLString(string(data.ExperienceData)))
	}

	languageDataValue := "NULL::jsonb"
	if isValidJSON(data.LanguageData) {
		languageDataValue = fmt.Sprintf("'%s'::jsonb", escapeSQLString(string(data.LanguageData)))
	}

	documentsDataValue := "NULL::jsonb"
	if isValidJSON(data.DocumentsData) {
		documentsDataValue = fmt.Sprintf("'%s'::jsonb", escapeSQLString(string(data.DocumentsData)))
	}

	// Handle empty strings and NULL for text fields
	commentsValue := "NULL::text"
	if data.Comments != "" {
		commentsValue = fmt.Sprintf("'%s'::text", escapeSQLString(data.Comments))
	}

	userRoleValue := "NULL::varchar"
	if data.UserRole != "" {
		userRoleValue = fmt.Sprintf("'%s'::varchar", escapeSQLString(data.UserRole))
	}

	// Build the complete DO block with named parameters (exactly like your working examples)
	query := fmt.Sprintf(`
DO $$
DECLARE
    v_task_id UUID := %s;
    v_sad_id BIGINT;
BEGIN
    CALL meivan.master_sad(
        p_action_type       => '%s'::varchar,
        p_task_id           => v_task_id,
        p_process_id        => %d::integer,
        p_employee_id       => '%s'::varchar,
        p_employee_name     => '%s'::varchar,
        p_assign_to         => '%s'::varchar,
        p_assigned_role     => '%s'::varchar,
        p_task_status_id    => %d::integer,
        p_activity_seq_no   => %d::integer,
        p_is_task_return    => %d::integer,
        p_is_task_approved  => %d::integer,
        p_email_flag        => %d::integer,
        p_template_id       => %d::integer,
        p_reject_flag       => %d::smallint,
        p_reject_role       => %s,
        p_initiated_by      => '%s'::varchar,
        p_updated_by        => '%s'::varchar,
        p_badge             => %d::integer,
        p_priority          => %d::integer,
        p_starred           => %d::integer,
        p_first_name        => '%s'::varchar,
        p_middle_name       => %s,
        p_last_name         => '%s'::varchar,
        p_gender            => '%s'::varchar,
        p_marital_status    => '%s'::varchar,
        p_father_name       => '%s'::varchar,
        p_mother_name       => '%s'::varchar,
        p_spouse_name       => %s,
        p_dob               => '%s'::date,
        p_age               => %d::integer,
        p_nationality       => '%s'::varchar,
        p_religion          => '%s'::varchar,
        p_caste_category    => '%s'::varchar,
        p_emergency_contact_no => '%s'::varchar,
        p_mobile_no         => '%s'::varchar,
        p_is_physically_challenged => %d::integer,
        p_percentage_of_disability  => %s,
        p_nature_of_disability      => %s,
        p_personal_email    => '%s'::varchar,
        p_aadhaar_no        => '%s'::varchar,
        p_mother_tongue     => '%s'::varchar,
        p_bank_name         => '%s'::varchar,
        p_ifsc_code         => '%s'::varchar,
        p_bank_acct_no      => '%s'::varchar,
        p_identification_marks => '%s'::varchar,
        p_pan_card_no       => '%s'::varchar,
        p_employee_type     => '%s'::varchar,
        p_department        => '%s'::varchar,
        p_designation       => '%s'::varchar,
        p_section           => '%s'::varchar,
        p_route_to          => '%s'::varchar,
        p_grade             => '%s'::varchar,
        p_emp_group         => '%s'::varchar,
        p_pay_info          => '%s'::varchar,
        p_basic_pay         => %f::numeric,
        p_non_practice_pay  => %f::numeric,
        p_name_of_pay_band  => '%s'::varchar,
        p_date_of_joining   => '%s'::date,
        p_date_of_confirmation => '%s'::date,
        p_office_room_no    => '%s'::varchar,
        p_office_extension_no => '%s'::varchar,
        p_is_active         => %d::integer,
        p_employee_status   => '%s'::varchar,
        p_date_of_retirement => '%s'::date,
        p_contact_data      => %s,
        p_dependents_data   => %s,
        p_education_data    => %s,
        p_experience_data   => %s,
        p_language_data     => %s,
        p_documents_data    => %s,
        p_comments          => %s,
        p_user_role         => %s,
        v_sad_id            => v_sad_id
    );

    RAISE NOTICE '✅ SAD created with ID: %%', v_sad_id;
END $$;
`,
		// Task and basic info
		taskIDValue, action, data.ProcessID,
		escapeSQLString(data.EmployeeID), escapeSQLString(data.EmployeeName),
		escapeSQLString(data.AssignTo), escapeSQLString(data.AssignedRole),
		data.TaskStatusID, data.ActivitySeqNo, data.IsTaskReturn, data.IsTaskApproved,
		data.EmailFlag, data.TemplateID, data.RejectFlag, rejectRoleValue,
		escapeSQLString(data.InitiatedBy), escapeSQLString(data.UpdatedBy),
		data.Badge, data.Priority, data.Starred,

		// Personal information
		escapeSQLString(data.FirstName), middleNameValue, escapeSQLString(data.LastName),
		escapeSQLString(data.Gender), escapeSQLString(data.MaritalStatus),
		escapeSQLString(data.FatherName), escapeSQLString(data.MotherName), spouseNameValue,
		data.DOB, data.Age, escapeSQLString(data.Nationality), escapeSQLString(data.Religion),
		escapeSQLString(data.CasteCategory), escapeSQLString(data.EmergencyContactNo),
		escapeSQLString(data.MobileNo), data.IsPhysicallyChallenged, percentageOfDisabilityValue,
		natureOfDisabilityValue, escapeSQLString(data.PersonalEmail), escapeSQLString(data.AadhaarNo),
		escapeSQLString(data.MotherTongue), escapeSQLString(data.BankName), escapeSQLString(data.IFSCCode),
		escapeSQLString(data.BankAcctNo), escapeSQLString(data.IdentificationMarks),
		escapeSQLString(data.PANCardNo), escapeSQLString(data.EmployeeType), escapeSQLString(data.Department),
		escapeSQLString(data.Designation), escapeSQLString(data.Section), escapeSQLString(data.RouteTo),
		escapeSQLString(data.Grade), escapeSQLString(data.EmpGroup), escapeSQLString(data.PayInfo),
		data.BasicPay, data.NonPracticePay, escapeSQLString(data.NameOfPayBand),
		data.DateOfJoining, data.DateOfConfirmation, escapeSQLString(data.OfficeRoomNo),
		escapeSQLString(data.OfficeExtensionNo), data.IsActive, escapeSQLString(data.EmployeeStatus),
		data.DateOfRetirement,

		// JSON data
		contactDataValue, dependentsDataValue, educationDataValue, experienceDataValue,
		languageDataValue, documentsDataValue,

		// Final fields
		commentsValue, userRoleValue,
	)

	return query
}

// escapeSQLString escapes single quotes in SQL strings
func escapeSQLString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

// isValidJSON checks if the byte array contains valid JSON
func isValidJSON(data json.RawMessage) bool {
	if data == nil || len(data) == 0 {
		return false
	}

	str := strings.TrimSpace(string(data))

	if len(str) == 0 || str == "null" {
		return false
	}

	// Basic validation - should start with [ or {
	if !strings.HasPrefix(str, "[") && !strings.HasPrefix(str, "{") {
		return false
	}

	// Try to parse it to validate
	var js interface{}
	return json.Unmarshal(data, &js) == nil
}
