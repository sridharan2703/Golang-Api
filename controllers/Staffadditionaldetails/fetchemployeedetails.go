// Package controllerssad contains structs and queries for Employee E-File Details API.
// --- Creator's Info ---
// Creator: Rovita
// Created On: 05-11-2025
// Last Modified By:
// Last Modified Date:
// This API fetches employee E-File details based on specified category

package controllerssad

import (
	"Hrmodule/auth"
	databasesad "Hrmodule/database/Staffadditionaldetails"
	"Hrmodule/utils"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

// APIResponseEmployeeEfile defines the structure of the JSON response
type APIResponseEmployeeEfile struct {
	Status  int         `json:"Status"`
	Message string      `json:"message"`
	Data    interface{} `json:"Data"`
}

// EmployeeEfileRequest defines the expected JSON request body
type EmployeeEfileRequest struct {
	Token      string `json:"token"`
	SessionID  string `json:"session_id"`
	EmployeeID string `json:"employee_id"`
	Category   string `json:"category"`
}

// Valid categories - expanded to include all categories from the SQL query
var validCategories = map[string]bool{
	"personaldetails":    true,
	"appointmentdetails": true,
	"educationdetails":   true,
	"experiencedetails":  true,
	"languagedetails":    true,
	"documentdetails":    true,
	"dependentdetails":   true,
	"nomineedetails":     true,
	"contactdetails":     true,
	"hindiproficiency":   true,
}

// encryptedResponse handles error responses by encrypting the message
func encryptedResponse(w http.ResponseWriter, statusCode int, msg string) {
	response := APIResponseEmployeeEfile{
		Status:  statusCode,
		Message: msg,
		Data:    nil,
	}

	jsonResp, _ := json.Marshal(response)
	encrypted, err := utils.Encryptnew(jsonResp)
	if err != nil {
		http.Error(w, "Encryption failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]string{"Data": encrypted})
}

// validateCategory checks if the provided category is valid
func validateCategory(category string) bool {
	if category == "" {
		return false
	}
	// Convert to lowercase for case-insensitive comparison
	category = strings.ToLower(strings.TrimSpace(category))
	return validCategories[category]
}

// getValidCategoriesString returns a string of valid categories for error messages
func getValidCategoriesString() string {
	categories := make([]string, 0, len(validCategories))
	for category := range validCategories {
		categories = append(categories, category)
	}
	return strings.Join(categories, ", ")
}

// EmployeeEfileDetails handles POST requests to fetch employee E-File details
// This is the main handler function that processes the API request
func EmployeeEfileDetails(w http.ResponseWriter, r *http.Request) {
	// Validate HTTP method - only POST requests are allowed
	if r.Method != http.MethodPost {
		encryptedResponse(w, http.StatusMethodNotAllowed, "Method not allowed, use POST")
		return
	}

	// Read and parse the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		encryptedResponse(w, http.StatusBadRequest, "Unable to read body")
		return
	}

	// Unmarshal JSON request body into struct
	var req EmployeeEfileRequest
	if err := json.Unmarshal(body, &req); err != nil {
		encryptedResponse(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	// Validate required fields
	if req.SessionID == "" {
		encryptedResponse(w, http.StatusBadRequest, "session_id is required")
		return
	}

	if req.Category == "" {
		encryptedResponse(w, http.StatusBadRequest, "category is required")
		return
	}

	// Validate category parameter
	if !validateCategory(req.Category) {
		errorMsg := "Invalid category. Valid categories are: " + getValidCategoriesString()
		encryptedResponse(w, http.StatusBadRequest, errorMsg)
		return
	}

	// Load encryption key from database using session ID
	if err := utils.LoadKeyFromDB(req.SessionID); err != nil {
		encryptedResponse(w, http.StatusInternalServerError, "Failed to load encryption key: "+err.Error())
		return
	}

	// Set token in header for authentication
	if req.Token != "" {
		r.Header.Set("token", req.Token)
	}

	// Authenticate the request
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	// Create logged handler for request logging and processing
	loggedHandler := auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Fetch employee data based on employee ID and category
		data, err := databasesad.GetEmployeeEFile(req.EmployeeID, req.Category)
		if err != nil {
			encryptedResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Prepare success response
		response := APIResponseEmployeeEfile{
			Status:  200,
			Message: "Success",
			Data:    data,
		}

		// Marshal and encrypt the response
		jsonResp, _ := json.Marshal(response)
		encrypted, err := utils.Encryptnew(jsonResp)
		if err != nil {
			encryptedResponse(w, http.StatusInternalServerError, "Encryption failed")
			return
		}

		// Log the response
		auth.SaveResponseLog(r, map[string]string{"Data": encrypted}, 200, "application/json", len(encrypted), string(body))

		// Send encrypted response
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"Data": encrypted})
	}))

	// Execute the handler
	loggedHandler.ServeHTTP(w, r)
}
