// Package databasecommon handles DB calls for InboxTasksRole API.
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On: 26-08-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 25-10-2025
package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelscommon "Hrmodule/models/common"
	"database/sql"
	"fmt"
	"net/url"
	"strings"
)

// InboxTasksRoleDatabase executes getinboxtasks_role
func InboxTasksRoleDatabase(decryptedData map[string]interface{}) ([]modelscommon.InboxTasksRole, int, error) {
	connectionString := credentials.Getdatabasemeivan()

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, 0, fmt.Errorf("DB open error: %v", err)
	}
	defer db.Close()

	// Extract order_type_id from decrypted data
	EmpID, ok := decryptedData["empid"].(string)
	if !ok || EmpID == "" {
		//return nil, 0, fmt.Errorf("missing 'empid' in request data")
	}

	// Extract order_type_id from decrypted data
	AssignedRole, ok := decryptedData["assignedrole"].(string)
	if !ok || AssignedRole == "" {
		//return nil, 0, fmt.Errorf("missing 'assignedrole' in request data")
	}
	// ✅ Decode any URL-encoded values like %20 → space
	if decodedRole, err := url.QueryUnescape(AssignedRole); err == nil {
		AssignedRole = decodedRole
	} else {
		// fallback: just replace %20 with space
		AssignedRole = strings.ReplaceAll(AssignedRole, "%20", " ")
	}

	// Run query
	rows, err := db.Query(modelscommon.MyQueryInboxTasksRole, EmpID, AssignedRole)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying DB: %v", err)
	}
	defer rows.Close()

	// Map results
	data, err := modelscommon.RetrieveInboxTasksRole(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving data: %v", err)
	}

	return data, len(data), nil
}
