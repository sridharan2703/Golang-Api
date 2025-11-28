// Package modelsofficeorder contains structs and queries for OfficeOrder_taskvisitdetails API.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 15-09-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 25-10-2025
package modelsofficeorder

import (
	"database/sql"
	"fmt"
)

// Query to fetch all task details by status description
var MyQueryTaskDetails = `
SELECT 
	id, task_id, process_id, cover_page_no, employee_id, employee_name, department, designation,
	visit_from, visit_to, duration, nature_of_visit, claim_type, city_town, country,
	header_html, order_no, order_date, to_column, subject, reference,
	body_html, signature_html, cc_to, footer_html, assign_to, assigned_role,
	task_status_id, activity_seq_no, is_task_return, is_task_approved, email_flag,
	template_id, reject_flag, reject_role, original_order_no, order_type,
	initiated_by, initiated_on, updated_by, updated_on, created_at, updated_at, SM.status,priority
FROM meivan.pcr_m OM
JOIN meivan.statusmaster SM ON OM.task_status_id = SM.statusid
WHERE  statusdescription = $1
`

// Query to fetch task details by EmployeeID + CoverPageNo
var MyQueryTaskDetailsByEmployeeCoverPage = `
SELECT 
	id, task_id, process_id, cover_page_no, employee_id, employee_name, department, designation,
	visit_from, visit_to, duration, nature_of_visit, claim_type, city_town, country,
	header_html, order_no, order_date, to_column, subject, reference,
	body_html, signature_html, cc_to, footer_html, assign_to, assigned_role,
	task_status_id, activity_seq_no, is_task_return, is_task_approved, email_flag,
	template_id, reject_flag, reject_role, original_order_no, order_type,
	initiated_by, initiated_on, updated_by, updated_on, created_at, updated_at, SM.status,priority
FROM meivan.pcr_m OM
JOIN meivan.statusmaster SM ON OM.task_status_id = SM.statusid
WHERE task_status_id not in ('3') and statusdescription = $1 AND employee_id = $2 AND cover_page_no = $3
`

// Query to fetch task details by EmployeeID + CoverPageNo only (no status filter)
var MyQueryTaskDetailsByEmployeeCoverPageOnly = `
SELECT 
	id, task_id, process_id, cover_page_no, employee_id, employee_name, department, designation,
	visit_from, visit_to, duration, nature_of_visit, claim_type, city_town, country,
	header_html, order_no, order_date, to_column, subject, reference,
	body_html, signature_html, cc_to, footer_html, assign_to, assigned_role,
	task_status_id, activity_seq_no, is_task_return, is_task_approved, email_flag,
	template_id, reject_flag, reject_role, original_order_no, order_type,
	initiated_by, initiated_on, updated_by, updated_on, created_at, updated_at, SM.status,priority
FROM meivan.pcr_m OM
JOIN meivan.statusmaster SM ON OM.task_status_id = SM.statusid
WHERE task_status_id not in ('3') and employee_id = $1 AND cover_page_no = $2
`

var MyQueryTaskDetailsBycompleted = `
SELECT distinct
    hoo.task_id,pm.process_id, hoo.cover_page_no, hoo.employee_id, hoo.employee_name, hoo.department, 
    hoo.designation, hoo.visit_from, hoo.visit_to, hoo.duration, hoo.nature_of_visit, 
    hoo.claim_type, hoo.city_town, hoo.country, hoo.header_html, hoo.order_no, 
    hoo.order_date, hoo.to_column, hoo.subject, hoo.reference, hoo.body_html, 
    hoo.signature_html, hoo.cc_to, hoo.footer_html, hoo.original_order_no, hoo.order_type,
    hoo.initiated_by, hoo.initiated_on, hoo.updated_by, hoo.updated_on, 
    hoo.created_at, hoo.updated_at
FROM humanresources.office_order_pcr hoo
 JOIN meivan.pcr_m pm
        ON pm.cover_page_no = hoo.cover_page_no
WHERE NOT EXISTS (
    SELECT 1
    FROM meivan.pcr_m m2
    WHERE m2.cover_page_no = hoo.cover_page_no
      AND m2.task_status_id in ('6','4','22')   
);
`

