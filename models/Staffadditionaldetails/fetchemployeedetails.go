// Package modelssad contains structs and queries for Staff Additional details API.
//
// --- Creator's Info ---
// Creator: Vaishnavi
// Created On: 04-11-2025
// Last Modified By:  Rovita
// Last Modified Date: 12-1-2025
// This api is to feth the all the employee details
package modelssad

import (
	"database/sql"
	"fmt"
	"strings"
)

// -------------------- PostgreSQL Queries --------------------

// Complete query with all fields from all categories
const completeEFileQuery = `
SELECT
    COALESCE(ebi.employeeid, '') AS employeeid,
    COALESCE(ebi.firstname, '') AS firstname,
    COALESCE(ebi.middlename, '') AS middlename,
    COALESCE(ebi.lastname, '') AS lastname,
    COALESCE(cbv.combovalue, '') AS gender,
    COALESCE(cbv1.combovalue, '') AS marital_status,
    COALESCE(ebi.fatherorhusbandname, '') AS father_name,
    COALESCE(ebi.mothername, '') AS mother_name,
    COALESCE(ebi.spousename, '') AS spouse_name,
    COALESCE(TO_CHAR(ebi.dob, 'YYYY-MM-DD'), '') AS dob,
    COALESCE(EXTRACT(YEAR FROM AGE(ebi.dob))::text, '') AS age,
    COALESCE(ebi.nationality, '') AS nationality,
    COALESCE(ebi.birthstate, '') AS birth_state,
    COALESCE(ebi.placeofbirth, '') AS birth_district,
    COALESCE(ebi.birthtown, '') AS birth_place,
    COALESCE(ebi.hometown, '') AS hometown,
    COALESCE(rel.name, '') AS religion,
    COALESCE(cc.name, '') AS caste_category,
    COALESCE(ebi.emergencycontactno, '') AS emergency_contact_no,
    COALESCE(ebi.mobilenumber, '') AS mobile_no,
    COALESCE(ebi.isphysicallychallenged::text, '') AS is_physically_challenged,
    COALESCE(ebi.percentageofdisability::text, '') AS percentage_of_disability,
    COALESCE(ebi.physicallychallengedtype, '') AS nature_of_disability,
    COALESCE(ebi.primarymail, '') AS personal_email,
    COALESCE(ebd.aadhaar_no, '') AS aadhaar_no,
    COALESCE(ebi.mothertongue, '') AS mother_tongue,
    COALESCE(ebd.bank_name, '') AS bank_name,
    COALESCE(ebd.ifsc_code, '') AS ifsc_code,
    COALESCE(ebd.bank_account_no, '') AS bank_account_no,
    '' AS identification_marks,
    COALESCE(ebd.pan_card_no, '') AS pan_card_no,

    -- ================================
    -- APPOINTMENT DETAILS
    -- ================================
    COALESCE(ebi.displayname, '') AS employee_name,
    COALESCE(cbv2.combovalue, '') AS employee_type,
    COALESCE(dm.departmentname, '') AS department,
    COALESCE(dsm.designationname::text, '') AS designation,
    COALESCE(ead.section, '') AS section,
    '' AS route_to,
    COALESCE(epm.grade::text, '') AS grade,
    COALESCE(ead.employeegroup, '') AS employee_group,
    COALESCE(epm.paybandscale, '') AS pay_info,
    COALESCE(ead.basicpay::text, '') AS basic_pay,
    COALESCE(ead.nonpractisepay::text, '') AS non_practice_pay,
    COALESCE(epm.paybandname, '') AS name_of_pay_band,
    COALESCE(TO_CHAR(ead.doj, 'YYYY-MM-DD'), '') AS date_of_joining,
    COALESCE(TO_CHAR(ead.confirmationdate, 'YYYY-MM-DD'), '') AS date_of_confirmation,
    COALESCE(ead.roomno, '') AS office_room_no,
    COALESCE(ebi.extension, '') AS office_extension_no,
    COALESCE(ebi.employeestatus, '') AS is_active,
    COALESCE(ebi.employeestatus, '') AS employee_status,
    COALESCE(TO_CHAR(ead.retirementdate, 'YYYY-MM-DD'), '') AS date_of_retirement,

    -- ================================
    -- EDUCATIONAL QUALIFICATION DETAILS
    -- ================================
    COALESCE(eed.degree, '') AS degree_or_exam,
    COALESCE(eed.board, '') AS board_name,
    COALESCE(eed.institution, '') AS institution,
    COALESCE(eed.universityname, '') AS university_name,
    COALESCE(eed.universitycountry, '') AS education_country,
    COALESCE(eed.universitystate, '') AS education_state,
    COALESCE(eed.yearofpassing::text, '') AS month_year_of_passing,
    COALESCE(eed.registrationnumber, '') AS registration_no,
    COALESCE(eed.specialization, '') AS specialization,
    COALESCE(eed.mode, '') AS mode,
    COALESCE(eed.marks::text, '') AS percentage_of_marks,
    COALESCE(eed.obtainedmarks::text, '') AS obtained_marks,
    COALESCE(eed.class, '') AS class,

    -- ================================
    -- EXPERIENCE DETAILS
    -- ================================
    COALESCE(eep.name, '') AS organization_name,
    COALESCE(eep.address1, '') AS address1,
    COALESCE(eep.designation, '') AS designation_experience,
    COALESCE(TO_CHAR(eep.periodfrom, 'YYYY-MM-DD'), '') AS from_date,
    COALESCE(TO_CHAR(eep.periodto, 'YYYY-MM-DD'), '') AS to_date,
    COALESCE(eep.experience::text, '') AS total_experience,
    COALESCE(eep.payscale::text, '') AS pay_scale,
    COALESCE(eep.isgovtemployee::text, '') AS is_govt_employee,
    COALESCE(eep.employeetype, '') AS type_of_employment,

    -- ================================
    -- DEPENDENT INFO
    -- ================================
    COALESCE(edd.NAME, '') AS dependent_name,
    COALESCE(cbv3.combovalue, '') AS dependent_relationship,
    COALESCE(TO_CHAR(edd.dob, 'YYYY-MM-DD'), '') AS dependent_dob,
    COALESCE(edd.age::text, '') AS dependent_age,
    COALESCE(edd.maritalstatus, '') AS dependent_marital_status,
    COALESCE(bgm.bloodgroupname, '') AS dependent_blood_group,
    COALESCE(edd.gender, '') AS dependent_gender,
    COALESCE(edd.employmentstatus, '') AS dependent_employment_status,
    COALESCE(edd.aadharno, '') AS dependent_aadhaar_no,
    COALESCE(edd.istwins::text, '') AS is_twins,
    COALESCE(edd.dependantmobileno, '') AS dependent_mobile_no,
    COALESCE(CASE WHEN edd.status='A' THEN 'YES' ELSE 'NO' END, '') AS is_currently_dependent,
    COALESCE(edd.isinsured::text, '') AS opting_for_insurance,
    COALESCE(edd.ltc::text, '') AS opting_for_ltc,
    COALESCE(edd.ischilddisabled::text, '') AS is_person_disabled,
    COALESCE(edd.natureofchilddisability, '') AS dependent_nature_of_disability,

    -- ================================
    -- NOMINEE INFO
    -- ================================
    COALESCE(CASE WHEN nomineetype = 'gratuity' THEN edn.nomineename ELSE '' END, '') AS gratuity_nominee_name,
    COALESCE(CASE WHEN nomineetype = 'gratuity' THEN edn.relationship ELSE '' END, '') AS gratuity_nominee_relationship,
    COALESCE(TO_CHAR(CASE WHEN nomineetype = 'gratuity' THEN edn.dob ELSE NULL END, 'YYYY-MM-DD'), '') AS gratuity_nominee_dob,
    COALESCE(CASE WHEN nomineetype = 'gratuity' THEN edn.sharepercentage::text ELSE '' END, '') AS gratuity_nominee_share_percentage,

    COALESCE(CASE WHEN nomineetype = 'gtis' THEN edn.nomineename ELSE '' END, '') AS gtis_nominee_name,
    COALESCE(CASE WHEN nomineetype = 'gtis' THEN edn.relationship ELSE '' END, '') AS gtis_nominee_relationship,
    COALESCE(TO_CHAR(CASE WHEN nomineetype = 'gtis' THEN edn.dob ELSE NULL END, 'YYYY-MM-DD'), '') AS gtis_nominee_dob,
    COALESCE(CASE WHEN nomineetype = 'gtis' THEN edn.sharepercentage::text ELSE '' END, '') AS gtis_nominee_share_percentage,

    COALESCE(CASE WHEN nomineetype = 'nps' THEN edn.nomineename ELSE '' END, '') AS nps_nominee_name,
    COALESCE(CASE WHEN nomineetype = 'nps' THEN edn.relationship ELSE '' END, '') AS nps_nominee_relationship,
    COALESCE(TO_CHAR(CASE WHEN nomineetype = 'nps' THEN edn.dob ELSE NULL END, 'YYYY-MM-DD'), '') AS nps_nominee_dob,
    COALESCE(CASE WHEN nomineetype = 'nps' THEN edn.sharepercentage::text ELSE '' END, '') AS nps_nominee_share_percentage,

    COALESCE(CASE WHEN nomineetype = 'gpf' THEN edn.nomineename ELSE '' END, '') AS gpf_nominee_name,
    COALESCE(CASE WHEN nomineetype = 'gpf' THEN edn.relationship ELSE '' END, '') AS gpf_nominee_relationship,
    COALESCE(TO_CHAR(CASE WHEN nomineetype = 'gpf' THEN edn.dob ELSE NULL END, 'YYYY-MM-DD'), '') AS gpf_nominee_dob,
    COALESCE(CASE WHEN nomineetype = 'gpf' THEN edn.sharepercentage::text ELSE '' END, '') AS gpf_nominee_share_percentage,

    -- ================================
    -- CONTACT DETAILS
    -- ================================
    COALESCE(ecd.present_address1, '') AS present_address1,
    COALESCE(ecd.present_address2, '') AS present_address2,
    COALESCE(ecd.present_country, '') AS present_country,
    COALESCE(ecd.present_state, '') AS present_state,
    COALESCE(ecd.present_district, '') AS present_district,
    COALESCE(ecd.present_city, '') AS present_city,
    COALESCE(ecd.present_pincode, '') AS present_pincode,
    COALESCE(ecd.present_countrycode, '') AS present_countrycode,
    COALESCE(ecd.present_areacode, '') AS present_areacode,
    COALESCE(ecd.permanent_address1, '') AS permanent_address1,
    COALESCE(ecd.permanent_address2, '') AS permanent_address2,
    COALESCE(ecd.permanent_country, '') AS permanent_country,
    COALESCE(ecd.permanent_state, '') AS permanent_state,
    COALESCE(ecd.permanent_district, '') AS permanent_district,
    COALESCE(ecd.permanent_city, '') AS permanent_city,
    COALESCE(ecd.permanent_pincode, '') AS permanent_pincode,
    COALESCE(ecd.is_quarters::text, '') AS is_quarters,
    COALESCE(ecd.quarters_no, '') AS quarters_no,

    -- ================================
    -- LANGUAGE PROFICIENCY
    -- ================================
    COALESCE(eld.languagename, '') AS language,
    COALESCE(eld.reads::text, '') AS read,
    COALESCE(eld.writes::text, '') AS write,
    COALESCE(eld.speaks::text, '') AS speak,
    COALESCE(eld.official_lang_knowledge, '') AS hindi_level_of_knowledge,
    COALESCE(eld.official_lang_working, '') AS hindi_working_knowledge,
    COALESCE(eld.official_lang_proficiency, '') AS hindi_proficiency,

    -- ================================
    -- DOCUMENTS
    -- ================================
    COALESCE(emd.documentname, '') AS document_name


FROM humanresources.employeebasicinfo ebi
LEFT JOIN humanresources.employeeappointmentdetails ead 
    ON ebi.employeeid = ead.employeeid
LEFT JOIN humanresources.employeeeducationdetails eed 
    ON ebi.employeeid = eed.employeeid
LEFT JOIN humanresources.employeeexperiencedetails eep 
    ON ebi.employeeid = eep.employeeid
LEFT JOIN humanresources.employeelanguagedetails eld 
    ON ebi.employeeid = eld.employeeid
LEFT JOIN humanresources.employeecreationdocument emd 
    ON ebi.employeeid = emd.employeeid
	LEFT JOIN humanresources.employeedependantdetails edd 
    ON ebi.employeeid = edd.employeeid and edd.status='A'
		LEFT JOIN humanresources.employeenomineedetails edn 
    ON ebi.employeeid = edn.employeeid 
			LEFT JOIN humanresources.employeecontactdetails ecd 
    ON ebi.employeeid = ecd.employeeid 
				LEFT JOIN humanresources.employeebankdetails ebd 
    ON ebi.employeeid = ebd.employeeid 
			LEFT JOIN humanresources.departmentdesignationmapping ddm 
    ON ebi.employeeid = ddm.employeeid 
	left join humanresources.designationmaster dsm
	on ddm.designationid = dsm.designationid
	JOIN humanresources.combovaluesmaster cbv
	on ebi.gender=cbv.displayseq and cbv.comboname='Gender' and cbv.isactive='1'
		JOIN humanresources.combovaluesmaster cbv1
	on ebi.emppermaritalstatus=cbv1.displayseq and cbv1.comboname='MaritalStatus' and cbv1.isactive='1'
		JOIN humanresources.religion rel
	on ebi.religion=rel.id and rel.isactive='1'
			JOIN humanresources.castecategory cc
	on ebi.caste=cc.id and cc.isactive='1'
			JOIN humanresources.combovaluesmaster cbv2
	on ead.employeetype=cbv2.displayseq and cbv2.comboname='EmployeeType' and cbv2.isactive='1'
	left join humanresources.departmentmaster dm
	on ead.deptcode=dm.departmentcode
	left join  humanresources.employeepresentscalemaster epm
	on ead.presentscaleid = epm.id
			left JOIN humanresources.combovaluesmaster cbv3
	on edd.relationship=cbv3.displayseq and cbv3.comboname='RelationShip' and cbv3.isactive='1'
		left join  humanresources.bloodgroupmaster bgm 
	on edd.bloodgroup = bgm.id and bgm.isactive='1'
WHERE %s
ORDER BY ebi.employeeid ASC`

