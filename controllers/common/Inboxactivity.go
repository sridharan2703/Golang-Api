// Package controllerscommon contains APIs for updating workflow master records.
//
// This API updates master records with badge, priority, and starred status based on task ID.
//
// Path: Workflow Management
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 26-08-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 27-10-2025
package controllerscommon

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

type TaskUpdateRequest struct {
	Data string `json:"Data"` // pid||encryptedPayload
}

type APIResponse struct {
	Status       int    `json:"status"`
	Message      string `json:"message"`
	RowsAffected int64  `json:"rows_affected"`
}

// ------------------------------------------------------------
// DB CALL
// ------------------------------------------------------------
func UpdateTaskAttributes(taskid string, badge, priority, starred *int) (int64, error) {

	db, err := sql.Open("postgres", credentials.Getdatabasemeivan())
	if err != nil {
		return 0, fmt.Errorf("DB open error: %v", err)
	}
	defer db.Close()

	if badge == nil && priority == nil && starred == nil {
		return 0, fmt.Errorf("at least one field must be provided for update")
	}

	query := `SELECT meivan.update_task_attributes($1, $2, $3, $4)`
	var rowsAffected int64

	err = db.QueryRow(query, taskid, badge, priority, starred).Scan(&rowsAffected)
	if err != nil {
		return 0, fmt.Errorf("stored procedure execution error: %v", err)
	}

	return rowsAffected, nil
}

// ------------------------------------------------------------
// HANDLER
// ------------------------------------------------------------
func TaskUpdateHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request", http.StatusBadRequest)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	var req TaskUpdateRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// ------------------------------------------------------------
	// 1. Split PID || encryptedPayload
	// ------------------------------------------------------------
	parts := strings.Split(req.Data, "||")
	if len(parts) != 2 {
		http.Error(w, "Invalid Data format", http.StatusBadRequest)
		return
	}

	pid := parts[0]
	encryptedPart := parts[1]

	// ------------------------------------------------------------
	// 2. Load AES Key for PID
	// ------------------------------------------------------------
	key, err := utils.GetDecryptKey(pid)
	if err != nil {
		http.Error(w, "Failed to get Decryption Key", http.StatusUnauthorized)
		return
	}

	// ------------------------------------------------------------
	// 3. Decrypt AES Payload → JSON
	// ------------------------------------------------------------
	decryptedJSON, err := utils.DecryptAES(encryptedPart, key)
	if err != nil {
		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}

	var decrypted map[string]interface{}
	if err := json.Unmarshal([]byte(decryptedJSON), &decrypted); err != nil {
		http.Error(w, "Invalid decrypted JSON", http.StatusBadRequest)
		return
	}

	// ------------------------------------------------------------
	// 4. Extract Fields (taskid, badge, priority, starred, token)
	// ------------------------------------------------------------
	token, _ := decrypted["token"].(string)
	if token == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}
	r.Header.Set("token", token)

	taskID, _ := decrypted["taskid"].(string)
	if taskID == "" {
		http.Error(w, "Missing taskid", http.StatusBadRequest)
		return
	}

	// badge (optional)
	var badge *int
	if v, ok := decrypted["badge"].(float64); ok {
		tmp := int(v)
		badge = &tmp
	}

	// priority (optional)
	var priority *int
	if v, ok := decrypted["priority"].(float64); ok {
		tmp := int(v)
		priority = &tmp
	}

	// starred (optional)
	var starred *int
	if v, ok := decrypted["starred"].(float64); ok {
		tmp := int(v)
		starred = &tmp
	}

	// ------------------------------------------------------------
	// 5. Validate Token + IP + API name
	// ------------------------------------------------------------
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if err := auth.IsValidIDFromRequest(r); err != nil {
			http.Error(w, "Invalid TOKEN", http.StatusBadRequest)
			return
		}

		// ------------------------------------------------------------
		// 6. UPDATE DB
		// ------------------------------------------------------------
		rowsAffected, err := UpdateTaskAttributes(taskID, badge, priority, starred)

		var resp APIResponse
		if err != nil {
			resp = APIResponse{Status: 500, Message: err.Error()}
		} else if rowsAffected == 0 {
			resp = APIResponse{Status: 404, Message: "Task not found", RowsAffected: 0}
		} else {
			resp = APIResponse{
				Status:       200,
				Message:      "Record updated",
				RowsAffected: rowsAffected,
			}
		}

		// ------------------------------------------------------------
		// 7. Encrypt Response
		// ------------------------------------------------------------
		respJSON, _ := json.Marshal(resp)
		encryptedResp, err := utils.EncryptAES(string(respJSON), key)
		if err != nil {
			http.Error(w, "Encryption failed", http.StatusInternalServerError)
			return
		}

		final := map[string]string{
			"Data": fmt.Sprintf("%s||%s", pid, encryptedResp),
		}

		// ------------------------------------------------------------
		// 8. SEND RESPONSE
		// ------------------------------------------------------------
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(final)

	})).ServeHTTP(w, r)
}
