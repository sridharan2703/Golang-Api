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

	_ "github.com/lib/pq"
)

// Datadecryptkey represents the expected JSON structure for a session timeout request.
type Datadecryptkey struct {
	SessionID  string `json:"session_id"` // SessionID is the identifier of the session to be updated.
	Token      string `json:"token"`      // Token can also come from request body
	DecryptKey string `json:"decryptkey"` // decryptkey is now a STRING
}

// APIResponsefordatadecrypt defines the JSON response structure.
type APIResponsefordatadecrypt struct {
	Status     int    `json:"status"`     // HTTP-like status code
	Message    string `json:"message"`    // Human-readable message
	SessionID  string `json:"session_id"` // Echo back updated session_id
	DecryptKey string `json:"decryptkey"` // Echo back updated decryptkey
}

// datadecrypt updates decryptkey in session_data table
func datadecrypt(sessionId string, decryptkey string) error {
	connectionString := credentials.Getdatabasemeivan()

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return fmt.Errorf("DB open error: %v", err)
	}
	defer db.Close()

	query := `
		UPDATE session_data
		SET decryptkey = $1
		WHERE Session_Id = $2
	`

	_, err = db.Exec(query, decryptkey, sessionId)
	if err != nil {
		return fmt.Errorf("update error: %v", err)
	}

	return nil
}

// DatadecryptHandler handles POST requests to the /datadecrypt endpoint.
func DatadecryptHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	var req Datadecryptkey
	_ = json.Unmarshal(body, &req)

	if req.Token != "" {
		r.Header.Set("token", req.Token)
	}

	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	loggedHandler := auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := auth.IsValidIDFromRequest(r); err != nil {
			http.Error(w, "Invalid TOKEN provided", http.StatusBadRequest)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		if req.SessionID == "" {
			http.Error(w, "Missing required field: session_id", http.StatusBadRequest)
			return
		}

		err := datadecrypt(req.SessionID, req.DecryptKey)

		var response APIResponsefordatadecrypt
		if err != nil {
			response = APIResponsefordatadecrypt{
				Status:     500,
				Message:    "Failed to update session: " + err.Error(),
				SessionID:  req.SessionID,
				DecryptKey: req.DecryptKey,
			}
		} else {
			response = APIResponsefordatadecrypt{
				Status:     200,
				Message:    "Session updated successfully",
				SessionID:  req.SessionID,
				DecryptKey: req.DecryptKey,
			}
		}

		responseBytes, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Failed to serialize JSON", http.StatusInternalServerError)
			return
		}

		encrypted, err := utils.Encrypt(responseBytes)
		if err != nil {
			http.Error(w, "Encryption failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"Data": encrypted,
		})
	}))

	loggedHandler.ServeHTTP(w, r)
}
