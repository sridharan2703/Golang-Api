// Package auth provides authentication and authorization functionality,
// including client IP validation, token validation, and stored procedure integration.
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On:07-07-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 09-07-2025
//
//
//hi
package auth

import (
	credentials "Hrmodule/dbconfig"
	"Hrmodule/utils"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

// ResponseLogData represents the structure of data to be saved
type ResponseLogData struct {
	Timestamp      string              `json:"timestamp"`
	APIUrl         string              `json:"api_url"`
	Method         string              `json:"method"`
	ClientIP       string              `json:"client_ip"`
	Headers        map[string][]string `json:"headers"`
	QueryParams    map[string][]string `json:"query_params"`
	RequestBody    interface{}         `json:"request_body"`
	ResponseBody   interface{}         `json:"response_body"`
	StatusCode     int                 `json:"status_code"`
	ContentType    string              `json:"content_type"`
	ResponseSize   int                 `json:"response_size_bytes"`
	ProcessingTime string              `json:"processing_time"`
}

func IsValidIDFromRequest(r *http.Request) error {
	var token string

	// 1. Try to get token from the header
	token = r.Header.Get("token")

	// 2. Try to read from body only if it's POST and token is still empty
	if token == "" && r.Method == http.MethodPost {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return errors.New("unable to read request body")
		}
		// Restore the body so it can be read again later
		r.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

		var bodyData struct {
			Token string `json:"token"`
		}
		if err := json.Unmarshal(bodyBytes, &bodyData); err == nil {
			token = bodyData.Token
		}
	}

	// 3. Fallback to query string
	if token == "" {
		token = r.URL.Query().Get("token")
	}

	// 4. Token validation
	for _, char := range token {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return errors.New("invalid character in TOKEN")
		}
	}

	// Optional: match expected token
	// if token != "your_expected_token" {
	//     return errors.New("unauthorized")
	// }

	return nil
}

// LogRequestInfo logs the client's IP address and forwards the request to the given handler.
//
// Parameters:
//   - handler: The HTTP handler to wrap.
//
// Returns:
//   - A wrapped handler that logs the client IP before executing.
func LogRequestInfo(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		fmt.Printf("Client IP Address: %s\n", clientIP)
		handler(w, r)
	}
}

// ValidateAPI calls the stored procedure `API_Validation` to determine
// if the API access is valid, and logs the request and its result to a database.
//
// Parameters:
//   - APIName: The name of the API being accessed.
//   - clientIPAddress: The IP address of the requester.
//   - IDKey: The token or identifier used to validate the request.
//   - requestURL: The full URL of the incoming request.
//
// Returns:
//   - A boolean indicating if the request is valid.
//   - A status message returned from the stored procedure.
//   - An error if something goes wrong during validation or logging.
func ValidateAPI(APIName, clientIPAddress, IDKey, requestURL string) (bool, string, error) {

	// Get PostgreSQL connection string
	connectionString := credentials.Getdatabasemeivan()

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return false, "", fmt.Errorf("DB open error: %v", err)
	}
	defer db.Close()

	// Call stored function: api_validation_new
	var statusMessage string
	err = db.QueryRow(
		`SELECT api_hr.api_validation_new($1, $2, $3);`,
		APIName, clientIPAddress, IDKey,
	).Scan(&statusMessage)

	if err != nil {
		return false, "", fmt.Errorf("DB call failed: %v", err)
	}

	// Prepare status and error message
	status := ""
	errorMessage := ""

	if statusMessage == "Success" {
		status = "Success"
	} else {
		errorMessage = statusMessage
	}

	// Insert request log into client_request table
	_, err = db.Exec(`
		INSERT INTO api_hr.client_request 
			(ip_address, request_data, response_data, status, error, request_on, response_on, updated_on)
		VALUES 
			($1, $2, '', $3, $4, NOW(), NOW(), NOW())
	`,
		clientIPAddress, requestURL, status, errorMessage,
	)

	if err != nil {
		return false, "", fmt.Errorf("Insert log failed: %v", err)
	}

	// Return whether request is valid
	return statusMessage == "Success", statusMessage, nil
}

