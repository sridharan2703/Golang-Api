// Package controllersofficeorder handles HTTP APIs for ReturnDropdown API.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 29-10-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 29-10-2025
package controllersofficeorder

import (
	"Hrmodule/auth"
	databaseofficeorder "Hrmodule/database/officeorder"
	"Hrmodule/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type APIResponseReturnDropdown struct {
	Status  int         `json:"Status"`
	Message string      `json:"message"`
	Data    interface{} `json:"Data"`
}

type ReturnDropdownRequest struct {
	Data string `json:"Data"`
}

// ReturnDropdown — Controller handler
func ReturnDropdown(w http.ResponseWriter, r *http.Request) {
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

	var req ReturnDropdownRequest

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

	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := auth.IsValidIDFromRequest(r); err != nil {
			http.Error(w, "Invalid TOKEN", http.StatusBadRequest)
			return
		}

		data, totalCount, err := databaseofficeorder.GetReturnDropdownFromDB(decryptedData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := APIResponseReturnDropdown{
			Status:  200,
			Message: "Success",
			Data: map[string]interface{}{
				"No Of Records": totalCount,
				"Records":       data,
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
