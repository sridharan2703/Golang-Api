// Package databaseofficeorder handles DB access for CC Roles Master.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 20-11-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 20-11-2025
package databaseofficeorder

import (
	credentials "Hrmodule/dbconfig"
	modelsofficeorder "Hrmodule/models/officeorder"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func GetCcRolesFromDB(decryptedData map[string]interface{}) ([]modelsofficeorder.CcRoleStruct, int, error) {

	connectionString := credentials.Getdatabasemeivan()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, 0, fmt.Errorf("DB open error: %v", err)
	}
	defer db.Close()

	employeeID, ok := decryptedData["employee_id"].(string)
	if !ok || employeeID == "" {
		return nil, 0, fmt.Errorf("missing 'employee_id' in request data")
	}

	rows, err := db.Query(modelsofficeorder.MyQueryCcRoles, employeeID)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	data, err := modelsofficeorder.RetrieveCcRoles(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("retrieving result failed: %v", err)
	}

	return data, len(data), nil
}
