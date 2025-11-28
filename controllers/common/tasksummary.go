// Package commoncontrollers exposes API for Tasksummary.
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On:21-11-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 21-11-2025
package controllerscommon

import (
	"Hrmodule/auth"
	databasecommon "Hrmodule/database/common"
	"Hrmodule/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type APITaskSummary struct {
	Status  int         `json:"Status"`
	Message string      `json:"message"`
	Data    interface{} `json:"Data"`
}

type TaskSummaryRequest struct {
	Data string `json:"Data"`
}

func TaskSummary(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req TaskSummaryRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
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
		http.Error(w, "Key fetch failed", http.StatusUnauthorized)
		return
	}

	decryptedJSON, err := utils.DecryptAES(encryptedPart, key)
	if err != nil {
		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}

	var decryptedData map[string]interface{}
	json.Unmarshal([]byte(decryptedJSON), &decryptedData)

	// // REQUIRED FIELDS
	// if decryptedData["employee_id"] == nil ||
	// 	decryptedData["task_status_id"] == nil ||
	// 	decryptedData["priority"] == nil {
	// 	http.Error(w, "employee_id, task_status_id, priority required", http.StatusBadRequest)
	// 	return
	// }

	token, _ := decryptedData["token"].(string)
	if token == "" {
		http.Error(w, "Token missing", http.StatusBadRequest)
		return
	}

	r.Header.Set("token", token)

	// TOKEN + IP VALIDATION
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if err := auth.IsValidIDFromRequest(r); err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		result, count, err := databasecommon.GetTaskSummaryFromDB(decryptedData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := APITaskSummary{
			Status:  200,
			Message: "Success",
			Data: map[string]interface{}{
				"No Of Records": count,
				"Records":       result,
			},
		}

		respJSON, _ := json.Marshal(response)

		encryptedResp, err := utils.EncryptAES(string(respJSON), key)
		if err != nil {
			http.Error(w, "Encryption failed", http.StatusInternalServerError)
			return
		}

		finalResp := map[string]string{
			"Data": fmt.Sprintf("%s||%s", pid, encryptedResp),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(finalResp)

	})).ServeHTTP(w, r)
}
