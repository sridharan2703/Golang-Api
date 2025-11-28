// Package controllersofficeorder contains structs and handles insertion into the OfficeOrder table, along with next task assignment and rejection logic.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 01-11-2025
//
// Last Modified By:  Sridharan
//
// Last Modified Date: 21-11-2025
package controllersofficeorder

import (
	"Hrmodule/auth"
	credentials "Hrmodule/dbconfig"
	"Hrmodule/utils"

	// "bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/tls"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Global variable to hold the encryption key as []byte.
var secretKeyde []byte

// ======================
// Initialization (init)
// ======================
func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file. Ensure it is in the application's execution directory.")
	}

	key := os.Getenv("ENCRYPTION_KEY")
	if key == "" {
		panic("ENCRYPTION_KEY not set in .env file")
	}
	if len(key) != 32 {
		panic(fmt.Sprintf("ENCRYPTION_KEY must be 32 bytes for AES-256 GCM, got %d bytes", len(key)))
	}

	secretKeyde = []byte(key)
}

// ======================
// Custom Types for Nullable Fields
// ======================
type NullTime struct {
	sql.NullTime
}

func (nt *NullTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		nt.Valid = false
		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	if str == "" {
		nt.Valid = false
		return nil
	}

	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, str); err == nil {
			nt.Time = t
			nt.Valid = true
			return nil
		}
	}

	return fmt.Errorf("invalid date format: %s", str)
}

type NullInt struct {
	sql.NullInt64
}

func (ni *NullInt) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ni.Valid = false
		return nil
	}

	var val int64
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}

	ni.Int64 = val
	ni.Valid = true
	return nil
}

type NullString struct {
	sql.NullString
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		ns.Valid = false
		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	if str == "" {
		ns.Valid = false
		return nil
	}

	ns.String = str
	ns.Valid = true
	return nil
}

// ======================
// Request Wrapper for Encrypted Data
// ======================
type EncryptedRequest struct {
	Data string `json:"Data"`
}

// ======================
// Request & Response Structs
// ======================
type PCRInsertRequest struct {
	// Authentication & Session
	Token     string `json:"token"`
	PID       string `json:"P_id"` // Added for encryption/decryption
	SessionID string `json:"session_id"`

	// Submission Control
	TypeOfSubmit string `json:"typeofsubmit"` // "draft" or "submit"

	// PCR Basic Information
	CoverPageNo   string   `json:"p_cover_page_no"`
	EmployeeID    string   `json:"p_employee_id"`
	EmployeeName  string   `json:"p_employee_name"`
	Department    string   `json:"p_department"`
	Designation   string   `json:"p_designation"`
	VisitFrom     NullTime `json:"p_visit_from"`
	VisitTo       NullTime `json:"p_visit_to"`
	Duration      int      `json:"p_duration"`
	NatureOfVisit string   `json:"p_nature_of_visit"`
	ClaimType     string   `json:"p_claim_type"`
	CityTown      string   `json:"p_city_town"`
	Country       string   `json:"p_country"`

	// Order Document Fields
	HeaderHTML    string   `json:"p_header_html"`
	OrderNo       string   `json:"p_order_no"`
	OrderDate     NullTime `json:"p_order_date"`
	ToColumn      string   `json:"p_to_column"`
	Subject       string   `json:"p_subject"`
	Reference     string   `json:"p_reference"`
	BodyHTML      string   `json:"p_body_html"`
	SignatureHTML string   `json:"p_signature_html"`
	CCTo          string   `json:"p_cc_to"`
	FooterHTML    string   `json:"p_footer_html"`

	// Workflow Assignment Fields
	AssignTo     string `json:"p_assign_to"`
	AssignedRole string `json:"p_assigned_role"`

	// Task Status & Activity
	TaskStatusID         int  `json:"p_task_status_id"`
	CurrentActivitySeqNo int  `json:"p_current_activity_seq_no"`
	ActivitySeqNo        int  `json:"p_activity_seq_no"`
	IsTaskReturn         bool `json:"p_is_task_return"`
	IsTaskApproved       bool `json:"p_is_task_approved"`

	// Audit Fields
	InitiatedBy string   `json:"p_initiated_by"`
	InitiatedOn NullTime `json:"p_initiated_on"`
	UpdatedBy   string   `json:"p_updated_by"`
	UpdatedOn   NullTime `json:"p_updated_on"`

	// Email & Template
	EmailFlag  bool    `json:"p_email_flag"`
	TemplateID NullInt `json:"p_template_id"`

	// Rejection Handling
	RejectFlag int    `json:"p_reject_flag"`
	RejectRole string `json:"p_reject_role"`

	// Additional Order Info
	OriginalOrderNo string `json:"p_original_order_no"`
	OrderType       string `json:"p_order_type"`
	Remarks         string `json:"p_remarks"`
	UserRole        string `json:"p_user_role"`

	// Process Configuration
	ProcessID int    `json:"p_process_id"`
	Priority  string `json:"p_priority"`

	// NextActivity API Specific Fields
	TaskID         NullString `json:"p_task_id"`
	Role           NullString `json:"p_role"`
	RequestedUser  NullString `json:"p_requested_user"`
	ReturnToRole   NullString `json:"p_return_to_role"`
	ReturnToUser   NullString `json:"p_return_to_user"`
	SendBackToMe   *bool      `json:"p_send_back_to_me"`
	SendBackToUser NullString `json:"p_send_back_to_user"`
}

