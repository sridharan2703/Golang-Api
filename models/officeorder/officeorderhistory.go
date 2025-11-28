// Package modelsofficeorder contains structs and queries for OfficeOrderHistory
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 14-11-2025
//
// Last Modified By:
//
// Last Modified Date:
package modelsofficeorder

import (
	"database/sql"
	"fmt"
)

var MyQueryOrderHistory = (`
WITH base AS (
    SELECT regexp_replace($1, '/A[0-9]+(/CAN)?$', '') AS base_order
)
SELECT
    order_no,
    task_status_id
FROM meivan.pcr_m m, base b
WHERE m.order_no LIKE b.base_order || '%'
AND m.task_status_id = 3   -- only completed
ORDER BY length(m.order_no), m.order_no
`)

type OrderHistoryStruct struct {
	OrderNo      *string `json:"order_no"`
	TaskStatusID *int    `json:"task_status_id"`
}

func RetrieveOrderHistory(rows *sql.Rows) ([]OrderHistoryStruct, error) {
	var list []OrderHistoryStruct
	for rows.Next() {
		var s OrderHistoryStruct
		err := rows.Scan(
			&s.OrderNo,
			&s.TaskStatusID,
		)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %v", err)
		}
		list = append(list, s)
	}
	return list, nil
}
