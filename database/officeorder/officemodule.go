// Package databaseordersubmodule handles DB access for Order Sub Module Master.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 15-09-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 15-09-2025
package databaseofficeorder

import (
	credentials "Hrmodule/dbconfig"
	modelsofficeorder "Hrmodule/models/officeorder"
	"database/sql"
	"fmt"
)

func GetOrderSubModuleFromDB(decryptedData map[string]interface{}) ([]modelsofficeorder.OrderSubModuleStruct, int, error) {
	connectionString := credentials.Getdatabasemeivan()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, 0, fmt.Errorf("DB open error: %v", err)

	}
	defer db.Close()
	// Extract order_type_id from decrypted data
	orderTypeID, ok := decryptedData["order_type_id"].(string)
	if !ok || orderTypeID == "" {
		return nil, 0, fmt.Errorf("missing 'order_type_id' in request data")
	}

	rows, err := db.Query(modelsofficeorder.MyQueryOrderSubModule, orderTypeID)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	data, err := modelsofficeorder.RetrieveOrderSubModule(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("retrieving result failed: %v", err)
	}

	return data, len(data), nil
}
