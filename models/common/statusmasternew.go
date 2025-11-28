// Package modelsstatusmaster contains data structures and queries
// related to Status Master.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 29-09-2025
// Last Modified By: Sridharan
// Last Modified Date: 29-09-2025
package modelscommon

import (
	"database/sql"
	"fmt"
)

var MyQueryStatusMasternew = (`
SELECT statusid, statusdescription
FROM meivan.statusmaster
WHERE statusdescription = $1
`)

type StatusMasternewStruct struct {
	StatusID          *int    `json:"statusid"`
	StatusDescription *string `json:"statusdescription"`
}

func RetrieveStatusMasternew(rows *sql.Rows) ([]StatusMasternewStruct, error) {
	var list []StatusMasternewStruct
	for rows.Next() {
		var s StatusMasternewStruct
		err := rows.Scan(
			&s.StatusID,
			&s.StatusDescription,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		list = append(list, s)
	}
	return list, nil
}
