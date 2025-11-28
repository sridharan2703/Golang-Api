// Package databaseofficeorder handles database operations for OfficeOrdercount.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 15-09-2025
//
// Last Modified By: Ramya
//
// Last Modified Date: 30-09-2025
//
// databaseofficeorder/postgres.go
package databaseofficeorder

import (
	credentials "Hrmodule/dbconfig"
	modelsofficeorder "Hrmodule/models/officeorder"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// GetCombinedNeedGenerate fetches Postgres counts (add MSSQL later if needed)
func GetCombinedNeedGenerate() (modelsofficeorder.CombinedNeedGenerate, error) {
	var result modelsofficeorder.CombinedNeedGenerate

	connectionString := credentials.GetdatabaseWF_officeorder()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return result, fmt.Errorf("DB open error: %v", err)
	}
	defer db.Close()

	rows, err := db.Query(modelsofficeorder.MyQueryNeedGeneratePostgres)
	if err != nil {
		return result, fmt.Errorf("error querying Postgres: %v", err)
	}
	defer rows.Close()

	records, err := modelsofficeorder.RetrieveNeedGeneratePostgres(rows)
	if err != nil {
		return result, fmt.Errorf("error retrieving Postgres data: %v", err)
	}

	if len(records) > 0 {
		result.Postgres = records[0]
	}

	return result, nil
}