// GetEFileQueryByCategory returns the complete query with employee ID filter
func GetEFileQueryByCategory(category string) string {
	whereClause := "ebi.employeeid = $1"
	return fmt.Sprintf(completeEFileQuery, whereClause)
}

// GetEFileQueryAllByCategory returns the complete query for all employees
func GetEFileQueryAllByCategory(category string) string {
	whereClause := "1=1"
	return fmt.Sprintf(completeEFileQuery, whereClause)
}

// EmployeeEfileDetails struct to hold complete employee E-File data
type EmployeeEfileDetails struct {
	// Personal Details
	EmployeeID             string `json:"employeeid"`
	FirstName              string `json:"firstname"`
	MiddleName             string `json:"middlename"`
	LastName               string `json:"lastname"`
	Gender                 string `json:"gender"`
	MaritalStatus          string `json:"marital_status"`
	FatherName             string `json:"father_name"`
	MotherName             string `json:"mother_name"`
	SpouseName             string `json:"spouse_name"`
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
	BankName               string `json:"bank_name"`
	IFSCCode               string `json:"ifsc_code"`
	BankAccountNo          string `json:"bank_account_no"`
	IdentificationMarks    string `json:"identification_marks"`
	PanCardNo              string `json:"pan_card_no"`

	// Appointment Details
	EmployeeName       string `json:"employee_name"`
	EmployeeType       string `json:"employee_type"`
	Department         string `json:"department"`
	Designation        string `json:"designation"`
	Section            string `json:"section"`
	RouteTo            string `json:"route_to"`
	Grade              string `json:"grade"`
	EmployeeGroup      string `json:"employee_group"`
	PayInfo            string `json:"pay_info"`
	BasicPay           string `json:"basic_pay"`
	NonPracticePay     string `json:"non_practice_pay"`
	NameOfPayBand      string `json:"name_of_pay_band"`
	DateOfJoining      string `json:"date_of_joining"`
	DateOfConfirmation string `json:"date_of_confirmation"`
	OfficeRoomNo       string `json:"office_room_no"`
	OfficeExtensionNo  string `json:"office_extension_no"`
	IsActive           string `json:"is_active"`
	EmployeeStatus     string `json:"employee_status"`
	DateOfRetirement   string `json:"date_of_retirement"`

	// Education Details
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

	// Experience Details
	OrganizationName string `json:"organization_name"`
	Address1         string `json:"address1"`
	DesignationExp   string `json:"designation_experience"`
	FromDate         string `json:"from_date"`
	ToDate           string `json:"to_date"`
	TotalExperience  string `json:"total_experience"`
	PayScale         string `json:"pay_scale"`
	IsGovtEmployee   string `json:"is_govt_employee"`
	TypeOfEmployment string `json:"type_of_employment"`

	// Dependent Info
	DependentName             string `json:"dependent_name"`
	DependentRelationship     string `json:"dependent_relationship"`
	DependentDOB              string `json:"dependent_dob"`
	DependentAge              string `json:"dependent_age"`
	DependentMaritalStatus    string `json:"dependent_marital_status"`
	DependentBloodGroup       string `json:"dependent_blood_group"`
	DependentGender           string `json:"dependent_gender"`
	DependentEmploymentStatus string `json:"dependent_employment_status"`
	DependentAadhaarNo        string `json:"dependent_aadhaar_no"`
	IsTwins                   string `json:"is_twins"`
	DependentMobileNo         string `json:"dependent_mobile_no"`

	IsCurrentlyDependent        string `json:"is_currently_dependent"`
	OptingForInsurance          string `json:"opting_for_insurance"`
	OptingForLTC                string `json:"opting_for_ltc"`
	IsPersonDisabled            string `json:"is_person_disabled"`
	DependentNatureOfDisability string `json:"dependent_nature_of_disability"`

	// Nominee Info
	GratuityNomineeName            string `json:"gratuity_nominee_name"`
	GratuityNomineeRelationship    string `json:"gratuity_nominee_relationship"`
	GratuityNomineeDOB             string `json:"gratuity_nominee_dob"`
	GratuityNomineeSharePercentage string `json:"gratuity_nominee_share_percentage"`
	GTISNomineeName                string `json:"gtis_nominee_name"`
	GTISNomineeRelationship        string `json:"gtis_nominee_relationship"`
	GTISNomineeDOB                 string `json:"gtis_nominee_dob"`
	GTISNomineeSharePercentage     string `json:"gtis_nominee_share_percentage"`
	NPSNomineeName                 string `json:"nps_nominee_name"`
	NPSNomineeRelationship         string `json:"nps_nominee_relationship"`
	NPSNomineeDOB                  string `json:"nps_nominee_dob"`
	NPSNomineeSharePercentage      string `json:"nps_nominee_share_percentage"`
	GPFNomineeName                 string `json:"gpf_nominee_name"`
	GPFNomineeRelationship         string `json:"gpf_nominee_relationship"`
	GPFNomineeDOB                  string `json:"gpf_nominee_dob"`
	GPFNomineeSharePercentage      string `json:"gpf_nominee_share_percentage"`

	// Contact Details
	PresentAddress1    string `json:"present_address1"`
	PresentAddress2    string `json:"present_address2"`
	PresentCountry     string `json:"present_country"`
	PresentState       string `json:"present_state"`
	PresentDistrict    string `json:"present_district"`
	PresentCity        string `json:"present_city"`
	PresentPincode     string `json:"present_pincode"`
	PresentCountryCode string `json:"present_countrycode"`
	PresentAreaCode    string `json:"present_areacode"`
	PermanentAddress1  string `json:"permanent_address1"`
	PermanentAddress2  string `json:"permanent_address2"`
	PermanentCountry   string `json:"permanent_country"`
	PermanentState     string `json:"permanent_state"`
	PermanentDistrict  string `json:"permanent_district"`
	PermanentCity      string `json:"permanent_city"`
	PermanentPincode   string `json:"permanent_pincode"`
	IsQuarters         string `json:"is_quarters"`
	QuartersNo         string `json:"quarters_no"`

	// Language Proficiency
	Language string `json:"language"`
	Read     string `json:"read"`
	Write    string `json:"write"`
	Speak    string `json:"speak"`

	HindiLevelOfKnowledge string `json:"hindi_level_of_knowledge"`
	HindiWorkingKnowledge string `json:"hindi_working_knowledge"`
	HindiProficiency      string `json:"hindi_proficiency"`

	// Documents
	DocumentName string `json:"document_name"`
}

