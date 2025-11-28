// Package controllersofficeorder handles HTTP APIs for officeorder status update.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 15-09-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 15-09-2025
package controllersofficeorder

import (
	"Hrmodule/auth"
	credentials "Hrmodule/dbconfig"
	"Hrmodule/utils"
	"log"
	"strings"

	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

// Request structure for updating task status
type OfficeOrderUpdateTaskStatusRequest struct {
	Data string `json:"Data"`
}

// Standard response structure
type APIResponseforstatusupdate struct {
	Status  int         `json:"Status"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

// Helper to convert string to int with default fallback
func parseIntOrZeroforstatusupdate(s string) int {
	if s == "" {
		return 0
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

// Controller to update task status
func OfficeOrderUpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	// 1. Validate Request Method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	// 2. Read and Parse Request Body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req OfficeOrderUpdateTaskStatusRequest

	if err := json.Unmarshal(body, &req); err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// 3. Decrypt Data and Extract Token
	parts := strings.Split(req.Data, "||")
	if len(parts) != 2 {
		log.Printf("Invalid data format")
		http.Error(w, "Invalid Data format", http.StatusBadRequest)
		return
	}

	pid := parts[0]
	encryptedPart := parts[1]

	// Get decryption key from database
	key, err := utils.GetDecryptKey(pid)
	if err != nil {
		log.Printf("Key fetch failed: %v", err)
		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}

	// Decrypt the payload
	decryptedJSON, err := utils.DecryptAES(encryptedPart, key)
	if err != nil {
		log.Printf("Decryption error: %v", err)
		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}

	// Parse decrypted JSON to get token
	var decryptedData map[string]interface{}
	if err := json.Unmarshal([]byte(decryptedJSON), &decryptedData); err != nil {
		log.Printf("Invalid decrypted JSON: %v", err)
		http.Error(w, "Invalid decrypted data", http.StatusBadRequest)
		return
	}

	// 4️⃣ Extract fields
	token, _ := decryptedData["token"].(string)
	employeeID, _ := decryptedData["p_employeeid"].(string)
	coverPageNo, _ := decryptedData["p_coverpageno"].(string)
	taskstatusid, _ := decryptedData["p_taskstatusid"].(string)
	updatedby, _ := decryptedData["p_updatedby"].(string)

	if token == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}
	if employeeID == "" {
		http.Error(w, "Missing employeeid", http.StatusBadRequest)
		return
	}
	if coverPageNo == "" {
		http.Error(w, "Missing coverpageno", http.StatusBadRequest)
		return
	}
	token, ok := decryptedData["token"].(string)
	if !ok || token == "" {
		log.Printf("Token not found in decrypted data")
		http.Error(w, "Token not found", http.StatusBadRequest)
		return
	}

	// 4. Set Token Header for Authentication
	r.Header.Set("token", token)

	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := auth.IsValidIDFromRequest(r); err != nil {
			http.Error(w, "Invalid TOKEN", http.StatusBadRequest)
			return
		}

		connectionString := credentials.Getdatabasemeivan()
		db, err := sql.Open("postgres", connectionString)
		if err != nil {
			http.Error(w, fmt.Sprintf("DB connection error: %v", err), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		// Convert TaskStatusID to int
		taskStatusID := parseIntOrZeroforstatusupdate(taskstatusid)

		// Call the stored procedure
		_, err = db.Exec(`CALL updatetaskstatus($1, $2, $3, $4)`,
			coverPageNo, employeeID, taskStatusID, updatedby,
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("DB update error: %v", err), http.StatusInternalServerError)
			return
		}

		// Success response
		response := APIResponseforstatusupdate{
			Status:  200,
			Message: "Task status updated successfully",
			Data:    nil,
		}

		// 6️⃣ Marshal & encrypt before sending
		responseJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Response marshal failed", http.StatusInternalServerError)
			return
		}

		encryptedResponse, err := utils.EncryptAES(string(responseJSON), key)
		if err != nil {
			http.Error(w, "Response encryption failed", http.StatusInternalServerError)
			return
		}

		finalResp := map[string]string{
			"Data": fmt.Sprintf("%s||%s", pid, encryptedResponse),
		}

		// ✅ Save exactly what is sent to client
		auth.SaveResponseLog(
			r,
			finalResp,          // only final response
			http.StatusOK,      // status code
			"application/json", // content type
			len(responseJSON),  // size
			string(body),       // original request
		)

		// ✅ Send to client
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(finalResp)
	})).ServeHTTP(w, r)
}
