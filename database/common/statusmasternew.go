// Package databasestatusmaster handles DB access for Status Master.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 29-09-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 29-09-2025
package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelscommon "Hrmodule/models/common"
	"database/sql"
	"fmt"
)

func GetStatusMasternewFromDB(decryptedData map[string]interface{}) ([]modelscommon.StatusMasternewStruct, int, error) {
	connectionString := credentials.Getdatabasemeivan()

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, 0, fmt.Errorf("DB open error: %v", err)
	}
	defer db.Close()

	// Extract order_type_id from decrypted data
	StatusDescription, ok := decryptedData["statusdescription"].(string)
	if !ok || StatusDescription == "" {
		return nil, 0, fmt.Errorf("missing 'statusdescription' in request data")
	}

	rows, err := db.Query(modelscommon.MyQueryStatusMasternew, StatusDescription)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	data, err := modelscommon.RetrieveStatusMasternew(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("retrieving result failed: %v", err)
	}

	return data, len(data), nil
}