// CategoryResponse struct to hold category-specific response
// This generic struct can handle any category type
type CategoryResponse struct {
	EmployeeID string      `json:"employeeid"`
	Data       interface{} `json:"data"`
}

// RetrieveEmployeeEFile scans complete employee data from query results
// This function handles scanning all fields from the database
func RetrieveEmployeeEFile(rows *sql.Rows, category string) ([]CategoryResponse, error) {
	var list []CategoryResponse
	for rows.Next() {
		var e EmployeeEfileDetails
		err := scanCompleteEmployeeDetails(rows, &e)
		if err != nil {
			return nil, fmt.Errorf("error scanning employee data row: %v", err)
		}

		// Create category-specific response
		categoryData := createCategoryResponse(e, category)
		list = append(list, CategoryResponse{
			EmployeeID: e.EmployeeID,
			Data:       categoryData,
		})
	}
	return list, nil
}

// createCategoryResponse creates a category-specific response with only relevant fields
// This function filters the complete data to return only the requested category
func createCategoryResponse(e EmployeeEfileDetails, category string) interface{} {
	switch strings.ToLower(category) {
	case "personaldetails":
		return PersonalDetailsResponse{
			EmployeeID: e.EmployeeID,
			PersonalDetails: PersonalDetailsSection{
				FirstName:              e.FirstName,
				MiddleName:             e.MiddleName,
				LastName:               e.LastName,
				Gender:                 e.Gender,
				MaritalStatus:          e.MaritalStatus,
				DOB:                    e.DOB,
				Age:                    e.Age,
				Nationality:            e.Nationality,
				BirthState:             e.BirthState,
				BirthDistrict:          e.BirthDistrict,
				BirthPlace:             e.BirthPlace,
				Hometown:               e.Hometown,
				Religion:               e.Religion,
				CasteCategory:          e.CasteCategory,
				EmergencyContactNo:     e.EmergencyContactNo,
				MobileNo:               e.MobileNo,
				IsPhysicallyChallenged: e.IsPhysicallyChallenged,
				PercentageOfDisability: e.PercentageOfDisability,
				NatureOfDisability:     e.NatureOfDisability,
				PersonalEmail:          e.PersonalEmail,
				AadhaarNo:              e.AadhaarNo,
				MotherTongue:           e.MotherTongue,
				IdentificationMarks:    e.IdentificationMarks,
				PanCardNo:              e.PanCardNo,
			},
			FamilyDetails: FamilyDetailsSection{
				FatherName: e.FatherName,
				MotherName: e.MotherName,
				SpouseName: e.SpouseName,
			},
			BankDetails: BankDetailsSection{
				BankName:      e.BankName,
				IFSCCode:      e.IFSCCode,
				BankAccountNo: e.BankAccountNo,
			},
		}
	case "appointmentdetails":
		return AppointmentDetailsResponse{
			EmployeeID: e.EmployeeID,
			EmploymentInformation: EmploymentInformationSection{
				EmployeeName:   e.EmployeeName,
				EmployeeType:   e.EmployeeType,
				Department:     e.Department,
				Designation:    e.Designation,
				Section:        e.Section,
				RouteTo:        e.RouteTo,
				Grade:          e.Grade,
				EmployeeGroup:  e.EmployeeGroup,
				EmployeeStatus: e.EmployeeStatus,
				IsActive:       e.IsActive,
			},
			EmploymentDetails: EmploymentDetailsSection{
				PayInfo:            e.PayInfo,
				BasicPay:           e.BasicPay,
				NonPracticePay:     e.NonPracticePay,
				NameOfPayBand:      e.NameOfPayBand,
				DateOfJoining:      e.DateOfJoining,
				DateOfConfirmation: e.DateOfConfirmation,
				DateOfRetirement:   e.DateOfRetirement,
				OfficeRoomNo:       e.OfficeRoomNo,
				OfficeExtensionNo:  e.OfficeExtensionNo,
			},
		}
	case "educationdetails":
		return EducationDetailsResponse{
			EmployeeID:         e.EmployeeID,
			DegreeOrExam:       e.DegreeOrExam,
			BoardName:          e.BoardName,
			Institution:        e.Institution,
			UniversityName:     e.UniversityName,
			EducationCountry:   e.EducationCountry,
			EducationState:     e.EducationState,
			MonthYearOfPassing: e.MonthYearOfPassing,
			RegistrationNo:     e.RegistrationNo,
			Specialization:     e.Specialization,
			Mode:               e.Mode,
			PercentageOfMarks:  e.PercentageOfMarks,
			ObtainedMarks:      e.ObtainedMarks,
			Class:              e.Class,
		}
	case "experiencedetails":
		return ExperienceDetailsResponse{
			EmployeeID:       e.EmployeeID,
			OrganizationName: e.OrganizationName,
			Address1:         e.Address1,
			DesignationExp:   e.DesignationExp,
			FromDate:         e.FromDate,
			ToDate:           e.ToDate,
			TotalExperience:  e.TotalExperience,
			PayScale:         e.PayScale,
			IsGovtEmployee:   e.IsGovtEmployee,
			TypeOfEmployment: e.TypeOfEmployment,
		}
	case "languagedetails":
		return LanguageDetailsResponse{
			EmployeeID: e.EmployeeID,
			LanguageProficiency: LanguageProficiencySection{
				Language: e.Language,
				Read:     e.Read,
				Write:    e.Write,
				Speak:    e.Speak,
			},
			HindiProficiency: HindiProficiencySection{
				LevelOfKnowledge: e.HindiLevelOfKnowledge,
				WorkingKnowledge: e.HindiWorkingKnowledge,
				Proficiency:      e.HindiProficiency,
			},
		}
	case "documentdetails":
		return DocumentDetailsResponse{
			EmployeeID:   e.EmployeeID,
			DocumentName: e.DocumentName,
		}
	case "dependentdetails":
		return DependentDetailsResponse{
			EmployeeID:                e.EmployeeID,
			DependentName:             e.DependentName,
			DependentRelationship:     e.DependentRelationship,
			DependentDOB:              e.DependentDOB,
			DependentAge:              e.DependentAge,
			DependentMaritalStatus:    e.DependentMaritalStatus,
			DependentBloodGroup:       e.DependentBloodGroup,
			DependentGender:           e.DependentGender,
			DependentEmploymentStatus: e.DependentEmploymentStatus,
			DependentAadhaarNo:        e.DependentAadhaarNo,
			IsTwins:                   e.IsTwins,
			DependentMobileNo:         e.DependentMobileNo,

			IsCurrentlyDependent:        e.IsCurrentlyDependent,
			OptingForInsurance:          e.OptingForInsurance,
			OptingForLTC:                e.OptingForLTC,
			IsPersonDisabled:            e.IsPersonDisabled,
			DependentNatureOfDisability: e.DependentNatureOfDisability,
		}

	case "nomineedetails":
		// For nominee details, return the structured response directly without nesting
		return map[string]interface{}{
			"employeeid": e.EmployeeID,

			"gratuity_nominee": map[string]string{
				"name":             e.GratuityNomineeName,
				"relationship":     e.GratuityNomineeRelationship,
				"dob":              e.GratuityNomineeDOB,
				"share_percentage": e.GratuityNomineeSharePercentage,
			},
			"gtis_nominee": map[string]string{
				"name":             e.GTISNomineeName,
				"relationship":     e.GTISNomineeRelationship,
				"dob":              e.GTISNomineeDOB,
				"share_percentage": e.GTISNomineeSharePercentage,
			},
			"nps_nominee": map[string]string{
				"name":             e.NPSNomineeName,
				"relationship":     e.NPSNomineeRelationship,
				"dob":              e.NPSNomineeDOB,
				"share_percentage": e.NPSNomineeSharePercentage,
			},
			"gpf_nominee": map[string]string{
				"name":             e.GPFNomineeName,
				"relationship":     e.GPFNomineeRelationship,
				"dob":              e.GPFNomineeDOB,
				"share_percentage": e.GPFNomineeSharePercentage,
			},
		}
	case "contactdetails":
		return ContactDetailsResponse{
			EmployeeID: e.EmployeeID,
			CurrentAddress: CurrentAddressSection{
				Address1: e.PresentAddress1,
				Address2: e.PresentAddress2,
				Country:  e.PresentCountry,
				State:    e.PresentState,
				District: e.PresentDistrict,
				City:     e.PresentCity,
				Pincode:  e.PresentPincode,
			},
			PermanentAddress: PermanentAddressSection{
				Address1: e.PermanentAddress1,
				Address2: e.PermanentAddress2,
				Country:  e.PermanentCountry,
				State:    e.PermanentState,
				District: e.PermanentDistrict,
				City:     e.PermanentCity,
				Pincode:  e.PermanentPincode,
			},
		}
	case "hindiproficiency":
		return HindiProficiencyResponse{
			EmployeeID:            e.EmployeeID,
			HindiLevelOfKnowledge: e.HindiLevelOfKnowledge,
			HindiWorkingKnowledge: e.HindiWorkingKnowledge,
			HindiProficiency:      e.HindiProficiency,
		}
	default:
		// Default to personal details if category is not specifically handled
		return PersonalDetailsResponse{
			EmployeeID: e.EmployeeID,
			PersonalDetails: PersonalDetailsSection{
				FirstName: e.FirstName,
				LastName:  e.LastName,
			},
		}
	}
}

