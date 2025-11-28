// Package modelsofficeorder contains structs and queries for OfficeOrder_modules API.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 15-09-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 15-09-2025
package modelsofficeorder

import (
	"database/sql"
	"fmt"
)

var MyQueryOrderSubModule = (`
SELECT module_id, process_code, process_name, description
FROM meivan.process_master
WHERE id = $1
`)

type OrderSubModuleStruct struct {
	Module_id    *int    `json:"module_id"`
	Process_code *string `json:"process_code"`
	Process_name *string `json:"process_name"`
	Description  *string `json:"description"`
}

func RetrieveOrderSubModule(rows *sql.Rows) ([]OrderSubModuleStruct, error) {
	var list []OrderSubModuleStruct
	for rows.Next() {
		var s OrderSubModuleStruct
		err := rows.Scan(
			&s.Module_id,
			&s.Process_code,
			&s.Process_name,
			&s.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		list = append(list, s)
	}
	return list, nil
}
