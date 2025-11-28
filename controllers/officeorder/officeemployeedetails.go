// Package controllersofficeorder handles HTTP APIs for officeorder visit details.
//
// --- Creator's Info ---
// Creator: Sridharan
// Created On: 15-09-2025
// Last Modified By: Sridharan
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

// APIResponseVisitDetails defines the structure of API responses
type APIResponseVisitDetails struct {
	Status  int         `json:"Status"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data,omitempty"`
}

// VisitDetailsRequest defines incoming encrypted payload format
type VisitDetailsRequest struct {
	Data string `json:"Data"`
}

// -----------------------------------------------------------------------------
// ✅ Main Handler: Ordervisitdetails
// -----------------------------------------------------------------------------
func Ordervisitdetails(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		encryptAndRespondforvisitdetails(w, r, nil, APIResponseVisitDetails{
			Status:  http.StatusMethodNotAllowed,
			Message: "Method not allowed, use POST",
		})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		encryptAndRespondforvisitdetails(w, r, body, APIResponseVisitDetails{
			Status:  http.StatusBadRequest,
			Message: "Unable to read request body",
		})
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	var req VisitDetailsRequest
	
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Split Data → pid || encryptedPart
	parts := strings.Split(req.Data, "||")
	if len(parts) != 2 {
		http.Error(w, "Invalid Data format", http.StatusBadRequest)
		return
	}
	pid := parts[0]
	encryptedPart := parts[1]

	key, err := utils.GetDecryptKey(pid)
	if err != nil {
		http.Error(w, "Failed to get decryption key", http.StatusUnauthorized)
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

	token, ok := decryptedData["token"].(string)
	if !ok || token == "" {
		http.Error(w, "Token not found", http.StatusBadRequest)
		return
	}

	r.Header.Set("token", token)

	// ✅ Authentication
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		encryptAndRespondforvisitdetails(w, r, body, APIResponseVisitDetails{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized access - token validation failed",
		})
		return
	}

	// ✅ Log and process
	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := auth.IsValidIDFromRequest(r); err != nil {
			http.Error(w, "Invalid TOKEN provided", http.StatusBadRequest)
			return
		}

		data, totalCount, err := database.GetVisitDetailsFromDB(decryptedData)
		if err != nil {
			encryptAndRespondforvisitdetails(w, r, body, APIResponseVisitDetails{
				Status:  http.StatusInternalServerError,
				Message: "Database error: " + err.Error(),
			})
			return
		}

		response := APIResponseVisitDetails{
			Status:  http.StatusOK,
			Message: "Success",
			Data: map[string]interface{}{
				"No Of Records": totalCount,
				"Records":       data,
			},
		}

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

		auth.SaveResponseLog(r, finalResp, http.StatusOK, "application/json", len(responseJSON), string(body))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(finalResp)
	})).ServeHTTP(w, r)
}

// -----------------------------------------------------------------------------
// ✅ Helper: encryptAndRespondforvisitdetails
// -----------------------------------------------------------------------------
func encryptAndRespondforvisitdetails(w http.ResponseWriter, r *http.Request, body []byte, resp APIResponseVisitDetails) {
	// Try to extract pid and key again (if present)
	var pid, key string

	if len(body) > 0 {
		var req VisitDetailsRequest
		if err := json.Unmarshal(body, &req); err == nil {
			parts := strings.Split(req.Data, "||")
			if len(parts) == 2 {
				pid = parts[0]
				if k, err := utils.GetDecryptKey(pid); err == nil {
					key = k
				}
			}
		}
	}

	responseJSON, _ := json.MarshalIndent(resp, "", "  ")

	// Encrypt only if key is available
	encryptedResponse := string(responseJSON)
	if key != "" {
		enc, err := utils.EncryptAES(string(responseJSON), key)
		if err == nil {
			encryptedResponse = fmt.Sprintf("%s||%s", pid, enc)
		}
	}

	finalResp := map[string]string{
		"Data": encryptedResponse,
	}

	auth.SaveResponseLog(r, finalResp, resp.Status, "application/json", len(responseJSON), string(body))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)
	json.NewEncoder(w).Encode(finalResp)
}
