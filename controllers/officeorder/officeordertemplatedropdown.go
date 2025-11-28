// Package controllersofficeorder handles HTTP APIs for officeorder status dropdown.
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
	databaseofficeorder "Hrmodule/database/officeorder"
	"Hrmodule/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// APIResponseDropdownValues defines the structure of the API response
type APIResponseDropdownValues struct {
	Status  int         `json:"Status"`
	Message string      `json:"message"`
	Data    interface{} `json:"Data"`
}

// DropdownValueAPIRequest defines the request body for dropdown values
type DropdownValueAPIRequest struct {
	Data string `json:"Data"`
}

// DropdownValuesHandler handles POST requests to fetch dropdown values
func DropdownValuesHandler(w http.ResponseWriter, r *http.Request) {
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

	var req DropdownValueAPIRequest

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

	token, ok := decryptedData["token"].(string)
	if !ok || token == "" {
		http.Error(w, "Token not found", http.StatusBadRequest)
		return
	}
	r.Header.Set("token", token)

	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		// That function already manages its own response.
		return
	}

	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := auth.IsValidIDFromRequest(r); err != nil {
			//encryptAndRespondfordropdown(w, http.StatusBadRequest, "Invalid TOKEN")
			return
		}

		data, totalCount, err := databaseofficeorder.GetDropdownValuesFromDB(decryptedData)
		if err != nil {
			//	encryptAndRespondfordropdown(w, http.StatusInternalServerError, err.Error())
			return
		}

		response := APIResponseDropdownValues{
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
