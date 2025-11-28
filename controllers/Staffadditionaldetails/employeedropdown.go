// Package modelssad contains structs and queries for Staff Additional details API.

// --- Creator's Info ---
// Creator: Vaishnavi
// Created On: 04-11-2025
// Last Modified By:  Rovita
// Last Modified Date: 12-1-2025
// This api is to feth the all active employees.
package controllerssad

import (
	"Hrmodule/auth"
	databasesad "Hrmodule/database/Staffadditionaldetails"
	"Hrmodule/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// APIResponseEmployeeDropdown defines the structure of the JSON response
type APIResponseEmployeeDropdown struct {
	Status  int         `json:"Status"`
	Message string      `json:"message"`
	Data    interface{} `json:"Data"`
}

// EmployeeDropdownRequest defines the expected request structure
type EmployeeDropdownRequest struct {
	Data string `json:"Data"`
}

func EmployeeDropdown(w http.ResponseWriter, r *http.Request) {
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

	var req EmployeeDropdownRequest
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

	var decryptedData map[string]interface{}
	if err := json.Unmarshal([]byte(decryptedJSON), &decryptedData); err != nil {
		http.Error(w, "Invalid decrypted data", http.StatusBadRequest)
		return
	}

	// ✅ Validate Type inside decrypted data (only once)
	typeVal, ok := decryptedData["type"].(string)
	if !ok || typeVal != "adminefile" {
		http.Error(w, "Invalid type", http.StatusBadRequest)
		return
	}

	// Debug: Print decrypted data
	fmt.Printf("Decrypted data: %+v\n", decryptedData)

	// Extract token
	token, ok := decryptedData["token"].(string)
	if !ok || token == "" {
		http.Error(w, "Token not found in decrypted data", http.StatusBadRequest)
		return
	}

	// Check for P_id with different possible field names (case-insensitive)
	var p_id string
	if val, ok := decryptedData["P_id"].(string); ok && val != "" {
		p_id = val
	} else if val, ok := decryptedData["p_id"].(string); ok && val != "" {
		p_id = val
	} else if val, ok := decryptedData["PId"].(string); ok && val != "" {
		p_id = val
	} else if val, ok := decryptedData["pid"].(string); ok && val != "" {
		p_id = val
	} else {
		p_id = pid
		fmt.Printf("Using pid from request as P_id: %s\n", p_id)
	}

	if p_id != "" && p_id != pid {
		fmt.Printf("P_id mismatch: received %s, expected %s\n", p_id, pid)
		http.Error(w, "P_id mismatch", http.StatusBadRequest)
		return
	}

	r.Header.Set("token", token)

	// 4️⃣ Authentication check
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	if err := auth.IsValidIDFromRequest(r); err != nil {
		http.Error(w, "Invalid TOKEN", http.StatusBadRequest)
		return
	}

	// 5️⃣ Business logic - fetch employee dropdown data
	data, err := databasesad.GetEmployeeDropdown()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := APIResponseEmployeeDropdown{
		Status:  200,
		Message: "Success",
		Data:    data,
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
