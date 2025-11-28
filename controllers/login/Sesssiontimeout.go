// Package common contains APIs that are commonly used across the application and are grouped together for reusability.
//
// This API marks a user session as inactive upon logout or timeout.
//
// Path: Login Page
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 09-07-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 09-07-2025
package controllerslogin

import (
	"Hrmodule/auth"
	credentials "Hrmodule/dbconfig"
	"Hrmodule/utils"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

// Incoming encrypted request
type SessionTimeoutEncryptedRequest struct {
	Data string `json:"Data"` // pid||encryptedPayload
}

// API response
type APIResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Update session logout in DB
func UpdateSessionLogout(sessionId string, idleTimeout int) error {

	db, err := sql.Open("postgres", credentials.Getdatabasemeivan())
	if err != nil {
		return fmt.Errorf("DB open error: %v", err)
	}
	defer db.Close()

	// Correct query
	query := `
		UPDATE session_data 
		SET Is_Active = 0, idletimeout = $2, Logout_Date = NOW() 
		WHERE Session_Id = $1`

	_, err = db.Exec(query, sessionId, idleTimeout)
	if err != nil {
		return fmt.Errorf("update error: %v", err)
	}

	return nil
}

func SessionTimeoutHandler(w http.ResponseWriter, r *http.Request) {

	// Only POST allowed
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read full body
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(rawBody))

	// Parse JSON for Data
	var encReq SessionTimeoutEncryptedRequest
	if err := json.Unmarshal(rawBody, &encReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Split: pid || encryptedPayload
	parts := strings.Split(encReq.Data, "||")
	if len(parts) != 2 {
		http.Error(w, "Invalid Data format", http.StatusBadRequest)
		return
	}

	pid := parts[0]
	encryptedInput := parts[1]

	// Fetch AES Key from DB using P_ID
	key, err := utils.GetDecryptKey(pid)
	if err != nil {
		http.Error(w, "Failed to fetch decryption key", http.StatusUnauthorized)
		return
	}

	// Decrypt AES → JSON
	decryptedJSON, err := utils.DecryptAES(encryptedInput, key)
	if err != nil {
		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}

	// Parse decrypted JSON
	var decryptedData map[string]interface{}
	if err := json.Unmarshal([]byte(decryptedJSON), &decryptedData); err != nil {
		http.Error(w, "Invalid decrypted JSON", http.StatusBadRequest)
		return
	}

	// Extract token & session_id
	token, _ := decryptedData["token"].(string)
	//sessionID, _ := decryptedData["session_id"].(string)

	// Extract idletimeout (float64 → int)
	idletimeout := 0
	if v, ok := decryptedData["idletimeout"].(float64); ok {
		idletimeout = int(v)
	}

	// Validate required fields
	if token == "" {
		http.Error(w, "Missing required fields in decrypted payload", http.StatusBadRequest)
		return
	}

	// Inject token into header
	r.Header.Set("token", token)

	// Step: validate token + IP
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	// Wrap next stage
	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Validate token
		if err := auth.IsValidIDFromRequest(r); err != nil {
			http.Error(w, "Invalid TOKEN provided", http.StatusBadRequest)
			return
		}

		// Update the session
		err := UpdateSessionLogout(pid, idletimeout)

		var apiResp APIResponse
		if err != nil {
			apiResp = APIResponse{
				Status:  500,
				Message: "Failed to update session: " + err.Error(),
			}
		} else {
			apiResp = APIResponse{
				Status:  200,
				Message: fmt.Sprintf("Session updated successfully (idletimeout=%d)", idletimeout),
			}
		}

		// Encrypt response JSON
		jsonResp, _ := json.Marshal(apiResp)
		encryptedResp, _ := utils.EncryptAES(string(jsonResp), key)

		// Final output: pid||encryptedResponse
		finalResponse := map[string]string{
			"Data": fmt.Sprintf("%s||%s", pid, encryptedResp),
		}

		// Write JSON output
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(finalResponse)

	})).ServeHTTP(w, r)
}
