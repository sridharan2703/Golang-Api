// Package databasecommon handles database connections and queries related to SessionData.
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On: 25-08-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 25-08-2025
package databaselogin

import (
	credentials "Hrmodule/dbconfig"
	modelslogin "Hrmodule/models/login"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// SessionDatadatabase executes query and returns SessionData list
func SessionDatadatabase(decryptedData map[string]interface{}) ([]modelslogin.SessionDataStructure, int, error) {
	// Connection string
	connectionString := credentials.Getdatabasemeivan()

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, 0, fmt.Errorf("DB open error: %v", err)
	}
	defer db.Close()

	SessionID, ok := decryptedData["Session_id"].(string)
	if !ok || SessionID == "" {
		return nil, 0, fmt.Errorf("missing 'Session_id' in request data")
	}

	// Execute query
	rows, err := db.Query(modelslogin.MyQuerySessionData, SessionID)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	// Map results
	sessionDataList, err := modelslogin.RetrieveSessionData(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving data: %v", err)
	}

	return sessionDataList, len(sessionDataList), nil
}