// NextActivity API Request/Response structs
type NextActivityRequest struct {
	Token           string            `json:"token"`
	ProcessID       int               `json:"ProcessId"`
	CurrentActivity int               `json:"Current_Activity"`
	TaskID          *string           `json:"TaskId"`
	IsApproved      int               `json:"IsApproved"`
	IsReturn        int               `json:"IsReturn"`
	Role            *string           `json:"Role"`
	EmployeeID      *string           `json:"EmployeeId"`
	Conditions      map[string]string `json:"Conditions"`
	RequestedUser   *string           `json:"RequestedUser"`
	ReturnToRole    *string           `json:"ReturnToRole"`
	ReturnToUser    *string           `json:"ReturnToUser"`
	SendBackToMe    bool              `json:"SendBackToMe"`
	SendBackToUser  *string           `json:"SendBackToUser"`
}

type NextActivityResponse struct {
	Status  interface{} `json:"Status"` // Can be string or number
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

// Helper method to get Status as string
func (r *NextActivityResponse) GetStatusString() string {
	switch v := r.Status.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.Itoa(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

type NextActivityData struct {
	NoOfRecords int                  `json:"No Of Records"`
	Records     []NextActivityRecord `json:"Records"`
}

type NextActivityRecord struct {
	NextActivities string  `json:"next_activities"`
	RoleNames      string  `json:"role_names"`
	EmailFlag      int     `json:"emailflag"`
	TemplateID     *int    `json:"template_id"`
	AssignTo       *string `json:"assign_to"`
	RejectFlag     int     `json:"reject_flag"`
}

// Email Queue API Request/Response structs
type EmailQueueRequest struct {
	ProcessID  int    `json:"ProcessId"`
	TaskID     string `json:"TaskId"`
	TemplateID int    `json:"TemplateId"`
	Token      string `json:"token"`
}

type EmailQueueResponse struct {
	Status  interface{} `json:"Status"` // Can be string or number
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

// Helper method to get Status as string
func (r *EmailQueueResponse) GetStatusString() string {
	switch v := r.Status.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.Itoa(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

type APIResponsePCRInsert struct {
	Status  int         `json:"Status"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

// ======================
// HTML Helper Functions
// ======================

// stripHTMLTags removes all HTML tags from a string
func stripHTMLTags(html string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	text := re.ReplaceAllString(html, " ")

	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	text = strings.ReplaceAll(text, "&#39;", "'")

	text = strings.TrimSpace(text)
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	return text
}

// ======================
// Helper Functions
// ======================

func (req *PCRInsertRequest) validate() error {

	if req.TypeOfSubmit == "" {
		return fmt.Errorf("typeofsubmit is required")
	}
	typeOfSubmit := strings.ToLower(req.TypeOfSubmit)
	if typeOfSubmit != "draft" && typeOfSubmit != "submit" {
		return fmt.Errorf("typeofsubmit must be either 'draft' or 'submit'")
	}
	if req.CoverPageNo == "" {
		return fmt.Errorf("p_cover_page_no is required")
	}
	if req.EmployeeID == "" {
		return fmt.Errorf("p_employee_id is required")
	}
	if req.EmployeeName == "" {
		return fmt.Errorf("p_employee_name is required")
	}
	if req.ProcessID == 0 {
		return fmt.Errorf("p_process_id is required and cannot be 0 or null")
	}
	return nil
}

func (req *PCRInsertRequest) validatePCRData() error {
	if strings.ToLower(req.TypeOfSubmit) == "draft" {
		return nil
	}

	if req.Department == "" {
		return fmt.Errorf("p_department is required")
	}
	if req.Designation == "" {
		return fmt.Errorf("p_designation is required")
	}
	if req.VisitFrom.Valid && req.VisitTo.Valid && req.VisitFrom.Time.After(req.VisitTo.Time) {
		return fmt.Errorf("p_visit_from cannot be after p_visit_to")
	}
	if req.Duration < 0 {
		return fmt.Errorf("p_duration cannot be negative")
	}
	if req.NatureOfVisit == "" {
		return fmt.Errorf("p_nature_of_visit is required")
	}
	if len(req.NatureOfVisit) > 200 {
		return fmt.Errorf("p_nature_of_visit exceeds 200 characters")
	}
	if req.CityTown == "" {
		return fmt.Errorf("p_city_town is required")
	}
	if len(req.CityTown) > 200 {
		return fmt.Errorf("p_city_town exceeds 200 characters")
	}
	if req.Country == "" {
		return fmt.Errorf("p_country is required")
	}
	if len(req.Country) > 200 {
		return fmt.Errorf("p_country exceeds 200 characters")
	}
	if req.ToColumn == "" {
		return fmt.Errorf("p_to_column is required")
	}
	if req.Subject == "" {
		return fmt.Errorf("p_subject is required")
	}
	if req.InitiatedBy == "" {
		return fmt.Errorf("p_initiated_by is required")
	}
	if req.UpdatedBy == "" {
		return fmt.Errorf("p_updated_by is required")
	}
	if req.ProcessID <= 0 {
		return fmt.Errorf("p_process_id must be greater than 0")
	}

	if req.ActivitySeqNo <= 0 {
		return fmt.Errorf("p_activity_seq_no is required for submit")
	}

	return nil
}

func toSQLValue(nt NullTime) interface{} {
	if nt.Valid {
		return nt.Time
	}
	return nil
}

func toSQLInt(ni NullInt) interface{} {
	if ni.Valid {
		return ni.Int64
	}
	return nil
}

func toSQLString(s string) interface{} {
	if s != "" {
		return s
	}
	return nil
}

func toNullString(ns NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

// ======================
// Encryption/Decryption Helper Functions
// ======================

// encryptAndRespond encrypts the response and sends it back
func encryptAndRespond(w http.ResponseWriter, r *http.Request, pid, key string, status int, payload interface{}) {
	responseJSON, _ := json.MarshalIndent(payload, "", "  ")
	encryptedResponse := string(responseJSON)

	if key != "" {
		if enc, err := utils.EncryptAES(string(responseJSON), key); err == nil {
			encryptedResponse = fmt.Sprintf("%s||%s", pid, enc)
		}
	}

	finalResp := map[string]string{"Data": encryptedResponse}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(finalResp)
}

// ======================
// Main Handler with Encryption/Decryption
// ======================

func PCRInsert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	// Step 1: Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Step 2: Parse the outer JSON wrapper
	var encReq EncryptedRequest
	if err := json.Unmarshal(body, &encReq); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Step 3: Split PID and encrypted data
	parts := strings.Split(encReq.Data, "||")
	if len(parts) != 2 {
		http.Error(w, "Invalid Data format", http.StatusBadRequest)
		return
	}
	pid := parts[0]
	encryptedPart := parts[1]

	// Step 4: Get decryption key using PID
	key, err := utils.GetDecryptKey(pid)
	if err != nil {
		http.Error(w, "Failed to get decryption key", http.StatusUnauthorized)
		return
	}

	// Step 5: Decrypt the data
	decryptedJSON, err := utils.DecryptAES(encryptedPart, key)
	if err != nil {
		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}

	// Step 6: Parse decrypted JSON into PCRInsertRequest
	var req PCRInsertRequest
	if err := json.Unmarshal([]byte(decryptedJSON), &req); err != nil {
		http.Error(w, "Invalid decrypted data", http.StatusBadRequest)
		return
	}

	// Store PID for later use in response
	req.PID = pid

	// Step 7: Validate required fields
	if err := req.validate(); err != nil {
		encryptAndRespond(w, r, pid, key, http.StatusBadRequest,
			map[string]string{"error": err.Error()})
		return
	}

	if err := req.validatePCRData(); err != nil {
		encryptAndRespond(w, r, pid, key, http.StatusBadRequest,
			map[string]string{"error": "Data validation failed: " + err.Error()})
		return
	}

	// Step 8: Token validation
	if req.Token != "" {
		r.Header.Set("token", req.Token)
	}

	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	// Step 9: Process the request with auth logging
	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := auth.IsValidIDFromRequest(r); err != nil {
			encryptAndRespond(w, r, pid, key, http.StatusBadRequest,
				map[string]string{"error": "Invalid TOKEN provided"})
			return
		}

		// Determine if this is a submit action
		isSubmit := strings.ToLower(req.TypeOfSubmit) == "submit"

		fmt.Printf("=== BEFORE NextActivity API ===\n")
		fmt.Printf("ActivitySeqNo: %d\n", req.ActivitySeqNo)
		fmt.Printf("CurrentActivitySeqNo: %d\n", req.CurrentActivitySeqNo)
		fmt.Printf("AssignedRole: %s\n", req.AssignedRole)
		fmt.Printf("TaskStatusID: %d\n", req.TaskStatusID)
		fmt.Printf("EmailFlag: %v\n", req.EmailFlag)
		fmt.Printf("TemplateID: %v\n", req.TemplateID)
		fmt.Printf("RejectFlag: %d\n", req.RejectFlag)

		// STEP 1: Call NextActivity API ONLY if typeofsubmit is "submit"
		if isSubmit {
			fmt.Printf("\n🔄 TypeOfSubmit is 'submit' - Calling NextActivity API...\n")

			_, err := callNextActivityForPCR(&req, pid, key)
			if err != nil {
				encryptAndRespond(w, r, pid, key, http.StatusInternalServerError,
					map[string]string{"error": fmt.Sprintf("NextActivity API error: %v", err)})
				return
			}

			fmt.Printf("\n=== AFTER NextActivity API ===\n")
			fmt.Printf("ActivitySeqNo: %d (UPDATED)\n", req.ActivitySeqNo)
			fmt.Printf("CurrentActivitySeqNo: %d\n", req.CurrentActivitySeqNo)
			fmt.Printf("AssignedRole: %s (UPDATED)\n", req.AssignedRole)
			fmt.Printf("AssignTo: %s (UPDATED)\n", req.AssignTo)
			fmt.Printf("EmailFlag: %v (UPDATED)\n", req.EmailFlag)
			fmt.Printf("TemplateID: %v (UPDATED)\n", req.TemplateID)
			fmt.Printf("RejectFlag: %d (UPDATED)\n", req.RejectFlag)
		} else {
			fmt.Printf("\n📝 TypeOfSubmit is 'draft' - Skipping NextActivity API\n")
		}

		// STEP 2: Single Database Insert/Update
		fmt.Printf("\n💾 Calling Database Stored Procedure (SINGLE CALL)...\n")
		if err := executePCRInsert(&req); err != nil {
			encryptAndRespond(w, r, pid, key, http.StatusInternalServerError,
				map[string]string{"error": fmt.Sprintf("Database operation error: %v", err)})
			return
		}

		if isSubmit {
			fmt.Printf("✅ PCR record inserted/updated with workflow data for cover_page_no: %s\n", req.CoverPageNo)
			fmt.Printf("   Final DB Values: ActivitySeqNo=%d, AssignedRole=%s, TaskStatusID=%d\n",
				req.ActivitySeqNo, req.AssignedRole, req.TaskStatusID)
		} else {
			fmt.Printf("✅ PCR record saved as draft for cover_page_no: %s\n", req.CoverPageNo)
		}

		// STEP 2.5: Trigger Office Order PDF ONLY when next_activities = 6
		if isSubmit && req.ActivitySeqNo == 6 {
			fmt.Printf("📝 next_activities = 6 → triggering Office Order PDF API...\n")

			taskID, err := getTaskIDFromDBorderno(req.OrderNo)
			if err != nil {
				fmt.Printf("⚠️ Warning: Unable to fetch task_id for Office Order PDF API: %v\n", err)
			} else {
				if err := callOfficeOrderPDF(&req, taskID, pid, key); err != nil {
					fmt.Printf("⚠️ Warning: Office Order PDF API call failed: %v\n", err)
				} else {
					fmt.Printf("✅ Office Order PDF API successfully triggered for task_id: %s\n", taskID)
				}
			}
		}

		// STEP 3: Call Email Queue API if it's a submit AND EmailFlag is true
		if isSubmit && req.EmailFlag {
			taskID, err := getTaskIDFromDB(req.CoverPageNo)
			if err != nil {
				fmt.Printf("⚠️  Warning: Failed to get task_id for email queue: %v\n", err)
			} else {
				if err := callEmailQueueAPI(&req, taskID, pid, key); err != nil {
					fmt.Printf("⚠️  Warning: Email queue API call failed: %v\n", err)
				} else {
					fmt.Printf("✅ Email queue API called successfully for task_id: %s\n", taskID)
				}
			}
		}

		// Prepare success message
		successMsg := "PCR record saved as draft successfully"
		if isSubmit {
			successMsg = "PCR record submitted and updated with workflow data successfully"
		}

		// Send encrypted success response
		response := APIResponsePCRInsert{
			Status:  200,
			Message: successMsg,
			Data:    nil,
		}

		encryptAndRespond(w, r, pid, key, http.StatusOK, response)

	})).ServeHTTP(w, r)
}

// ======================
// Get Task ID from Database
// ======================
func getTaskIDFromDB(coverPageNo string) (string, error) {
	db, err := sql.Open("postgres", credentials.Getdatabasemeivan())
	if err != nil {
		return "", fmt.Errorf("DB connection failed: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return "", fmt.Errorf("DB ping failed: %w", err)
	}

	var taskID string
	query := "SELECT task_id FROM meivan.pcr_m WHERE cover_page_no = $1"

	err = db.QueryRow(query, coverPageNo).Scan(&taskID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no task_id found for cover_page_no: %s", coverPageNo)
		}
		return "", fmt.Errorf("failed to query task_id: %w", err)
	}

	return taskID, nil
}

func getTaskIDFromDBorderno(orderno string) (string, error) {
	db, err := sql.Open("postgres", credentials.Getdatabasemeivan())
	if err != nil {
		return "", fmt.Errorf("DB connection failed: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return "", fmt.Errorf("DB ping failed: %w", err)
	}

	var taskID string
	query := "SELECT task_id FROM meivan.pcr_m WHERE order_no = $1"

	err = db.QueryRow(query, orderno).Scan(&taskID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no task_id found for order_no: %s", orderno)
		}
		return "", fmt.Errorf("failed to query task_id: %w", err)
	}

	return taskID, nil
}

// ======================
// Office Order PDF API Call with Encryption/Decryption
// ======================
func callOfficeOrderPDF(req *PCRInsertRequest, taskID, pid, key string) error {

	// Build JSON request (similar to EmailQueue style)
	officeReq := map[string]interface{}{
		"process_id": req.ProcessID,
		"token":      req.Token,
		"P_id":       pid,
		"task_id":    taskID,
		"status":     "completed",
	}

	jsonData, err := json.Marshal(officeReq)
	if err != nil {
		return fmt.Errorf("failed to marshal office order payload: %w", err)
	}

	// Encrypt JSON data
	encryptedData, err := utils.EncryptAES(string(jsonData), key)
	if err != nil {
		return fmt.Errorf("failed to encrypt office order request: %w", err)
	}

	// Wrap encrypted request as PID||EncryptedString
	encryptedReq := EncryptedRequest{
		Data: fmt.Sprintf("%s||%s", pid, encryptedData),
	}

	encryptedReqJSON, err := json.Marshal(encryptedReq)
	if err != nil {
		return fmt.Errorf("failed to marshal encrypted office request: %w", err)
	}

	fmt.Printf("=== Office Order PDF API Request ===\n")
	fmt.Printf("URL: https://hr-api-uat.iitm.ac.in:8081/pdfapi\n")
	fmt.Printf("Encrypted Request: %s\n", string(encryptedReqJSON))

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Send API request
	resp, err := client.Post(
		"https://hr-api-uat.iitm.ac.in:8081/pdfapi",
		"application/json",
		strings.NewReader(string(encryptedReqJSON)),
	)
	if err != nil {
		return fmt.Errorf("office order API request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read office order response: %w", err)
	}

	fmt.Printf("Office Order Encrypted Response: %s\n", string(respBody))

	// Parse encrypted response wrapper
	var encResp EncryptedRequest
	if err := json.Unmarshal(respBody, &encResp); err != nil {
		return fmt.Errorf("failed to unmarshal encrypted response: %w", err)
	}

	// Split PID||encryptedData
	parts := strings.Split(encResp.Data, "||")
	if len(parts) != 2 {
		return fmt.Errorf("invalid encrypted response format")
	}

	encryptedPart := parts[1]

	// Decrypt final JSON response
	decryptedJSON, err := utils.DecryptAES(encryptedPart, key)
	if err != nil {
		return fmt.Errorf("failed to decrypt office order response: %w", err)
	}

	fmt.Printf("Office Order Decrypted Response: %s\n", decryptedJSON)

	var officeResp struct {
		Status  interface{} `json:"status"`
		Message string      `json:"message"`
		Path    string      `json:"path"`
	}

	if err := json.Unmarshal([]byte(decryptedJSON), &officeResp); err != nil {
		return fmt.Errorf("failed to unmarshal decrypted JSON: %w", err)
	}

	// Convert status number/string to string
	statusStr := fmt.Sprintf("%v", officeResp.Status)

	if statusStr != "200" && statusStr != "HRM-BUS-SUC-0001" {
		return fmt.Errorf("office order API error %s: %s", statusStr, officeResp.Message)
	}

	fmt.Printf("📄 PDF Generated: %s\n", officeResp.Path)

	return nil
}

// ======================
// Email Queue API Call with Encryption/Decryption
// ======================
func callEmailQueueAPI(req *PCRInsertRequest, taskID, pid, key string) error {
	if !req.TemplateID.Valid {
		return fmt.Errorf("template_id is not valid")
	}

	emailReq := EmailQueueRequest{
		ProcessID:  req.ProcessID,
		TaskID:     taskID,
		TemplateID: int(req.TemplateID.Int64),
		Token:      req.Token,
	}

	jsonData, err := json.Marshal(emailReq)
	if err != nil {
		return fmt.Errorf("failed to marshal email queue request: %w", err)
	}

	// Encrypt the request
	encryptedData, err := utils.EncryptAES(string(jsonData), key)
	if err != nil {
		return fmt.Errorf("failed to encrypt email queue request: %w", err)
	}

	// Create encrypted request wrapper
	encryptedReq := EncryptedRequest{
		Data: fmt.Sprintf("%s||%s", pid, encryptedData),
	}

	encryptedReqJSON, err := json.Marshal(encryptedReq)
	if err != nil {
		return fmt.Errorf("failed to marshal encrypted email queue request: %w", err)
	}

	fmt.Printf("=== Email Queue API Request ===\n")
	fmt.Printf("URL: https://hr-api-uat.iitm.ac.in:5555/Emailqueue\n")
	fmt.Printf("Encrypted Request Body: %s\n", string(encryptedReqJSON))

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: tr,
	}

	resp, err := client.Post(
		"https://hr-api-uat.iitm.ac.in:5555/Emailqueue",
		"application/json",
		strings.NewReader(string(encryptedReqJSON)),
	)
	if err != nil {
		return fmt.Errorf("email queue API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read email queue response: %w", err)
	}

	fmt.Printf("Email Queue API Encrypted Response: %s\n", string(body))

	// Parse encrypted response
	var encResp EncryptedRequest
	if err := json.Unmarshal(body, &encResp); err != nil {
		return fmt.Errorf("failed to unmarshal encrypted email queue response: %w", err)
	}

	// Split PID and encrypted part
	parts := strings.Split(encResp.Data, "||")
	if len(parts) != 2 {
		return fmt.Errorf("invalid encrypted email queue response format")
	}
	encryptedPart := parts[1]

	// Decrypt response
	decryptedJSON, err := utils.DecryptAES(encryptedPart, key)
	if err != nil {
		return fmt.Errorf("failed to decrypt email queue response: %w", err)
	}

	fmt.Printf("Email Queue API Decrypted Response: %s\n", decryptedJSON)

	var emailResp EmailQueueResponse
	if err := json.Unmarshal([]byte(decryptedJSON), &emailResp); err != nil {
		return fmt.Errorf("failed to unmarshal email queue response: %w", err)
	}

	statusStr := emailResp.GetStatusString()

	// Check for success status - can be 200 (number) or "HRM-BUS-SUC-0001" (string)
	if statusStr != "200" && statusStr != "HRM-BUS-SUC-0001" {
		return fmt.Errorf("email queue API failed with status %s: %s", statusStr, emailResp.Message)
	}

	fmt.Printf("✅ Email Queue API Success with Status: %s\n", statusStr)

	return nil
}

// ======================
// NextActivity API Functions with Encryption/Decryption
// ======================

func decryptNextActivityGCM(encryptedData string) ([]byte, error) {
	if len(secretKeyde) == 0 {
		return nil, fmt.Errorf("encryption key not initialized. Check init() for .env loading errors")
	}

	data, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("base64 decode failed: %w", err)
	}

	nonceSize := 12
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short, expected at least %d bytes for nonce", nonceSize)
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	keyBytes := secretKeyde
	if len(keyBytes) != 32 {
		return nil, fmt.Errorf("encryption key must be 32 bytes for AES-256 GCM")
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("cipher creation failed: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("GCM cipher failed: %w", err)
	}

	plainText, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("GCM decryption failed (Auth tag mismatch or corruption): %w", err)
	}

	return plainText, nil
}

func callNextActivityForPCR(req *PCRInsertRequest, pid, key string) (interface{}, error) {
	choiceValue := req.SignatureHTML
	if choiceValue == "" {
		return nil, fmt.Errorf("p_signature_html is required for NextActivity API call")
	}

	choiceValue = stripHTMLTags(choiceValue)
	fmt.Printf("🔍 Using choice value (from p_signature_html): '%s'\n", choiceValue)

	currentActivity := req.CurrentActivitySeqNo
	if currentActivity <= 0 {
		currentActivity = req.ActivitySeqNo
		if currentActivity <= 0 {
			return nil, fmt.Errorf("either p_current_activity_seq_no or p_activity_seq_no must be greater than 0")
		}
		req.CurrentActivitySeqNo = currentActivity
		fmt.Printf("ℹ️  p_current_activity_seq_no not provided, using p_activity_seq_no=%d\n", currentActivity)
	}

	isApproved := 0
	if req.IsTaskApproved {
		isApproved = 1
	}

	isReturn := 0
	if req.IsTaskReturn {
		isReturn = 1
	}

	sendBackToMe := false
	if req.SendBackToMe != nil {
		sendBackToMe = *req.SendBackToMe
	}

	// Construct NextActivity API request
	nextActivityReq := NextActivityRequest{
		Token:           req.Token,
		ProcessID:       req.ProcessID,
		CurrentActivity: currentActivity,
		TaskID:          toNullString(req.TaskID),
		IsApproved:      isApproved,
		IsReturn:        isReturn,
		Role:            toNullString(req.Role),
		EmployeeID:      nil,
		Conditions: map[string]string{
			"choice": choiceValue,
		},
		RequestedUser:  toNullString(req.RequestedUser),
		ReturnToRole:   toNullString(req.ReturnToRole),
		ReturnToUser:   toNullString(req.ReturnToUser),
		SendBackToMe:   sendBackToMe,
		SendBackToUser: toNullString(req.SendBackToUser),
	}

	fmt.Printf("=== NextActivity API Parameters ===\n")
	fmt.Printf("CurrentActivity: %d (from p_current_activity_seq_no)\n", currentActivity)
	fmt.Printf("TaskID: %v\n", nextActivityReq.TaskID)
	fmt.Printf("IsApproved: %d (from IsTaskApproved: %v)\n", nextActivityReq.IsApproved, req.IsTaskApproved)
	fmt.Printf("IsReturn: %d (from IsTaskReturn: %v)\n", nextActivityReq.IsReturn, req.IsTaskReturn)
	fmt.Printf("Role: %v\n", nextActivityReq.Role)
	fmt.Printf("Choice (from SignatureHTML): %s\n", choiceValue)

	jsonData, err := json.Marshal(nextActivityReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal NextActivity request: %w", err)
	}

	// Encrypt the request
	encryptedData, err := utils.EncryptAES(string(jsonData), key)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt NextActivity request: %w", err)
	}

	// Create encrypted request wrapper
	encryptedReq := EncryptedRequest{
		Data: fmt.Sprintf("%s||%s", pid, encryptedData),
	}

	encryptedReqJSON, err := json.Marshal(encryptedReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal encrypted NextActivity request: %w", err)
	}

	fmt.Printf("=== NextActivity API Request ===\n")
	fmt.Printf("URL: https://hr-api-uat.iitm.ac.in:5555/NextActivity\n")
	fmt.Printf("Encrypted Request Body: %s\n", string(encryptedReqJSON))

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: tr,
	}

	resp, err := client.Post(
		"https://hr-api-uat.iitm.ac.in:5555/NextActivity",
		"application/json",
		strings.NewReader(string(encryptedReqJSON)),
	)
	if err != nil {
		return nil, fmt.Errorf("NextActivity API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read NextActivity response: %w", err)
	}

	fmt.Printf("NextActivity API Encrypted Response: %s\n", string(body))

	// Parse encrypted response
	var encResp EncryptedRequest
	if err := json.Unmarshal(body, &encResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal encrypted NextActivity response: %w", err)
	}

	// Split PID and encrypted part
	parts := strings.Split(encResp.Data, "||")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid encrypted NextActivity response format")
	}
	encryptedPart := parts[1]

	// Decrypt response
	decryptedJSON, err := utils.DecryptAES(encryptedPart, key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt NextActivity response: %w", err)
	}

	fmt.Printf("NextActivity API Decrypted Response: %s\n", decryptedJSON)

	var nextActivityResp NextActivityResponse
	if err := json.Unmarshal([]byte(decryptedJSON), &nextActivityResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal decrypted NextActivity response: %w", err)
	}

	// Check for success status - can be 200 (number) or "HRM-BUS-SUC-0001" (string)
	statusStr := nextActivityResp.GetStatusString()

	if statusStr != "200" && statusStr != "HRM-BUS-SUC-0001" {
		return nil, fmt.Errorf("NextActivity API failed with status %s: %s", statusStr, nextActivityResp.Message)
	}

	fmt.Printf("✅ NextActivity API Success with Status: %s\n", statusStr)

	dataBytes, err := json.Marshal(nextActivityResp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal NextActivity data for specific unmarshaling: %w", err)
	}

	var workflowData NextActivityData
	if err := json.Unmarshal(dataBytes, &workflowData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal workflow data: %w", err)
	}

	if len(workflowData.Records) == 0 {
		return nil, fmt.Errorf("NextActivity API returned success but no records were found")
	}

	record := workflowData.Records[0]

	if record.NextActivities == "" {
		return nil, fmt.Errorf("NextActivity API returned empty next_activities - workflow configuration may be missing for ProcessID=%d, CurrentActivity=%d, Choice='%s'. Please verify workflow setup",
			req.ProcessID, req.ActivitySeqNo, choiceValue)
	}

	nextActivity, err := strconv.Atoi(record.NextActivities)
	if err != nil {
		return nil, fmt.Errorf("failed to convert next_activities '%s' to int: %w (ProcessID=%d, CurrentActivity=%d, Choice='%s')",
			record.NextActivities, err, req.ProcessID, req.ActivitySeqNo, choiceValue)
	}
	req.ActivitySeqNo = nextActivity

	req.AssignedRole = record.RoleNames

	req.EmailFlag = record.EmailFlag == 1

	if record.TemplateID != nil {
		req.TemplateID = NullInt{
			sql.NullInt64{
				Int64: int64(*record.TemplateID),
				Valid: true,
			},
		}
	} else {
		req.TemplateID = NullInt{
			sql.NullInt64{
				Valid: false,
			},
		}
	}

	if record.AssignTo != nil {
		req.AssignTo = *record.AssignTo
	} else {
		req.AssignTo = ""
	}

	req.RejectFlag = record.RejectFlag

	fmt.Printf("⭐️ PCRInsertRequest Updated with NextActivity Data:\n")
	fmt.Printf("   ActivitySeqNo=%d\n", req.ActivitySeqNo)
	fmt.Printf("   AssignedRole=%s\n", req.AssignedRole)
	fmt.Printf("   EmailFlag=%v\n", req.EmailFlag)
	fmt.Printf("   TemplateID=%v\n", req.TemplateID)
	fmt.Printf("   AssignTo=%s\n", req.AssignTo)
	fmt.Printf("   RejectFlag=%d\n", req.RejectFlag)

	return nextActivityResp.Data, nil
}

// ======================
// Database Execution
// ======================
func executePCRInsert(req *PCRInsertRequest) error {
	db, err := sql.Open("postgres", credentials.Getdatabasemeivan())
	if err != nil {
		return fmt.Errorf("DB connection failed: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("DB ping failed: %w", err)
	}

	query := `CALL meivan.pcr_insert_amen_can(
		$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,
		$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,
		$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,
		$31,$32,$33,$34,$35,$36,$37,$38,$39,$40,$41,$42,$43
	)`

	_, err = db.Exec(query,
		req.CoverPageNo,
		req.EmployeeID,
		req.EmployeeName,
		req.Department,
		req.Designation,
		toSQLValue(req.VisitFrom),
		toSQLValue(req.VisitTo),
		req.Duration,
		req.NatureOfVisit,
		req.ClaimType,
		req.CityTown,
		req.Country,
		req.HeaderHTML,
		req.OrderNo,
		toSQLValue(req.OrderDate),
		req.ToColumn,
		req.Subject,
		req.Reference,
		req.BodyHTML,
		req.SignatureHTML,
		req.CCTo,
		req.FooterHTML,
		req.AssignTo,
		req.AssignedRole,
		req.TaskStatusID,
		req.ActivitySeqNo,
		req.IsTaskReturn,
		req.IsTaskApproved,
		req.InitiatedBy,
		toSQLValue(req.InitiatedOn),
		req.UpdatedBy,
		toSQLValue(req.UpdatedOn),
		req.EmailFlag,
		toSQLInt(req.TemplateID),
		req.RejectFlag,
		toSQLString(req.RejectRole),
		toSQLString(req.OriginalOrderNo),
		toSQLString(req.OrderType),
		toSQLString(req.Remarks),
		toSQLString(req.UserRole),
		req.ProcessID,
		req.CurrentActivitySeqNo,
		req.Priority,
	)
	if err != nil {
		return fmt.Errorf("stored procedure execution failed: %w", err)
	}

	return nil
}
