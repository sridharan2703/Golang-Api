// Package modelslogin contains data structures and DB scan logic for Session API.
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
package modelslogin

import (
	"database/sql"
	"fmt"
)

var MyQuerySessionDecryptKey = (`
SELECT decryptkey
FROM meivan.session_data
WHERE is_active=1 AND session_id = $1
`)

// SessionDecryptKey defines structure for session_data table response
type SessionDecryptKey struct {
	DecryptKey *string `json:"decryptkey"`
}

// RetrieveSessionDecryptKey scans rows into []SessionDecryptKey
func RetrieveSessionDecryptKey(rows *sql.Rows) ([]SessionDecryptKey, error) {
	var result []SessionDecryptKey

	for rows.Next() {
		var s SessionDecryptKey
		err := rows.Scan(
			&s.DecryptKey,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		result = append(result, s)
	}
	return result, nil
}
