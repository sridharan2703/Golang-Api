// Package controllersofficeorder handles HTTP APIs for officeorder remarks.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 15-09-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 27-10-2025
package controllersofficeorder

import (
	"Hrmodule/auth"
	databaseofficeorder "Hrmodule/database/officeorder"
	"Hrmodule/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// --- Response Structure ---
type APIResponseOfficeRole struct {
	Status  int         `json:"Status"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

// --- Request Structure ---
type OfficeCommentsRequest struct {
	Data string `json:"Data"`
}

// --- Handler Function ---
func OfficeComments(w http.ResponseWriter, r *http.Request) {
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

	// Step 2: Parse incoming JSON
	var req OfficeCommentsRequest
	if err := json.Unmarshal(body, &req); err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Step 3: Decrypt incoming data
	parts := strings.Split(req.Data, "||")
	if len(parts) != 2 {
		log.Printf("Invalid data format")
		http.Error(w, "Invalid Data format", http.StatusBadRequest)
		return
	}

	pid := parts[0]
	encryptedPart := parts[1]

	// Step 4: Fetch AES key from DB
	key, err := utils.GetDecryptKey(pid)
	if err != nil {
		log.Printf("Key fetch failed: %v", err)
		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}

	// Step 5: Decrypt request payload
	decryptedJSON, err := utils.DecryptAES(encryptedPart, key)
	if err != nil {
		log.Printf("Decryption error: %v", err)
		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}

	// Step 6: Parse decrypted JSON to map
	var decryptedData map[string]interface{}
	if err := json.Unmarshal([]byte(decryptedJSON), &decryptedData); err != nil {
		log.Printf("Invalid decrypted JSON: %v", err)
		http.Error(w, "Invalid decrypted data", http.StatusBadRequest)
		return
	}

	// Step 7: Extract token and authenticate
	token, ok := decryptedData["token"].(string)
	if !ok || token == "" {
		log.Printf("Token not found in decrypted data")
		http.Error(w, "Token not found", http.StatusBadRequest)
		return
	}
	r.Header.Set("token", token)

	// Step 8: Validate token
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	// Step 9: Wrap inside logging middleware
	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate again before DB access
		if err := auth.IsValidIDFromRequest(r); err != nil {
			http.Error(w, "Invalid TOKEN provided", http.StatusBadRequest)
			return
		}

		// Step 10: Call DB layer (now using decryptedData)
		data, totalCount, err := databaseofficeorder.GetOfficeCommentsFromDB(w, r, decryptedData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Step 11: Prepare API response
		response := APIResponseOfficeRole{
			Status:  200,
			Message: "Success",
			Data: map[string]interface{}{
				"No Of Records": totalCount,
				"Records":       data,
			},
		}

		// Step 12: Marshal + encrypt response
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

		// Step 13: Save log and send to client
		auth.SaveResponseLog(
			r,
			finalResp,          // final response
			http.StatusOK,      // status code
			"application/json", // content type
			len(responseJSON),  // size
			string(body),       // original request
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(finalResp)
	})).ServeHTTP(w, r)
}
