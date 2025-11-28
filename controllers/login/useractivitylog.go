// Package common contains APIs that are  used across the application for save the user activity
//
// This API marks a user session as inactive upon logout or timeout.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 09-07-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 09-07-2025
package controllerslogin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"Hrmodule/auth"
	credentials "Hrmodule/dbconfig"
	"Hrmodule/utils"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// UserActivityLog represents the data for a single log entry.
// The `json.RawMessage` type is used to handle the `event_data` JSONB column
// without needing to know its exact structure in advance.
type UserActivityLog struct {
	SessionID  string          `json:"session_id"`
	EventName  string          `json:"event_name"`
	EventData  json.RawMessage `json:"event_data"`
	EmployeeID string          `json:"employeeid"`
}

// InsertUserActivityLog handles the HTTP request to insert a new user activity log entry.
func InsertUserActivityLog(w http.ResponseWriter, r *http.Request) {

	// Handle token and IP validation, and log request info using a single wrapper.
	loggedHandler := auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Step 1: Handle authentication. This function should handle IP and token validation.
		if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
			return
		}

		// Step 2: Allow only POST method.
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Step 3: Decode the request body into the UserActivityLog struct.
		var logEntry UserActivityLog
		if err := json.NewDecoder(r.Body).Decode(&logEntry); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Step 4: Connect to the PostgreSQL database.
		connectionString := credentials.Getdatabasemeivan()
		db, err := sql.Open("postgres", connectionString)
		if err != nil {
			http.Error(w, fmt.Sprintf("Database connection error: %v", err), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		// Step 5: Prepare and execute the SQL INSERT statement with the new employeeid column.
		query := `
			INSERT INTO user_activity_logs (
				session_id,
				event_name,
				event_data,
				employeeid
			)
			VALUES ($1, $2, $3, $4)`

		_, err = db.Exec(
			query,
			logEntry.SessionID,
			logEntry.EventName,
			logEntry.EventData,
			logEntry.EmployeeID, // Pass the new field to the query
		)

		if err != nil {
			http.Error(w, "Insert failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Step 6: Prepare the JSON response.
		response := map[string]string{"message": "User activity log inserted successfully"}

		// Step 7: Marshal and encrypt the response.
		responseBytes, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Failed to serialize JSON", http.StatusInternalServerError)
			return
		}
		encrypted, err := utils.Encrypt(responseBytes)
		if err != nil {
			http.Error(w, "Encryption failed", http.StatusInternalServerError)
			return
		}

		// Step 8: Send the final encrypted response.
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"Data": encrypted,
		})
	}))

	// Execute the wrapped handler.
	loggedHandler.ServeHTTP(w, r)
}
