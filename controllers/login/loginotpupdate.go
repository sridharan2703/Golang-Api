// Package controllerslogin provides update the otp into database
//
// It ensures:
//   - Secure request validation using token-based authentication
//   - Validation of OTP details (username, mobile number, session ID, etc.)
//   - Automatic expiry check for OTPs using validity window
//   - Updates OTP status and verification timestamp on successful validation
//   - Encrypted response payloads for added security
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On: 26-08-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 26-08-2025
package controllerslogin

import (
	"Hrmodule/auth"
	credentials "Hrmodule/dbconfig"
	"Hrmodule/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

// ValidateOTPRequest represents the request body for OTP validation
type ValidateOTPRequest struct {
	Data string `json:"Data"`
}

func ValidateOTPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Step 1: Parse request body
	var req ValidateOTPRequest

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	parts := strings.Split(req.Data, "||")
	if len(parts) != 2 {
		http.Error(w, "Invalid Data format", http.StatusBadRequest)
		return
	}
	pid := parts[0]
	encryptedPart := parts[1]

	key, err := utils.GetDecryptKey(pid)
	if err != nil {
		http.Error(w, "Decryption key fetch failed", http.StatusUnauthorized)
		return
	}

	decryptedJSON, err := utils.DecryptAES(encryptedPart, key)
	if err != nil {
		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}

	var decryptedData map[string]interface{}
	if err := json.Unmarshal([]byte(decryptedJSON), &decryptedData); err != nil {
		http.Error(w, "Invalid decrypted data", http.StatusBadRequest)
		return
	}

	token, ok := decryptedData["token"].(string)
	if !ok || token == "" {
		http.Error(w, "Token not found", http.StatusBadRequest)
		return
	}
	r.Header.Set("token", token)

	// Step 3: Authenticate
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// username
		Username, ok := decryptedData["username"].(string)
		if !ok || strings.TrimSpace(Username) == "" {
			http.Error(w, "missing or invalid 'username'", http.StatusBadRequest)
			return
		}

		// mobileno (int64)
		var MobileNo int64
		switch v := decryptedData["mobileno"].(type) {
		case float64:
			MobileNo = int64(v)
		case string:
			var err error
			MobileNo, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				http.Error(w, "invalid 'mobileno'", http.StatusBadRequest)
				return
			}
		default:
			http.Error(w, "missing or invalid 'mobileno'", http.StatusBadRequest)
			return
		}

		// otp (int)
		var OTP int
		switch v := decryptedData["otp"].(type) {
		case float64:
			OTP = int(v)
		case string:
			var err error
			OTP64, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				http.Error(w, "invalid 'otp'", http.StatusBadRequest)
				return
			}
			OTP = int(OTP64)
		default:
			http.Error(w, "missing or invalid 'otp'", http.StatusBadRequest)
			return
		}

		// Step 7: DB Connection
		connectionString := credentials.Getdatabasemeivan()
		db, err := sql.Open("postgres", connectionString)
		if err != nil {
			http.Error(w, fmt.Sprintf("DB open error: %v", err), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		// Step 8: OTP Validation Query
		checkQuery := `
			SELECT 
				id,
				CASE WHEN (otpverifiedon IS NULL AND status = 0 AND otpvalidtill >= NOW()) 
					THEN '1' 
					ELSE '0' 
				END as validcheck
			FROM otp_details 
			WHERE username = $1 
			  AND mobileno = $2 
			  AND session_id = $3 
			  AND otp = $4
			  AND status = 0
			  AND otpverifiedon IS NULL 
			ORDER BY otpsendon DESC 
			LIMIT 1;
		`

		var id int
		var validCheck string

		err = db.QueryRow(checkQuery, Username, MobileNo, pid, OTP).Scan(&id, &validCheck)

		if err != nil {
			if err == sql.ErrNoRows {
				resp := map[string]interface{}{
					"success":    false,
					"message":    "Invalid OTP or OTP expired",
					"validcheck": "0",
				}

				// Encrypt response
				jsonResponse, _ := json.Marshal(resp)
				encryptedResponse, _ := utils.EncryptAES(string(jsonResponse), key)

				finalResp := map[string]string{
					"Data": fmt.Sprintf("%s||%s", pid, encryptedResponse),
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(finalResp)
				return
			}

			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// -----------------------------
		// OTP Valid → Update status
		// -----------------------------
		if validCheck == "1" {

			updateQuery := `
				UPDATE otp_details
				SET otpverifiedon = NOW(), status = 1
				WHERE id = $1
			`

			_, err = db.Exec(updateQuery, id)
			if err != nil {
				http.Error(w, "Error updating OTP verification: "+err.Error(), http.StatusInternalServerError)
				return
			}

			resp := map[string]interface{}{
				"success":    true,
				"message":    "OTP verified successfully",
				"id":         id,
				"session_id": pid,
			}

			jsonResponse, _ := json.Marshal(resp)
			encryptedResponse, _ := utils.EncryptAES(string(jsonResponse), key)

			finalResp := map[string]string{
				"Data": fmt.Sprintf("%s||%s", pid, encryptedResponse),
			}

			auth.SaveResponseLog(
				r,
				finalResp,
				http.StatusOK,
				"application/json",
				len(jsonResponse),
				string(body),
			)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(finalResp)
			return
		}

	})).ServeHTTP(w, r)
}
