// Package modelsofficeorder contains structs and queries for approval page status dropdown API.
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

// SQL Query for dropdown values
var MyQueryDropdownValues = (`
WITH Latest AS (
    SELECT *
    FROM meivan.pcr_m
    WHERE cover_page_no = $1
      AND employee_id = $2

    ORDER BY id DESC
    LIMIT 1
),
WFStatus AS (
    SELECT *
    FROM wf_integration.WF_officeorder
    WHERE coverpageno = $1
      AND employeeid = $2

)

SELECT 
    ROW_NUMBER() OVER (ORDER BY sort_order) AS id,
    dropdown_value
FROM (
    -- Completed
    SELECT 
        '' AS dropdown_value,
        1 AS sort_order
    FROM Latest L
    WHERE L.order_no LIKE '%/CAN%' and L.task_status_id = 3

    UNION ALL

    -- Amendment
    SELECT 
        'Amendment' AS dropdown_value,
        2 AS sort_order
    FROM Latest L
    WHERE L.task_status_id = 3
      AND L.order_no NOT LIKE '%/CAN%'
      --AND  L.order_no ~ '/A[0-9]+$'

    UNION ALL

    -- Cancellation
    SELECT 
        'Cancellation' AS dropdown_value,
        3 AS sort_order
    FROM Latest L
    WHERE L.task_status_id = 3
      AND L.order_no NOT LIKE '%/CAN%'
      --AND  L.order_no ~ '/A[0-9]+$'

    UNION ALL

    -- Continue Editing
    SELECT 
        'Continue Editing' AS dropdown_value,
        4 AS sort_order
    FROM Latest L
    WHERE L.task_status_id = 6

    UNION ALL

    -- New Order from task_status_id = 6
    SELECT 
        'New Order' AS dropdown_value,
        5 AS sort_order
    FROM Latest L
    WHERE L.task_status_id = 6

    UNION ALL

    -- New Order from pcr_m (task_status_id 0 or 1)
    SELECT 
        'New Order' AS dropdown_value,
        6 AS sort_order
    FROM Latest L
    WHERE L.task_status_id IN (0,1)

    UNION ALL

    -- New Order from WF_officeorder (status = 0)
    SELECT 
        'New Order' AS dropdown_value,
        7 AS sort_order
    FROM WFStatus W
    WHERE W.status = 0
) AS dropdowns
ORDER BY sort_order;

`)

// Struct for single dropdown value
type DropdownValueStruct struct {
	Id            *string `json:"Id"`
	DropdownValue *string `json:"dropdown_value"`
}

// Function to retrieve dropdown values from rows
func RetrieveDropdownValues(rows *sql.Rows) ([]DropdownValueStruct, error) {
	var list []DropdownValueStruct
	for rows.Next() {
		var d DropdownValueStruct
		err := rows.Scan(&d.Id, &d.DropdownValue)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		list = append(list, d)
	}
	return list, nil
}
