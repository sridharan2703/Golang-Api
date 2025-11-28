// Package controllersofficeorder handles HTTP APIs for officeordertaskvisitdetails.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 15-09-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 25-10-2025
package controllersofficeorder

import (
	"Hrmodule/auth"
	database "Hrmodule/database/officeorder"
	"Hrmodule/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// -----------------------------------------------------------------------------
// Response structure
// -----------------------------------------------------------------------------
type APIResponseTaskVisitDetails struct {
	Status  int         `json:"Status"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data,omitempty"`
}

// Request structure
type TaskVisitDetailsRequest struct {
	Data string `json:"Data"`
}

// -----------------------------------------------------------------------------
// ✅ Main Handler: OrderTaskVisitDetails
// -----------------------------------------------------------------------------
func OrderTaskVisitDetails(w http.ResponseWriter, r *http.Request) {
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
	r.Body = io.NopCloser(bytes.NewBuffer(body)) // Reusable
	defer r.Body.Close()

	// Step 3️⃣: Parse JSON
	var req TaskVisitDetailsRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Step 4️⃣: Split and decrypt
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

	// Step 5️⃣: Validate token
	token, ok := decryptedData["token"].(string)
	if !ok || token == "" {
		http.Error(w, "Token not found", http.StatusBadRequest)
		return
	}
	r.Header.Set("token", token)

	// Step 6️⃣: Authenticate request
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		sendEncryptedResponse(w, r, pid, key, body, APIResponseTaskVisitDetails{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized access - token validation failed",
		})
		return
	}

	// Step 7️⃣: Process inside logged handler
	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := auth.IsValidIDFromRequest(r); err != nil {
			sendEncryptedResponse(w, r, pid, key, body, APIResponseTaskVisitDetails{
				Status:  http.StatusBadRequest,
				Message: "Invalid TOKEN provided",
			})
			return
		}

		// ✅ Fetch from DB
		data, totalCount, err := database.GetTaskVisitDetailsFromDB(decryptedData)
		if err != nil {
			sendEncryptedResponse(w, r, pid, key, body, APIResponseTaskVisitDetails{
				Status:  http.StatusInternalServerError,
				Message: "Database error: " + err.Error(),
			})
			return
		}

		response := APIResponseTaskVisitDetails{
			Status:  http.StatusOK,
			Message: "Success",
			Data: map[string]interface{}{
				"No Of Records": totalCount,
				"Records":       data,
			},
		}

		sendEncryptedResponse(w, r, pid, key, body, response)
	})).ServeHTTP(w, r)
}

// -----------------------------------------------------------------------------
// ✅ Local Helper: sendEncryptedResponse
// -----------------------------------------------------------------------------
func sendEncryptedResponse(
	w http.ResponseWriter,
	r *http.Request,
	pid string,
	key string,
	body []byte,
	resp APIResponseTaskVisitDetails,
) {
	responseJSON, err := json.MarshalIndent(resp, "", "  ")
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

	// ✅ Log exactly what’s sent
	auth.SaveResponseLog(
		r,
		finalResp,
		resp.Status,
		"application/json",
		len(responseJSON),
		string(body),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)
	json.NewEncoder(w).Encode(finalResp)
}