// ================================
// PERSONAL DETAILS - Split into sections
// ================================

type PersonalDetailsResponse struct {
	EmployeeID      string                 `json:"employeeid"`
	PersonalDetails PersonalDetailsSection `json:"personal_details"`
	FamilyDetails   FamilyDetailsSection   `json:"family_details"`
	BankDetails     BankDetailsSection     `json:"bank_details"`
}

type PersonalDetailsSection struct {
	FirstName              string `json:"firstname"`
	MiddleName             string `json:"middlename"`
	LastName               string `json:"lastname"`
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
	IdentificationMarks    string `json:"identification_marks"`
	PanCardNo              string `json:"pan_card_no"`
}

type FamilyDetailsSection struct {
	FatherName string `json:"father_name"`
	MotherName string `json:"mother_name"`
	SpouseName string `json:"spouse_name"`
}

type BankDetailsSection struct {
	BankName      string `json:"bank_name"`
	IFSCCode      string `json:"ifsc_code"`
	BankAccountNo string `json:"bank_account_no"`
}

// ================================
// APPOINTMENT DETAILS - Split into sections
// ================================

type AppointmentDetailsResponse struct {
	EmployeeID            string                       `json:"employeeid"`
	EmploymentInformation EmploymentInformationSection `json:"employment_information"`
	EmploymentDetails     EmploymentDetailsSection     `json:"employment_details"`
}

