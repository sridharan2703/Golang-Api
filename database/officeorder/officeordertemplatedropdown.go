// Package databaseofficeorder handles database operations for status dropdown.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 15-09-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 30-09-2025
package databaseofficeorder

import (
	credentials "Hrmodule/dbconfig"
	modelsofficeorder "Hrmodule/models/officeorder"
	"database/sql"
	"fmt"
)

func GetDropdownValuesFromDB(decryptedData map[string]interface{}) ([]modelsofficeorder.DropdownValueStruct, int, error) {
	connectionString := credentials.GetGnanaThalamSchemaConnection("meivan")
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, 0, fmt.Errorf("DB open error: %v", err)
	}
	defer db.Close()
	// Extract order_type_id from decrypted data
	CoverPageNo, ok := decryptedData["coverpageno"].(string)
	if !ok || CoverPageNo == "" {
		return nil, 0, fmt.Errorf("missing 'coverpageno' in request data")
	}
	EmployeeID, ok := decryptedData["employeeid"].(string)
	if !ok || EmployeeID == "" {
		return nil, 0, fmt.Errorf("missing 'employeeid' in request data")
	}

	rows, err := db.Query(modelsofficeorder.MyQueryDropdownValues, CoverPageNo, EmployeeID)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	data, err := modelsofficeorder.RetrieveDropdownValues(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("retrieving result failed: %v", err)
	}

	return data, len(data), nil
}