// Responseset represents the standard API error response format.
type Responseset struct {
	Status  int      `json:"Status"`  // HTTP-like status code
	Message string   `json:"Message"` // Message describing the outcome
	Data    []string `json:"Data"`    // Additional data (usually empty for errors)
}

// HandleRequestforapiname_ipaddress_token validates a request by extracting relevant metadata (API name, IP address, token),
// invoking the `ValidateAPI` function, and returning an appropriate response.
//
// Parameters:
//   - w: The HTTP response writer.
//   - r: The HTTP request containing validation metadata.
//
// Returns:
//   - True if the request is authorized and passes all validations.
//   - False otherwise, and writes an error response directly to the client.
// func HandleRequestfor_apiname_ipaddress_token(w http.ResponseWriter, r *http.Request) bool {
// 	u, err := url.Parse(r.URL.String())
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return false
// 	}

func HandleRequestfor_apiname_ipaddress_token(w http.ResponseWriter, r *http.Request) bool {
	// Extract the values from the request
	u, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return false
	}

	// Extract the APIName from the URL path
	pathParts := strings.Split(u.Path, "/")
	var APIName string
	if len(pathParts) > 1 {
		APIName = pathParts[1] // Assuming "/Facultydetails" is part of the path
	}

	// Extract the clientIPAddress from the request
	clientIPAddress := strings.Split(r.RemoteAddr, ":")[0]

	// Extract the token from the header, body or query string
	var IDKey string

	// First check for token in header using 'X-Validation-Token'
	IDKey = r.Header.Get("token")

	// If the token is not in the header, check the body if it's a POST request
	if IDKey == "" && r.Method == http.MethodPost {
		var bodyData struct {
			Token string `json:"token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&bodyData); err == nil {
			IDKey = bodyData.Token
		}
	}

	// If the token is still not found, check the query string
	if IDKey == "" {
		queryValues := r.URL.Query()
		if idValues, ok := queryValues["token"]; ok && len(idValues) > 0 {
			IDKey = idValues[0]
		}
	}

	// Get the entire request URL as a string
	requestURL := r.URL.String()

	// Validate the API using the token, client IP, and APIName
	isValid, statusMessage, err := ValidateAPI(APIName, clientIPAddress, IDKey, requestURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return false
	}

	// Handle various API status messages
	switch statusMessage {
	case "Invalid_Key":
		return respondWithError(w, 400, "Invalid_Key")
	case "Invalid_APIName":
		return respondWithError(w, 401, "Invalid_APIName")
	case "Invalid_IPAddress":
		return respondWithError(w, 402, "Invalid_IPAddress")
	case "Inactive_APIName":
		return respondWithError(w, 403, "Inactive_APIName")
	case "Inactive_Vendor":
		return respondWithError(w, 404, "Inactive_Vendor")
	case "Inactive_Ip_Address":
		return respondWithError(w, 405, "Inactive_Ip_Address")
	case "UnauthorizedUser":
		return respondWithError(w, 406, "UnauthorizedUser")
	case "Invalid_RollNo":
		return respondWithError(w, 407, "Invalid_RollNo")
	}

	// If validation fails, return a forbidden error
	if !isValid {
		return respondWithError(w, http.StatusForbidden, statusMessage)
	}

	return true
}

// respondWithError writes an encrypted JSON error response to the client.
//
// It builds a structured error object (`Responseset`), marshals it to JSON,
// encrypts the response using AES-GCM, and sends it as a JSON object with an "encrypted" key.
//
// Parameters:
//   - w: The HTTP response writer.
//   - statusCode: The HTTP status code to send.
//   - message: The error message.
//
// Returns:
//   - false (for convenience use in calling code).
func respondWithError(w http.ResponseWriter, statusCode int, message string) bool {
	response := Responseset{
		Status:  statusCode,
		Message: message,
		Data:    []string{},
	}

	responseJSON, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	// Encrypt the error response
	encrypted, err := utils.Encrypt(responseJSON)
	if err != nil {
		http.Error(w, "Encryption failed", http.StatusInternalServerError)
		return false
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"Data": encrypted,
	})

	return false
}

// ensureOriginHeader checks and adds "Origin" header if missing due to casing issues
func ensureOriginHeader(r *http.Request) map[string][]string {
	headers := make(map[string][]string)
	for k, v := range r.Header {
		headers[k] = v
	}

	// Ensure "Origin" is preserved if provided in lowercase (some clients do this)
	if _, ok := headers["Origin"]; !ok {
		if v := r.Header.Get("origin"); v != "" {
			headers["Origin"] = []string{v}
		}
	}

	return headers
}

// Save the respone log of user activity
func SaveResponseLog(
	r *http.Request,
	responseBody interface{},
	statusCode int,
	contentType string,
	responseSize int,
	requestBodyRaw string,
) {
	startTime := time.Now() // start timer

	var parsedRequestBody map[string]interface{}
	var pid string

	// Extract P_id from request body
	if requestBodyRaw != "" {

		// Try normal JSON parsing first
		if err := json.Unmarshal([]byte(requestBodyRaw), &parsedRequestBody); err == nil {

			// CASE 1: Direct JSON { "P_id": "xxxx" }
			if id, ok := parsedRequestBody["P_id"].(string); ok {
				pid = id
			}

			// CASE 2: Encrypted format: { "Data": "P_id||EncryptedString" }
			if pid == "" {
				if data, ok := parsedRequestBody["Data"].(string); ok {
					parts := strings.Split(data, "||")
					if len(parts) > 0 && parts[0] != "" {
						pid = parts[0] // first part is P_id
					}
				}
			}
		}
	}

	// If P_id missing → use default
	if pid == "" {
		pid = "unknown_pid"
	}

	// Timestamp format for folder + filename
	now := time.Now()
	dateFolder := now.Format("02_01_2006") // e.g. 29_11_2025
	timeStamp := now.Format("15_04_05")

	apiName := strings.TrimPrefix(r.URL.Path, "/")

	// Build directory: /var/log/Hrmodule/DATE/PID/
	baseDir := "/var/log/Hrmodule"

	// ✔ Ensure Hrmodule folder exists
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		fmt.Println("Error creating base Hrmodule directory:", err)
		return
	}

	pidFolder := filepath.Join(baseDir, dateFolder, pid)

	// ✔ Create date + PID folder
	if err := os.MkdirAll(pidFolder, 0755); err != nil {
		fmt.Println("Error creating log directory:", err)
		return
	}

	// File: APIName__Time.json
	fileName := fmt.Sprintf("%s__%s.json", apiName, timeStamp)
	fullFilePath := filepath.Join(pidFolder, fileName)

	headers := ensureOriginHeader(r)

	var finalRequestBody interface{}
	if parsedRequestBody != nil {
		finalRequestBody = parsedRequestBody
	} else {
		finalRequestBody = requestBodyRaw
	}

	logData := ResponseLogData{
		Timestamp:      now.Format("2006-01-02 15:04:05.000"),
		APIUrl:         r.URL.Path,
		Method:         r.Method,
		ClientIP:       r.RemoteAddr,
		Headers:        headers,
		QueryParams:    r.URL.Query(),
		RequestBody:    finalRequestBody,
		ResponseBody:   responseBody,
		StatusCode:     statusCode,
		ContentType:    contentType,
		ResponseSize:   responseSize,
		ProcessingTime: time.Since(startTime).String(),
	}

	logJSON, err := json.MarshalIndent(logData, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling log data:", err)
		return
	}

	if err := os.WriteFile(fullFilePath, logJSON, 0644); err != nil {
		fmt.Println("Error writing log file:", err)
		return
	}

	fmt.Println("Response log saved to", fullFilePath)
}
