// Package models contains data structures and database access logic for the DefaultRoleName page.
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On:30-07-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 10-11-2025
//
// Path:Login Page
package modelscommon

import (
	//	modelstable "Hrmodule/models/tables"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var MyQueryDefaultRoleName = `
SELECT A.Employeeid as UserID,loginname as Username,
    B.campuscode || ' ' ||
    CASE WHEN A.sectionid IS NOT NULL THEN D.sectioncode ELSE A.departmentcode END || ' ' || 
    C.rolename AS RoleName
FROM humanresources.employeerolemapping A
JOIN humanresources.campus B ON A.campusid = B.id
JOIN meivan.rolemaster C ON A.roleid = C.id
LEFT JOIN humanresources.section D ON A.sectionid = D.id
join humanresources.employeebasicinfo e
on A.employeeid=e.employeeid
WHERE e.loginname = $1
`

// DefaultRoleNamestructure defines the structure of DefaultRoleName
type DefaultRoleNamestructure struct {
	UserID   *string `json:"UserID"`
	Username *string `json:"Username"`
	RoleName *string `json:"RoleName"`
}

// RetrieveDefaultRoleName scans rows into DefaultRoleNamestructure slice
func RetrieveDefaultRoleName(rows *sql.Rows) ([]DefaultRoleNamestructure, error) {
	var DefaultRoleNameapi []DefaultRoleNamestructure

	for rows.Next() {
		var DRN DefaultRoleNamestructure
		err := rows.Scan(
			&DRN.UserID,
			&DRN.Username,
			&DRN.RoleName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		DefaultRoleNameapi = append(DefaultRoleNameapi, DRN)
	}

	return DefaultRoleNameapi, nil
}
