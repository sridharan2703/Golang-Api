// Package modelsofficeorder contains structs and queries for OfficeOrder API.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 15-09-2025
// Last Modified By:  Ramya
// Last Modified Date: 30-09-2025
// modelsofficeorder/postgres.go
package modelsofficeorder

import (
	"database/sql"
	"fmt"
)

// -------------------- PostgreSQL --------------------

var MyQueryNeedGeneratePostgres = `
SELECT 
    SUM(CASE WHEN w.status = 4 THEN 1 ELSE 0 END) AS ongoing,
   (
        SELECT COUNT(*)
        FROM humanresources.office_order_pcr hoo
        WHERE NOT EXISTS (
            SELECT 1 
            FROM meivan.pcr_m m2
            WHERE m2.cover_page_no = hoo.cover_page_no
              AND m2.task_status_id in ('6','4','22')
        )
    ) AS complete,

    SUM(CASE WHEN w.status = 6 THEN 1 ELSE 0 END) AS saveandhold,
    SUM(CASE WHEN w.status IN ('0','1') THEN 1 ELSE 0 END) AS needtogenerate
FROM wf_integration.WF_officeorder w;
`

type NeedGeneratePostgres struct {
	Ongoing        int `json:"ongoing"`
	Complete       int `json:"complete"`
	SaveAndHold    int `json:"saveandhold"`
	NeedToGenerate int `json:"need_to_generate"`
}

// CombinedNeedGenerate combines data from multiple DBs (Postgres, MSSQL, etc.)
type CombinedNeedGenerate struct {
	Postgres NeedGeneratePostgres `json:"postgres"`
	// Add other DB results here later
}

func RetrieveNeedGeneratePostgres(rows *sql.Rows) ([]NeedGeneratePostgres, error) {
	var list []NeedGeneratePostgres
	for rows.Next() {
		var n NeedGeneratePostgres
		err := rows.Scan(
			&n.Ongoing,
			&n.Complete,
			&n.SaveAndHold,
			&n.NeedToGenerate,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning Postgres row: %v", err)
		}
		list = append(list, n)
	}
	return list, nil
}
