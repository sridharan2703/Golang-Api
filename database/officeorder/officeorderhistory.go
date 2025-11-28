// Package databaseofficeorder handles DB access for History of officeorder API.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 29-10-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 29-10-2025
package databaseofficeorder

import (
	credentials "Hrmodule/dbconfig"
	modelsofficeorder "Hrmodule/models/officeorder"
	"database/sql"
	"fmt"
)

func GetOrderHistoryFromDB(decryptedData map[string]interface{}) ([]modelsofficeorder.OrderHistoryStruct, int, error) {
	connectionString := credentials.Getdatabasemeivan()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, 0, fmt.Errorf("DB open error: %v", err)
	}
	defer db.Close()
	// Extract order_type_id from decrypted data
	orderNo, ok := decryptedData["orderNo"].(string)
	if !ok || orderNo == "" {
		return nil, 0, fmt.Errorf("missing 'orderNo' in request data")
	}
	rows, err := db.Query(modelsofficeorder.MyQueryOrderHistory, orderNo)
	if err != nil {
		return nil, 0, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	data, err := modelsofficeorder.RetrieveOrderHistory(rows)
	if err != nil {
		return nil, 0, err
	}

	return data, len(data), nil
}
