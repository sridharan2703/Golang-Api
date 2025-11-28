// Package modelsofficeorder contains structs and queries for CC Roles API.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 20-11-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 20-11-2025
package modelsofficeorder

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

var MyQueryCcRoles = (`
SELECT 
    role_name,
    ARRAY_AGG(email_id ORDER BY email_id) AS email_list,
    ARRAY_AGG(name ORDER BY name) AS name_list,
    ARRAY_AGG(employee_id ORDER BY employee_id) AS employee_list
FROM (
    SELECT DISTINCT 
        r.name,
        r.email_id,
        r.role_name,
        r.employee_id,
        CASE 
            WHEN r.role_name = d.departmentcode || ' HOD' THEN 1
            WHEN r.role_name ILIKE '%DEAN%' THEN 2
            WHEN r.role_name ILIKE '%DEPUTY REGISTRAR%' THEN 3
            WHEN r.role_name = 'ADMIN1 ASSISTANT REGISTRAR' THEN 4
            WHEN r.role_name = 'ADMIN2 ASSISTANT REGISTRAR' THEN 5
            ELSE 99
        END AS role_priority
    FROM humanresources.employeeappointmentdetails e
    JOIN humanresources.departmentmaster d
        ON UPPER(e.deptcode) = UPPER(d.departmentcode)
    JOIN meivan.cc_to_recipients r
        ON TRUE
    WHERE e.employeeid =$1
      AND (
            r.role_name = d.departmentcode || ' HOD'
         OR r.role_name ILIKE '%DEAN%'
         OR r.role_name ILIKE '%REGISTRAR%'
      )
) AS t
GROUP BY role_name, role_priority
ORDER BY role_priority, role_name;
`)

type CcRoleStruct struct {
	RoleName     *string  `json:"role_name"`
	EmailList    []string `json:"email_list"`
	NameList     []string `json:"name_list"`
	EmployeeList []string `json:"employee_list"`
}

func RetrieveCcRoles(rows *sql.Rows) ([]CcRoleStruct, error) {
	var list []CcRoleStruct

	for rows.Next() {
		var s CcRoleStruct

		err := rows.Scan(
			&s.RoleName,
			pq.Array(&s.EmailList),
			pq.Array(&s.NameList),
			pq.Array(&s.EmployeeList),
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		list = append(list, s)
	}
	return list, nil
}