// Struct representing a full task record
type TaskDetails struct {
	ID              *int64   `json:"id"`
	TaskID          *string  `json:"task_id"`
	ProcessID       *int     `json:"process_id"`
	CoverPageNo     *string  `json:"cover_page_no"`
	EmployeeID      *string  `json:"employee_id"`
	EmployeeName    *string  `json:"employee_name"`
	Department      *string  `json:"department"`
	Designation     *string  `json:"designation"`
	VisitFrom       *string  `json:"visit_from"`
	VisitTo         *string  `json:"visit_to"`
	Duration        *float64 `json:"duration"`
	NatureOfVisit   *string  `json:"nature_of_visit"`
	ClaimType       *string  `json:"claim_type"`
	CityTown        *string  `json:"city_town"`
	Country         *string  `json:"country"`
	HeaderHTML      *string  `json:"header_html"`
	OrderNo         *string  `json:"order_no"`
	OrderDate       *string  `json:"order_date"`
	ToColumn        *string  `json:"to_column"`
	Subject         *string  `json:"subject"`
	Reference       *string  `json:"reference"`
	BodyHTML        *string  `json:"body_html"`
	SignatureHTML   *string  `json:"signature_html"`
	CCTo            *string  `json:"cc_to"`
	FooterHTML      *string  `json:"footer_html"`
	AssignTo        *string  `json:"assign_to"`
	AssignedRole    *string  `json:"assigned_role"`
	TaskStatusID    *int     `json:"task_status_id"`
	ActivitySeqNo   *int     `json:"activity_seq_no"`
	IsTaskReturn    *bool    `json:"is_task_return"`
	IsTaskApproved  *bool    `json:"is_task_approved"`
	EmailFlag       *bool    `json:"email_flag"`
	TemplateID      *int     `json:"template_id"`
	RejectFlag      *int     `json:"reject_flag"`
	RejectRole      *string  `json:"reject_role"`
	OriginalOrderNo *string  `json:"original_order_no"`
	OrderType       *string  `json:"order_type"`
	InitiatedBy     *string  `json:"initiated_by"`
	InitiatedOn     *string  `json:"initiated_on"`
	UpdatedBy       *string  `json:"updated_by"`
	UpdatedOn       *string  `json:"updated_on"`
	CreatedAt       *string  `json:"created_at"`
	UpdatedAt       *string  `json:"updated_at"`
	Status          *int     `json:"status"`
	Priority        *int     `json:"priority"`
}

// Function to map DB rows to TaskDetails slice
func RetrieveTaskDetails(rows *sql.Rows) ([]TaskDetails, error) {
	var list []TaskDetails
	for rows.Next() {
		var v TaskDetails
		err := rows.Scan(
			&v.ID, &v.TaskID, &v.ProcessID, &v.CoverPageNo, &v.EmployeeID, &v.EmployeeName,
			&v.Department, &v.Designation, &v.VisitFrom, &v.VisitTo, &v.Duration,
			&v.NatureOfVisit, &v.ClaimType, &v.CityTown, &v.Country,
			&v.HeaderHTML, &v.OrderNo, &v.OrderDate, &v.ToColumn, &v.Subject, &v.Reference,
			&v.BodyHTML, &v.SignatureHTML, &v.CCTo, &v.FooterHTML,
			&v.AssignTo, &v.AssignedRole, &v.TaskStatusID, &v.ActivitySeqNo,
			&v.IsTaskReturn, &v.IsTaskApproved, &v.EmailFlag, &v.TemplateID,
			&v.RejectFlag, &v.RejectRole, &v.OriginalOrderNo, &v.OrderType,
			&v.InitiatedBy, &v.InitiatedOn, &v.UpdatedBy, &v.UpdatedOn,
			&v.CreatedAt, &v.UpdatedAt, &v.Status, &v.Priority,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		list = append(list, v)
	}
	return list, nil
}

// Function to map DB rows from MyQueryTaskDetailsBycompleted to TaskDetails slice
func RetrieveCompletedTaskDetails(rows *sql.Rows) ([]TaskDetails, error) {
	var list []TaskDetails
	for rows.Next() {
		var v TaskDetails
		err := rows.Scan(
			&v.TaskID, &v.ProcessID, &v.CoverPageNo, &v.EmployeeID, &v.EmployeeName,
			&v.Department, &v.Designation, &v.VisitFrom, &v.VisitTo, &v.Duration,
			&v.NatureOfVisit, &v.ClaimType, &v.CityTown, &v.Country,
			&v.HeaderHTML, &v.OrderNo, &v.OrderDate, &v.ToColumn, &v.Subject, &v.Reference,
			&v.BodyHTML, &v.SignatureHTML, &v.CCTo, &v.FooterHTML,
			&v.OriginalOrderNo, &v.OrderType,
			&v.InitiatedBy, &v.InitiatedOn, &v.UpdatedBy, &v.UpdatedOn,
			&v.CreatedAt, &v.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		list = append(list, v)
	}
	return list, nil
}
