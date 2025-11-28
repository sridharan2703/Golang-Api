// Package databaselogin handles DB calls for SessionDecryptKey API.
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On: 15-09-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 15-09-2025
package databaselogin

import (
	credentials "Hrmodule/dbconfig"
	modelslogin "Hrmodule/models/login"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

// Request body for SessionDecryptKey
type SessionDecryptKeyRequest struct {
	SessionID string `json:"session_id"`
}

// SessionDecryptKeyDatabase executes the query
func SessionDecryptKeyDatabase(w http.ResponseWriter, r *http.Request) ([]modelslogin.SessionDecryptKey, int, error) {
	connectionString := credentials.Getdatabasemeivan()

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		http.Error(w, fmt.Sprintf("DB open error: %v", err), http.StatusInternalServerError)
		return nil, 0, err
	}
	defer db.Close()

	// Decode request body
	var req SessionDecryptKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, 0, fmt.Errorf("invalid request body: %v", err)
	}
	defer r.Body.Close()

	if req.SessionID == "" {
		return nil, 0, fmt.Errorf("missing 'session_id' in request body")
	}

	// Execute query
	rows, err := db.Query(modelslogin.MyQuerySessionDecryptKey, req.SessionID)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying DB: %v", err)
	}
	defer rows.Close()

	// Map results
	data, err := modelslogin.RetrieveSessionDecryptKey(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving data: %v", err)
	}

	return data, len(data), nil
}
