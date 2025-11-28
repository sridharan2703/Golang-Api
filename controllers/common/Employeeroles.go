// Package commoncontrollers handles HTTP request routing, authentication,
// and API response formatting for the application.
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On:30-07-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 30-07-2025
package controllerscommon

import (
	"Hrmodule/auth"
	database "Hrmodule/database/common"
	"Hrmodule/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// APIResponseforDefaultRoleName defines the standard structure of the API response.
type APIResponseforDefaultRoleName struct {
	Status  int         `json:"Status"`
	Message string      `json:"message"`
	Data    interface{} `json:"Data"`
}

// Struct for request body (for token injection + flexibility for other fields later)
type DefaultRoleNameRequest struct {
	Data string `json:"Data"`
}

// DefaultRoleName handles the HTTP POST request to fetch DefaultRoleName data for Employees.
func DefaultRoleName(w http.ResponseWriter, r *http.Request) {
	// Allow only POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	// Step 1: Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req DefaultRoleNameRequest

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

	token, ok := decryptedData["token"].(string)
	if !ok || token == "" {
		log.Printf("Token not found in decrypted data")
		http.Error(w, "Token not found", http.StatusBadRequest)
		return
	}

	// 4. Set Token Header for Authentication
	r.Header.Set("token", token)

	// 5. Authentication Check
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	// Step 4: Logging middleware
	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Step 5: Validate token
		if err := auth.IsValidIDFromRequest(r); err != nil {
			http.Error(w, "Invalid TOKEN provided", http.StatusBadRequest)
			return
		}

		// Step 6: DB query
		DefaultRoleNameData, totalCount, err := database.DefaultRoleNamedatabase(decryptedData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Step 7: Response struct
		response := APIResponseforDefaultRoleName{
			Status:  200,
			Message: "Success",
			Data: map[string]interface{}{
				"No Of Records": totalCount,
				"Records":       DefaultRoleNameData,
			},
		}

		// Step 8: Marshal to JSON
		responseJSON, err := json.MarshalIndent(response, "", "  ")

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