type EmploymentInformationSection struct {
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

type EmploymentDetailsSection struct {
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

// ================================
// CONTACT DETAILS - Split into sections
// ================================

type ContactDetailsResponse struct {
	EmployeeID       string                  `json:"employeeid"`
	CurrentAddress   CurrentAddressSection   `json:"current_address"`
	PermanentAddress PermanentAddressSection `json:"permanent_address"`
}

type CurrentAddressSection struct {
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	Country  string `json:"country"`
	State    string `json:"state"`
	District string `json:"district"`
	City     string `json:"city"`
	Pincode  string `json:"pincode"`
}

type PermanentAddressSection struct {
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	Country  string `json:"country"`
	State    string `json:"state"`
	District string `json:"district"`
	City     string `json:"city"`
	Pincode  string `json:"pincode"`
}

// ================================
// LANGUAGE DETAILS - Split into sections
// ================================

type LanguageDetailsResponse struct {
	EmployeeID          string                     `json:"employeeid"`
	LanguageProficiency LanguageProficiencySection `json:"language_proficiency"`
	HindiProficiency    HindiProficiencySection    `json:"hindi_proficiency"`
}

type LanguageProficiencySection struct {
	Language string `json:"language"`
	Read     string `json:"read"`
	Write    string `json:"write"`
	Speak    string `json:"speak"`
}

type HindiProficiencySection struct {
	LevelOfKnowledge string `json:"level_of_knowledge"`
	WorkingKnowledge string `json:"working_knowledge"`
	Proficiency      string `json:"proficiency"`
}

// ================================
// NOMINEE DETAILS - Split into sections
// ================================

type NomineeDetailsResponse struct {
	EmployeeID      string         `json:"employeeid"`
	GeneralNominee  NomineeSection `json:"general_nominee"`
	GratuityNominee NomineeSection `json:"gratuity_nominee"`
	GTISNominee     NomineeSection `json:"gtis_nominee"`
	NPSNominee      NomineeSection `json:"nps_nominee"`
	GPFNominee      NomineeSection `json:"gpf_nominee"`
}

type NomineeSection struct {
	Name            string `json:"name"`
	Relationship    string `json:"relationship"`
	DOB             string `json:"dob"`
	SharePercentage string `json:"share_percentage"`
}

// ================================
// OTHER CATEGORIES (unchanged)
// ================================

type EducationDetailsResponse struct {
	EmployeeID         string `json:"employeeid"`
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

type ExperienceDetailsResponse struct {
	EmployeeID       string `json:"employeeid"`
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

type DocumentDetailsResponse struct {
	EmployeeID   string `json:"employeeid"`
	DocumentName string `json:"document_name"`
}

type DependentDetailsResponse struct {
	EmployeeID                string `json:"employeeid"`
	DependentName             string `json:"dependent_name"`
	DependentRelationship     string `json:"dependent_relationship"`
	DependentDOB              string `json:"dependent_dob"`
	DependentAge              string `json:"dependent_age"`
	DependentMaritalStatus    string `json:"dependent_marital_status"`
	DependentBloodGroup       string `json:"dependent_blood_group"`
	DependentGender           string `json:"dependent_gender"`
	DependentEmploymentStatus string `json:"dependent_employment_status"`
	DependentAadhaarNo        string `json:"dependent_aadhaar_no"`
	IsTwins                   string `json:"is_twins"`
	DependentMobileNo         string `json:"dependent_mobile_no"`

	IsCurrentlyDependent        string `json:"is_currently_dependent"`
	OptingForInsurance          string `json:"opting_for_insurance"`
	OptingForLTC                string `json:"opting_for_ltc"`
	IsPersonDisabled            string `json:"is_person_disabled"`
	DependentNatureOfDisability string `json:"dependent_nature_of_disability"`
}

type HindiProficiencyResponse struct {
	EmployeeID            string `json:"employeeid"`
	HindiLevelOfKnowledge string `json:"hindi_level_of_knowledge"`
	HindiWorkingKnowledge string `json:"hindi_working_knowledge"`
	HindiProficiency      string `json:"hindi_proficiency"`
}

// scanCompleteEmployeeDetails scans all fields from the complete query
// This function handles NULL values by using COALESCE in SQL, so all fields will be empty strings if NULL
func scanCompleteEmployeeDetails(rows *sql.Rows, e *EmployeeEfileDetails) error {
	return rows.Scan(
		// ================================
		// PERSONAL DETAILS (27 fields)
		// ================================
		&e.EmployeeID,
		&e.FirstName,
		&e.MiddleName,
		&e.LastName,
		&e.Gender,
		&e.MaritalStatus,
		&e.FatherName,
		&e.MotherName,
		&e.SpouseName,
		&e.DOB,
		&e.Age,
		&e.Nationality,
		&e.BirthState,
		&e.BirthDistrict,
		&e.BirthPlace,
		&e.Hometown,
		&e.Religion,
		&e.CasteCategory,
		&e.EmergencyContactNo,
		&e.MobileNo,
		&e.IsPhysicallyChallenged,
		&e.PercentageOfDisability,
		&e.NatureOfDisability,
		&e.PersonalEmail,
		&e.AadhaarNo,
		&e.MotherTongue,
		&e.BankName,
		&e.IFSCCode,
		&e.BankAccountNo,
		&e.IdentificationMarks,
		&e.PanCardNo,

		// ================================
		// APPOINTMENT DETAILS (19 fields)
		// ================================
		&e.EmployeeName,
		&e.EmployeeType,
		&e.Department,
		&e.Designation,
		&e.Section,
		&e.RouteTo,
		&e.Grade,
		&e.EmployeeGroup,
		&e.PayInfo,
		&e.BasicPay,
		&e.NonPracticePay,
		&e.NameOfPayBand,
		&e.DateOfJoining,
		&e.DateOfConfirmation,
		&e.OfficeRoomNo,
		&e.OfficeExtensionNo,
		&e.IsActive,
		&e.EmployeeStatus,
		&e.DateOfRetirement,

		// ================================
		// EDUCATIONAL QUALIFICATION DETAILS (13 fields)
		// ================================
		&e.DegreeOrExam,
		&e.BoardName,
		&e.Institution,
		&e.UniversityName,
		&e.EducationCountry,
		&e.EducationState,
		&e.MonthYearOfPassing,
		&e.RegistrationNo,
		&e.Specialization,
		&e.Mode,
		&e.PercentageOfMarks,
		&e.ObtainedMarks,
		&e.Class,

		// ================================
		// EXPERIENCE DETAILS (9 fields)
		// ================================
		&e.OrganizationName,
		&e.Address1,
		&e.DesignationExp,
		&e.FromDate,
		&e.ToDate,
		&e.TotalExperience,
		&e.PayScale,
		&e.IsGovtEmployee,
		&e.TypeOfEmployment,

		// ================================
		// DEPENDENT INFO (17 fields)
		// ================================
		&e.DependentName,
		&e.DependentRelationship,
		&e.DependentDOB,
		&e.DependentAge,
		&e.DependentMaritalStatus,
		&e.DependentBloodGroup,
		&e.DependentGender,
		&e.DependentEmploymentStatus,
		&e.DependentAadhaarNo,
		&e.IsTwins,
		&e.DependentMobileNo,
		&e.IsCurrentlyDependent,
		&e.OptingForInsurance,
		&e.OptingForLTC,
		&e.IsPersonDisabled,
		&e.DependentNatureOfDisability,

		// ================================
		// NOMINEE INFO (20 fields)
		// ================================
		&e.GratuityNomineeName,
		&e.GratuityNomineeRelationship,
		&e.GratuityNomineeDOB,
		&e.GratuityNomineeSharePercentage,
		&e.GTISNomineeName,
		&e.GTISNomineeRelationship,
		&e.GTISNomineeDOB,
		&e.GTISNomineeSharePercentage,
		&e.NPSNomineeName,
		&e.NPSNomineeRelationship,
		&e.NPSNomineeDOB,
		&e.NPSNomineeSharePercentage,
		&e.GPFNomineeName,
		&e.GPFNomineeRelationship,
		&e.GPFNomineeDOB,
		&e.GPFNomineeSharePercentage,

		// ================================
		// CONTACT DETAILS (18 fields)
		// ================================
		&e.PresentAddress1,
		&e.PresentAddress2,
		&e.PresentCountry,
		&e.PresentState,
		&e.PresentDistrict,
		&e.PresentCity,
		&e.PresentPincode,
		&e.PresentCountryCode,
		&e.PresentAreaCode,
		&e.PermanentAddress1,
		&e.PermanentAddress2,
		&e.PermanentCountry,
		&e.PermanentState,
		&e.PermanentDistrict,
		&e.PermanentCity,
		&e.PermanentPincode,
		&e.IsQuarters,
		&e.QuartersNo,

		// ================================
		// LANGUAGE PROFICIENCY (4 fields)
		// ================================
		&e.Language,
		&e.Read,
		&e.Write,
		&e.Speak,

		// ================================
		// HINDI LANGUAGE PROFICIENCY (3 fields)
		// ================================
		&e.HindiLevelOfKnowledge,
		&e.HindiWorkingKnowledge,
		&e.HindiProficiency,

		// ================================
		// DOCUMENTS (1 field)
		// ================================
		&e.DocumentName,
	)

}
