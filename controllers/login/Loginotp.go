// Package controllerslogin provides handlers for managing OTP-based
// authentication, including secure insertion of OTP records into the database,
// session tracking, request validation, and encrypted API responses.
//
// It ensures:
//   - Secure request validation using token-based authentication
//   - Insertion of OTP details (username, mobile number, session ID, etc.)
//   - Automatic expiry of OTPs after 45 seconds
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

type OTPDetails struct {
	Data string `json:"Data"`
}

// InsertOTPHandler inserts a new OTPDetails row
func InsertOTPHandler(w http.ResponseWriter, r *http.Request) {

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

	var req OTPDetails
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

	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// ----------------------
		// Extract fields safely
		// ----------------------

		// username
		Username, ok := decryptedData["username"].(string)
		if !ok || Username == "" {
			http.Error(w, "missing 'username' in request data", http.StatusBadRequest)
			return
		}

		// session_id

		// mobileno INT64
		var MobileNo int64
		switch v := decryptedData["mobileno"].(type) {
		case float64:
			MobileNo = int64(v)
		case string:
			MobileNo, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				http.Error(w, "invalid 'mobileno'", http.StatusBadRequest)
				return
			}
		default:
			http.Error(w, "missing or invalid 'mobileno'", http.StatusBadRequest)
			return
		}

		// otp INT
		var OTP int
		switch v := decryptedData["otp"].(type) {
		case float64:
			OTP = int(v)
		case string:
			OTP, err = strconv.Atoi(v)
			if err != nil {
				http.Error(w, "invalid 'otp'", http.StatusBadRequest)
				return
			}
		default:
			http.Error(w, "missing or invalid 'otp'", http.StatusBadRequest)
			return
		}

		// ----------------------
		// Insert into DB
		// ----------------------
		connectionString := credentials.Getdatabasemeivan()
		db, err := sql.Open("postgres", connectionString)
		if err != nil {
			http.Error(w, fmt.Sprintf("DB open error: %v", err), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		query := `
			INSERT INTO otp_details 
			(username, mobileno, otp, otpsendon, status, otpvalidtill, session_id, resend)
			VALUES ($1, $2, $3, NOW(), 0, NOW() + interval '45 seconds', $4, 0)
			RETURNING id;
		`

		var id int
		err = db.QueryRow(query,
			Username,
			MobileNo,
			OTP,
			pid,
		).Scan(&id)

		if err != nil {
			http.Error(w, "Error inserting: "+err.Error(), http.StatusInternalServerError)
			return
		}

		resp := map[string]interface{}{
			"message":    "OTP record inserted successfully",
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

	})).ServeHTTP(w, r)
}
