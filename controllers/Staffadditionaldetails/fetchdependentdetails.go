// Package controllerssad contains structs and queries for Staff Additional details API.

// --- Creator's Info ---
// Creator: Vaishnavi
// Created On: 04-11-2025
// Last Modified By:  Rovita
// Last Modified Date: 12-1-2025
// This api is to feth the all the dependent details
package controllerssad

import (
	"Hrmodule/auth"
	databasesad "Hrmodule/database/Staffadditionaldetails"
	"Hrmodule/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type EmployeeDependentDetailsRequest struct {
	Data string `json:"Data"`
}

type DependentAPIResponse struct {
	Status  int         `json:"Status"`
	Message string      `json:"message"`
	Data    interface{} `json:"Data"`
}

func EmployeeDependentDetailsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("=== EmployeeDependentDetailsHandler Started ===")

	// 1️⃣ Validate Request Method
	if r.Method != http.MethodPost {
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	// 2️⃣ Read and Parse Request Body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Raw request body: %s", string(body))

	var req EmployeeDependentDetailsRequest
	if err := json.Unmarshal(body, &req); err != nil {
		log.Printf("JSON unmarshal error: %v", err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	log.Printf("Request Data received: %s", req.Data)

	// 3️⃣ Split and decrypt
	parts := strings.Split(req.Data, "||")
	if len(parts) != 2 {
		log.Printf("Invalid Data format, parts: %v", parts)
		http.Error(w, "Invalid Data format", http.StatusBadRequest)
		return
	}
	pid := parts[0]
	encryptedPart := parts[1]

	log.Printf("PID: %s", pid)

	key, err := utils.GetDecryptKey(pid)
	if err != nil {
		log.Printf("Decryption key fetch failed: %v", err)
		http.Error(w, "Decryption key fetch failed", http.StatusUnauthorized)
		return
	}

	decryptedJSON, err := utils.DecryptAES(encryptedPart, key)
	if err != nil {
		log.Printf("Decryption failed: %v", err)
		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}

	log.Printf("Decrypted JSON: %s", decryptedJSON)

	var decryptedData map[string]interface{}
	if err := json.Unmarshal([]byte(decryptedJSON), &decryptedData); err != nil {
		log.Printf("Invalid decrypted data: %v", err)
		http.Error(w, "Invalid decrypted data", http.StatusBadRequest)
		return
	}

	token, ok := decryptedData["token"].(string)
	if !ok || token == "" {
		log.Printf("Token not found in decrypted data: %+v", decryptedData)
		http.Error(w, "Token not found", http.StatusBadRequest)
		return
	}
	r.Header.Set("token", token)

	// 4️⃣ Authentication check
	log.Println("Performing authentication check...")
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		log.Println("Authentication failed")
		return
	}

	// Validate token from request
	if err := auth.IsValidIDFromRequest(r); err != nil {
		log.Printf("Token validation failed: %v", err)
		http.Error(w, "Invalid TOKEN", http.StatusBadRequest)
		return
	}

	// 5️⃣ Extract employee ID and business logic
	employeeID, ok := decryptedData["employeeid"].(string)
	if !ok || employeeID == "" {
		log.Printf("Missing employeeid in decrypted data: %+v", decryptedData)
		http.Error(w, "Missing 'employeeid' in request data", http.StatusBadRequest)
		return
	}

	log.Printf("Fetching dependent details for employee: %s", employeeID)
	data, err := databasesad.FetchEmployeeDependentDetails(employeeID)
	if err != nil {
		log.Printf("Database fetch error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Database response: %+v", data)

	response := DependentAPIResponse{
		Status:  200,
		Message: "Success",
		Data:    data,
	}

	// 6️⃣ Marshal response to JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("Response marshal failed: %v", err)
		http.Error(w, "Response marshal failed", http.StatusInternalServerError)
		return
	}

	log.Printf("Response JSON: %s", string(responseJSON))

	// 7️⃣ Encrypt the response
	encryptedResponse, err := utils.EncryptAES(string(responseJSON), key)
	if err != nil {
		log.Printf("Response encryption failed: %v", err)
		http.Error(w, "Response encryption failed", http.StatusInternalServerError)
		return
	}

	finalResp := map[string]string{
		"Data": fmt.Sprintf("%s||%s", pid, encryptedResponse),
	}

	log.Printf("Final encrypted response prepared")

	// ✅ Save exactly what is sent to client
	auth.SaveResponseLog(
		r,
		finalResp,          // only final response
		http.StatusOK,      // status code
		"application/json", // content type
		len(responseJSON),  // size
		string(body),       // original request
	)

	// ✅ Send encrypted response to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(finalResp); err != nil {
		log.Printf("Error sending response: %v", err)
		return
	}

	log.Println("=== EmployeeDependentDetailsHandler Completed Successfully ===")
}
