// Package controllersofficeorder handles HTTP APIs for officeorder count.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 15-09-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 06-11-2025
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

// APIResponseNeedGenerate defines the structure of the JSON response
// returned by the NeedGenerateHandler API.
type APIResponseNeedGenerate struct {
	Status  int         `json:"Status"`
	Message string      `json:"message"`
	Data    interface{} `json:"Data"`
}

// NeedGenerateRequest defines the expected JSON request body structure
// This now expects encrypted Data format: "P_id||encrypted_payload"
type NeedGenerateRequest struct {
	Data string `json:"Data"`
}

// NeedGenerateHandler handles POST requests to check the count
//
// It validates the request method, decrypts the incoming data to extract the token,
// verifies session and token, and calls the database layer to get counts of orders
// that need generation, ongoing, saved and held, and completed.
func NeedGenerateHandler(w http.ResponseWriter, r *http.Request) {
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

	var req NeedGenerateRequest
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

	// 6. Log Request Info and Process Business Logic
	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get data from database
		data, err := databaseofficeorder.GetCombinedNeedGenerate()
		if err != nil {
			log.Printf("Database error in GetCombinedNeedGenerate: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Prepare the final response structure
		response := APIResponseNeedGenerate{
			Status:  200,
			Message: "Success",
			Data: map[string]interface{}{
				"need_to_generate": data.Postgres.NeedToGenerate,
				"ongoing":          data.Postgres.Ongoing,
				"saveandhold":      data.Postgres.SaveAndHold,
				"complete":         data.Postgres.Complete,
			},
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
