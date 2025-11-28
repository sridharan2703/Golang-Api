// Package controllerscommon contains structs and queries for EmployeePresentScaleMaster.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 17-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all EmployeePresentScaleMaster.

package controllerscommon

import (
	"Hrmodule/auth"
	databasesad "Hrmodule/database/common"
	"Hrmodule/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// APIResponseEmployeePresentScaleMaster defines the structure of the JSON response
type APIResponseEmployeePresentScaleMaster struct {
	Status  int         `json:"Status"`
	Message string      `json:"message"`
	Data    interface{} `json:"Data"`
	PId     string      `json:"P_id,omitempty"`
}

// EmployeePresentScaleMasterRequest defines the expected request structure
type EmployeePresentScaleMasterRequest struct {
	Data string `json:"Data"`
}

// EmployeePresentScaleMasterDecryptedRequestData defines the structure of decrypted data
type EmployeePresentScaleMasterDecryptedRequestData struct {
	Token      string `json:"token"`
	PId        string `json:"P_id"`
	GradeGroup string `json:"gradegroup"`
}

func EmployeePresentScaleMaster(w http.ResponseWriter, r *http.Request) {
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

	var req EmployeePresentScaleMasterRequest
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

	var decryptedData EmployeePresentScaleMasterDecryptedRequestData
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

	// 5️⃣ Business logic - fetch EmployeePresentScaleMaster data with gradegroup filter
	data, err := databasesad.GetEmployeePresentScaleMaster(decryptedData.GradeGroup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create response with all fields
	response := APIResponseEmployeePresentScaleMaster{
		Status:  200,
		Message: "Success",
		Data:    data,
		PId:     pid, // Include P_id in the response
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