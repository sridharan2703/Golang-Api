// // Package modelsofficeorder contains structs and queries for approval page OfficeOrder_approval_remarks API.
// //
// // --- Creator's Info ---
// // Creator: Sridharan
// //
// // Created On: 15-09-2025
// //
// // Last Modified By: Sridharan
// //
// // Last Modified Date: 15-09-2025
// package modelsofficeorder

// import (
// 	"database/sql"
// 	"fmt"
// 	"time" // Needed for the 'updatedon' column type
// )
// var MyQueryOfficeComments = (`
// 	SELECT
// 		C.remarks,
// 		C.updatedby,
// 		C.updatedon
// 	FROM
// 		meivan.officeorder_comments AS C
// 	INNER JOIN
// 		meivan.officeorder_master AS M
// 	ON
// 		C.officeorderid = M.officeorderid
// 	WHERE
// 		M.employeeid = $1 AND M.coverpageno = $2;
// `)

//  type OfficeCommentStructure struct {
// 	Remarks   *string    `json:"Remarks"`    // "remarks" (text)
// 	UpdatedBy *string    `json:"UpdatedBy"`  // "updatedby" (text)
// 	UpdatedOn *time.Time `json:"UpdatedOn"`  // "updatedon" (timestamp with timezone)
// }

//  func RetrieveOfficeComments(rows *sql.Rows) ([]OfficeCommentStructure, error) {
// 	var commentList []OfficeCommentStructure

// 	for rows.Next() {
// 		var comment OfficeCommentStructure

// 		    err := rows.Scan(
// 			&comment.Remarks,
// 			&comment.UpdatedBy,
// 			&comment.UpdatedOn,
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("error scanning OfficeCommentStructure row: %v", err)
// 		}
// 		commentList = append(commentList, comment)
// 	}

// 	return commentList, nil
// }

// Package modelsofficeorder contains structs and queries for approval page OfficeOrder_approval_remarks API.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 15-09-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 27-10-2025


// package modelsofficeorder

// import (
// 	"database/sql"
// 	"fmt"
// 	"time"
// )

// // --- Queries for fetching remarks based on inputs ---

// // Case 1: Using employee_id + cover_page_no
// var QueryOfficeCommentsByEmployee = (`
// 	SELECT 
// 		user_id,
// 		COALESCE(user_role, '') AS user_role,
// 		remarks,
// 		updated_by,
// 		updated_on
// 	FROM 
// 		meivan.pdftemplaterecords($1, $2, NULL)order by updated_on desc ;
// `)

// // Case 2: Using order_no only
// var QueryOfficeCommentsByOrderNo = (`
// 	SELECT 
// 		user_id,
// 		COALESCE(user_role, '') AS user_role,
// 		remarks,
// 		updated_by,
// 		updated_on
// 	FROM 
// 		meivan.pdftemplaterecords(NULL, NULL, $1)order by updated_on desc ;
// `)

// // --- Struct for Result Data ---
// type OfficeCommentStructure struct {
// 	UserID    *string    `json:"UserID"`    // user_id
// 	UserRole  *string    `json:"UserRole"`  // user_role
// 	Remarks   *string    `json:"Remarks"`   // remarks
// 	UpdatedBy *string    `json:"UpdatedBy"` // updated_by
// 	UpdatedOn *time.Time `json:"UpdatedOn"` // updated_on (timestamp with/without timezone)
// }

// // --- Function to read from rows ---
// func RetrieveOfficeComments(rows *sql.Rows) ([]OfficeCommentStructure, error) {
// 	var comments []OfficeCommentStructure

// 	for rows.Next() {
// 		var comment OfficeCommentStructure

// 		err := rows.Scan(
// 			&comment.UserID,
// 			&comment.UserRole,
// 			&comment.Remarks,
// 			&comment.UpdatedBy,
// 			&comment.UpdatedOn,
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("error scanning OfficeCommentStructure row: %v", err)
// 		}

// 		comments = append(comments, comment)
// 	}

// 	return comments, nil
// }

//17/11/2025

package modelsofficeorder

import (
	"database/sql"
	"fmt"
)

// --- Query for fetching comments using process_id + task_id ---
var QueryOfficeComments = `
    SELECT 
        user_display,
        user_role,
        remarks,
        updated_on
    FROM meivan.getcomments($1, $2)
    ORDER BY updated_on DESC;
`

// --- Struct for Result Data ---
type OfficeCommentStructure struct {
	UserDisplay *string `json:"UserID"` // user_display (id - name)
	UserRole    *string `json:"UserRole"`    // user_role
	Remarks     *string `json:"Remarks"`     // remarks
	UpdatedOn   *string `json:"UpdatedOn"`   // updated_on (formatted: YYYY-MM-DD HH:MI)
}

// --- Function to read from rows ---
func RetrieveOfficeComments(rows *sql.Rows) ([]OfficeCommentStructure, error) {
	var comments []OfficeCommentStructure

	for rows.Next() {
		var comment OfficeCommentStructure

		err := rows.Scan(
			&comment.UserDisplay,
			&comment.UserRole,
			&comment.Remarks,
			&comment.UpdatedOn,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning OfficeCommentStructure row: %v", err)
		}

		comments = append(comments, comment)
	}

	return comments, nil
}
