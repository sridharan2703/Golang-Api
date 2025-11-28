// // Package controllerssad handles API logic for Staff Additional Details.
// //
// // --- Creator's Info ---
// // Creator: Rovita
// // Created On: 11-11-2025
// // Last Modified By:
// // Last Modified Date:
// // Description: API to insert or update active employee personal basic details.
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

// APIResponseBasicDetails defines the structure of the JSON response
type APIResponseBasicDetails struct {
	Status  int         `json:"Status"`
	Message string      `json:"message"`
	Data    interface{} `json:"Data"`
}

// SadBasicRequest defines the expected request structure
type SadBasicRequest struct {
	Data string `json:"Data"`
}

// SadBasicDetails handles both insert and update operations using master_sad stored procedure
func SadBasicDetails(w http.ResponseWriter, r *http.Request) {
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

	var req SadBasicRequest
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

	// Extract required fields from decrypted data
	token, ok := decryptedData["token"].(string)
	if !ok || token == "" {
		http.Error(w, "Token not found", http.StatusBadRequest)
		return
	}

	action, ok := decryptedData["action"].(string)
	if !ok || action == "" {
		http.Error(w, "Action is required", http.StatusBadRequest)
		return
	}

	// Extract the data payload
	dataPayload, ok := decryptedData["data"].(map[string]interface{})
	if !ok {
		http.Error(w, "Data payload is required", http.StatusBadRequest)
		return
	}

	// Validate and fix JSON data fields if needed
	fixJSONDataFields(dataPayload)

	// Convert map to JSON and then to MasterSadPayload
	dataJSON, err := json.Marshal(dataPayload)
	if err != nil {
		http.Error(w, "Invalid data format", http.StatusBadRequest)
		return
	}

	var payload databasesad.MasterSadPayload
	if err := json.Unmarshal(dataJSON, &payload); err != nil {
		http.Error(w, fmt.Sprintf("Invalid data structure: %v", err), http.StatusBadRequest)
		return
	}

	// Set token in header for authentication
	r.Header.Set("token", token)

	// 4️⃣ Authentication check
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	// Validate token from request
	if err := auth.IsValidIDFromRequest(r); err != nil {
		http.Error(w, "Invalid TOKEN", http.StatusBadRequest)
		return
	}

	// 5️⃣ Business logic - call master stored procedure
	var resultMsg string
	switch strings.ToLower(action) {
	case "submit", "saveasdraft":
		sadID, err := databasesad.CallMasterSADStoredProcedure(action, payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if strings.ToLower(action) == "submit" {
			resultMsg = fmt.Sprintf("Employee details submitted successfully. SAD ID: %d", sadID)
		} else {
			resultMsg = fmt.Sprintf("Employee details saved as draft successfully. SAD ID: %d", sadID)
		}
	default:
		http.Error(w, "Invalid action. Use 'submit' or 'saveasdraft'.", http.StatusBadRequest)
		return
	}

	response := APIResponseBasicDetails{
		Status:  200,
		Message: resultMsg,
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
}

// Helper function to fix JSON data fields
func fixJSONDataFields(data map[string]interface{}) {
	jsonFields := []string{"contact_data", "dependents_data", "education_data", "experience_data", "language_data", "documents_data"}

	for _, field := range jsonFields {
		if val, exists := data[field]; exists {
			switch v := val.(type) {
			case string:
				// If it's already a string, ensure it's valid JSON
				if v == "" || v == "null" {
					data[field] = nil
				}
			case []interface{}:
				// If it's an array, convert to JSON string
				if len(v) == 0 {
					data[field] = nil
				} else {
					jsonBytes, err := json.Marshal(v)
					if err == nil {
						data[field] = string(jsonBytes)
					}
				}
			case map[string]interface{}:
				// If it's an object, convert to JSON string
				jsonBytes, err := json.Marshal(v)
				if err == nil {
					data[field] = string(jsonBytes)
				}
			case nil:
				// Already nil, leave as is
			default:
				// Unknown type, set to nil
				data[field] = nil
			}
		}
	}
}
