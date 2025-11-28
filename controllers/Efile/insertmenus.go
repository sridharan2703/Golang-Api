// Package controllersefile contains structs and queries for InsertModules.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 24-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to insert data into category_role_map.

package controllersefile

import (
	"Hrmodule/auth"
	databaseefile "Hrmodule/database/Efile"
	"Hrmodule/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// APIResponseInsertModulesMaster defines the structure of the JSON response
type APIResponseInsertModulesMaster struct {
	Status       int    `json:"Status"`
	Message      string `json:"message"`
	RowsAffected int64  `json:"rows_affected,omitempty"`
	PId          string `json:"P_id,omitempty"`
}

// InsertModulesMasterRequest defines the expected request structure
type InsertModulesMasterRequest struct {
	Data string `json:"Data"`
}

// InsertModulesDecryptedRequestData defines the structure of decrypted data
type InsertModulesDecryptedRequestData struct {
	Token    string `json:"token"`
	PId      string `json:"P_id"`
	ModuleID string `json:"module_id"` // Can be comma-separated values
	RoleName string `json:"role_name"`
	Status   string `json:"status"`
}

func InsertModules(w http.ResponseWriter, r *http.Request) {
	// 1️⃣ Validate Request Method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	// 2️⃣ Read and Parse Request Body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req InsertModulesMasterRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// 3️⃣ Split and decrypt
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

	var decryptedData InsertModulesDecryptedRequestData
	if err := json.Unmarshal([]byte(decryptedJSON), &decryptedData); err != nil {
		http.Error(w, "Invalid decrypted data", http.StatusBadRequest)
		return
	}

	// Debug: Print decrypted data
	fmt.Printf("Decrypted data: %+v\n", decryptedData)

	// Validate required fields
	if decryptedData.Token == "" {
		http.Error(w, "Token not found in decrypted data", http.StatusBadRequest)
		return
	}

	if decryptedData.ModuleID == "" {
		http.Error(w, "module_id is required in decrypted data", http.StatusBadRequest)
		return
	}

	if decryptedData.RoleName == "" {
		http.Error(w, "role_name is required in decrypted data", http.StatusBadRequest)
		return
	}

	if decryptedData.Status == "" {
		http.Error(w, "status is required in decrypted data", http.StatusBadRequest)
		return
	}

	// Check for P_id consistency
	if decryptedData.PId != "" && decryptedData.PId != pid {
		fmt.Printf("P_id mismatch: received %s, expected %s\n", decryptedData.PId, pid)
		http.Error(w, "P_id mismatch", http.StatusBadRequest)
		return
	}

	r.Header.Set("token", decryptedData.Token)

	// 4️⃣ Authentication check
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	if err := auth.IsValidIDFromRequest(r); err != nil {
		http.Error(w, "Invalid TOKEN", http.StatusBadRequest)
		return
	}

	// 5️⃣ Parse module IDs from comma-separated string
	moduleIDs := strings.Split(decryptedData.ModuleID, ",")
	// Trim spaces from each module ID
	for i, moduleID := range moduleIDs {
		moduleIDs[i] = strings.TrimSpace(moduleID)
	}

	// Validate that we have at least one valid module ID
	validModuleIDs := []string{}
	for _, moduleID := range moduleIDs {
		if moduleID != "" {
			validModuleIDs = append(validModuleIDs, moduleID)
		}
	}

	if len(validModuleIDs) == 0 {
		http.Error(w, "No valid module IDs provided", http.StatusBadRequest)
		return
	}

	// 6️⃣ Business logic - insert multiple records into category_role_map
	rowsAffected, err := databaseefile.InsertMultipleCategoryRoleMap(
		validModuleIDs,
		decryptedData.RoleName,
		decryptedData.Status,
	)
	if err != nil {
		// Check if it's a duplicate error
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "already exists") {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create response
	response := APIResponseInsertModulesMaster{
		Status:       200,
		Message:      fmt.Sprintf("Successfully inserted %d records", rowsAffected),
		RowsAffected: rowsAffected,
		PId:          pid,
	}

	// 7️⃣ Marshal & encrypt before sending
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
		finalResp,
		http.StatusOK,
		"application/json",
		len(responseJSON),
		string(body),
	)

	// ✅ Send to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(finalResp)
}