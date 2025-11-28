// package controllersofficeorder handles fetch data from database and apply to template
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On: 22-10-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 25-10-2025
package controllersofficeorder

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"Hrmodule/auth"
	credentials "Hrmodule/dbconfig"
	"Hrmodule/utils"

	_ "github.com/lib/pq"
)

var PGDB *sql.DB

func init() {
	var err error
	connStr := credentials.Getdatabasemeivan()
	PGDB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("postgres open error: %v", err)
	}
	if err = PGDB.Ping(); err != nil {
		log.Fatalf("postgres ping error: %v", err)
	}
}

type LeaveInfo struct {
	LeaveName string `json:"leavetype"`
	StartDate string `json:"startdate"`
	EndDate   string `json:"enddate"`
}

type TemplateValues struct {
	FacultyName         string
	IDNumber            string
	Department          string
	Designation         string
	Destination         string
	ActivityDetails     string
	Place               string
	Country             string
	StartDate           string
	EndDate             string
	PerDiemAmount       float64
	Leaves              []LeaveInfo
	PrefixHolidays      string
	InterveningHolidays string
	SuffixHolidays      string
	RelievedDate        string
	OrderNo             string
	Orderdate           string
	Original_orderno    string
}

type OfficeOrderRequest struct {
	Data string `json:"Data"`
}

// Utility: clean HTML list tags
func CleanHTMLToNumberedList(html string) string {
	re := regexp.MustCompile(`(?is)<li>(.*?)</li>`)
	items := re.FindAllStringSubmatch(html, -1)
	var lines []string
	for i, match := range items {
		line := strings.TrimSpace(match[1])
		lines = append(lines, fmt.Sprintf("%d. %s", i+1, line))
	}
	if len(lines) > 0 {
		return strings.Join(lines, "\n")
	}
	reAll := regexp.MustCompile(`<[^>]*>`)
	return strings.TrimSpace(reAll.ReplaceAllString(html, ""))
}

// normalizeClaimType removes extra spaces and normalizes the claim type string
func normalizeClaimType(claimType string) string {
	// Split by comma, trim each part, rejoin with comma (no spaces)
	parts := strings.Split(claimType, ",")
	var normalized []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			normalized = append(normalized, trimmed)
		}
	}
	return strings.Join(normalized, ",")
}

// findTemplateForAmendmentCancellation tries multiple template name variations
func findTemplateForAmendmentCancellation(db *sql.DB, processType, claimType string, tmpl *struct {
	ToColumn, Subject, Reference, BodyHTML, FooterHTML sql.NullString
}) error {
	// Normalize the claim type (remove extra spaces)
	normalizedClaim := normalizeClaimType(claimType)

	// Generate all possible template type variations
	variations := []string{
		// 1. Try with normalized claim type (no spaces after commas)
		fmt.Sprintf("%s,%s", processType, normalizedClaim),
		// 2. Try with original claim type (with spaces)
		fmt.Sprintf("%s,%s", processType, claimType),
	}

	// If there are multiple claim types, try sorted variations
	if strings.Contains(normalizedClaim, ",") {
		parts := strings.Split(normalizedClaim, ",")

		// Try all permutations for common cases (up to 3 items)
		if len(parts) == 2 {
			// Try reverse order: "Project,CPDA" instead of "CPDA,Project"
			variations = append(variations,
				fmt.Sprintf("%s,%s,%s", processType, parts[1], parts[0]),
			)
		} else if len(parts) == 3 {
			// Try common orderings for 3 items
			variations = append(variations,
				fmt.Sprintf("%s,%s,%s,%s", processType, parts[0], parts[1], parts[2]),
				fmt.Sprintf("%s,%s,%s,%s", processType, parts[0], parts[2], parts[1]),
				fmt.Sprintf("%s,%s,%s,%s", processType, parts[1], parts[0], parts[2]),
				fmt.Sprintf("%s,%s,%s,%s", processType, parts[1], parts[2], parts[0]),
				fmt.Sprintf("%s,%s,%s,%s", processType, parts[2], parts[0], parts[1]),
				fmt.Sprintf("%s,%s,%s,%s", processType, parts[2], parts[1], parts[0]),
			)
		}
	}

	// Try each variation
	for _, templateType := range variations {
		fmt.Printf("🔍 Trying template type: '%s'\n", templateType)

		err := db.QueryRow(`
			SELECT to_column, subject, reference, body_html, footer_html
			FROM meivan.processtemplate
			WHERE template_type = $1
			LIMIT 1`, templateType).Scan(
			&tmpl.ToColumn, &tmpl.Subject, &tmpl.Reference, &tmpl.BodyHTML, &tmpl.FooterHTML,
		)

		if err == nil {
			fmt.Printf("✅ Template found with type: '%s'\n", templateType)
			return nil
		}
	}

	// If no exact match found, try case-insensitive and flexible matching
	fmt.Printf("🔍 Trying case-insensitive flexible search for: '%s' + '%s'\n", processType, normalizedClaim)

	query := `
		SELECT to_column, subject, reference, body_html, footer_html
		FROM meivan.processtemplate
		WHERE UPPER(REPLACE(template_type, ' ', '')) = UPPER(REPLACE($1, ' ', ''))
		LIMIT 1`

	templateType := fmt.Sprintf("%s,%s", processType, normalizedClaim)
	err := db.QueryRow(query, templateType).Scan(
		&tmpl.ToColumn, &tmpl.Subject, &tmpl.Reference, &tmpl.BodyHTML, &tmpl.FooterHTML,
	)

	if err == nil {
		fmt.Printf("✅ Template found with flexible matching for: '%s'\n", templateType)
		return nil
	}

	return fmt.Errorf("no template found for processType '%s' and claim_type '%s' (tried %d variations)",
		processType, claimType, len(variations))
}

// ---------------------------
// Date formatting helpers
// ---------------------------

// formatDateValue converts strings like RFC3339 or "yyyy-mm-dd" into "dd-mm-yyyy".
func formatDateValue(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	// Case 1: RFC3339 timestamp (contains "T")
	if strings.Contains(s, "T") {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			return t.Format("02-01-2006")
		}
		// Some DBs have "YYYY-MM-DDTHH:MM:SS" without timezone - try parsing that too
		if t, err := time.Parse("2006-01-02T15:04:05", s); err == nil {
			return t.Format("02-01-2006")
		}
	}

	// Case 2: yyyy-mm-dd
	if len(s) == 10 && strings.Count(s, "-") == 2 {
		if t, err := time.Parse("2006-01-02", s); err == nil {
			return t.Format("02-01-2006")
		}
	}

	// Already dd-mm-yyyy (quick sanity check)
	if len(s) == 10 && strings.Count(s, "-") == 2 {
		parts := strings.Split(s, "-")
		if len(parts[0]) == 2 && len(parts[1]) == 2 && len(parts[2]) == 4 {
			return s
		}
	}

	// Fallback - return original
	return s
}

// convertBodyDates finds dates like "YYYY-MM-DDTHH:MM:SS" or "YYYY-MM-DD" inside HTML and converts to "DD-MM-YYYY".
func convertBodyDates(body string) string {
	if strings.TrimSpace(body) == "" {
		return body
	}

	// 1) Convert full timestamp with T: 2025-09-22T00:00:00 -> 22-09-2025
	re1 := regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2})T\d{2}:\d{2}:\d{2}`)
	body = re1.ReplaceAllString(body, "$3-$2-$1")

	// 2) Convert date-only: 2025-09-22 -> 22-09-2025
	// Use word boundaries to avoid touching numbers that are not dates.
	re2 := regexp.MustCompile(`\b(\d{4})-(\d{2})-(\d{2})\b`)
	body = re2.ReplaceAllString(body, "$3-$2-$1")

	return body
}

// ---------------------------
// Main Handler Function
// ---------------------------
func GetOfficeOrderDetailsfortemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req OfficeOrderRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	parts := strings.Split(req.Data, "||")
	if len(parts) != 2 {
		http.Error(w, "Invalid Data format", http.StatusBadRequest)
		return
	}
	pid := parts[0]
	encryptedPart := parts[1]

	key, err := utils.GetDecryptKey(pid)
	if err != nil {
		http.Error(w, "Decryption key fetch failed", http.StatusUnauthorized)
		return
	}

	decryptedJSON, err := utils.DecryptAES(encryptedPart, key)
	if err != nil {
		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}

	var decryptedData map[string]interface{}
	if err := json.Unmarshal([]byte(decryptedJSON), &decryptedData); err != nil {
		http.Error(w, "Invalid decrypted data", http.StatusBadRequest)
		return
	}

	token, _ := decryptedData["token"].(string)
	employeeID, _ := decryptedData["employeeid"].(string)
	coverPageNo, _ := decryptedData["coverpageno"].(string)
	processType, _ := decryptedData["processtype"].(string)

	if token == "" || employeeID == "" || coverPageNo == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	r.Header.Set("token", token)
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := auth.IsValidIDFromRequest(r); err != nil {
			http.Error(w, "Invalid TOKEN", http.StatusBadRequest)
			return
		}

		processType = strings.ToUpper(strings.TrimSpace(processType))

		var oo struct {
			EmployeeID            string
			FacultyName           string
			Facultynameoriginal   string
			Department            string
			Designation           string
			VisitFrom             string
			VisitTo               string
			Country               string
			CityTown              string
			CoverPageNo           string
			NatureOfParticipation string
			ClaimType             string
			OrderNo               sql.NullString
			Orderdate             string
			Original_orderno      sql.NullString
			LeaveDetails          string
		}

		var tmpl struct {
			ToColumn, Subject, Reference, BodyHTML, FooterHTML sql.NullString
		}

		// ✅ 1️⃣ FETCH BASED ON PROCESSTYPE
		if strings.EqualFold(processType, "Amendment") || strings.EqualFold(processType, "Cancellation") {
			// --- Fetch from meivan.pcr_m ---
			query := `
                 SELECT employee_id, employee_name, department, designation,
                       order_no,order_date,original_order_no, visit_from, visit_to, country, city_town,
                       cover_page_no, nature_of_visit, claim_type
                FROM humanresources.office_order_pcr
                WHERE  employee_id = $1 AND cover_page_no = $2`
			err := PGDB.QueryRow(query, employeeID, coverPageNo).Scan(
				&oo.EmployeeID, &oo.FacultyName, &oo.Department, &oo.Designation,
				&oo.OrderNo, &oo.Orderdate, &oo.Original_orderno, &oo.VisitFrom, &oo.VisitTo, &oo.Country, &oo.CityTown,
				&oo.CoverPageNo, &oo.NatureOfParticipation, &oo.ClaimType,
			)
			if err != nil {
				http.Error(w, fmt.Sprintf("No PCR record found: %v", err), http.StatusInternalServerError)
				return
			}

			// ✅ 2️⃣ Find template using smart matching
			templateFound := false
			processTypeTitle := strings.Title(strings.ToLower(processType))

			// Try to find template with flexible matching
			err = findTemplateForAmendmentCancellation(PGDB, processTypeTitle, oo.ClaimType, &tmpl)
			if err == nil {
				templateFound = true
			}

			if !templateFound {
				http.Error(w, fmt.Sprintf("Template not found for processType '%s' and claim_type '%s': %v", processType, oo.ClaimType, err), http.StatusInternalServerError)
				return
			}

		} else {
			// --- Otherwise fetch from wf_integration.WF_officeorder ---
			query := `
                SELECT 
    w.employeeid, 
    
    -- Prefix suffix value before facultyname
    CASE 
        WHEN c.combovalue IS NOT NULL 
        THEN c.combovalue || ' ' || w.facultyname
        ELSE w.facultyname
    END AS facultyname,
w.facultyname AS facultynameoriginal,
    w.department, 
    w.designation,
    w.visitfrom, 
    w.visitto, 
    w.country, 
    w.citytown, 
    w.coverpageno,
    w.NatureOfParticipation_value, 
    w.claimtype, 
    w.leavedetails

   -- e.suffix,
   -- c.combovalue AS suffix_value

FROM wf_integration.WF_officeorder w

LEFT JOIN humanresources.employeebasicinfo e 
       ON w.employeeid = e.employeeid

LEFT JOIN humanresources.combovaluesmaster c
       ON c.comboname = 'Suffix'
      AND c.displayseq = e.suffix     -- match suffix using displayseq

WHERE w.employeeid=$1 AND w.coverpageno=$2`

			if err := PGDB.QueryRow(query, employeeID, coverPageNo).Scan(
				&oo.EmployeeID, &oo.FacultyName, &oo.Facultynameoriginal, &oo.Department, &oo.Designation,
				&oo.VisitFrom, &oo.VisitTo, &oo.Country, &oo.CityTown,
				&oo.CoverPageNo, &oo.NatureOfParticipation, &oo.ClaimType, &oo.LeaveDetails,
			); err != nil {
				http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
				return
			}

			// For WF_officeorder there is no order_no column — keep as NULL
			oo.OrderNo = sql.NullString{String: "", Valid: false}

			// --- Fetch template by claim_type with flexible matching ---
			normalizedClaim := normalizeClaimType(oo.ClaimType)

			// Try exact match first
			err := PGDB.QueryRow(`
                SELECT to_column, subject, reference, body_html, footer_html
                FROM meivan.processtemplate
                WHERE template_type=$1
                LIMIT 1`, normalizedClaim).Scan(
				&tmpl.ToColumn, &tmpl.Subject, &tmpl.Reference, &tmpl.BodyHTML, &tmpl.FooterHTML,
			)

			// If exact match fails, try case-insensitive without spaces
			if err != nil {
				fmt.Printf("🔍 Exact match failed for '%s', trying flexible matching\n", normalizedClaim)

				err = PGDB.QueryRow(`
					SELECT to_column, subject, reference, body_html, footer_html
					FROM meivan.processtemplate
					WHERE UPPER(REPLACE(template_type, ' ', '')) = UPPER(REPLACE($1, ' ', ''))
					LIMIT 1`, normalizedClaim).Scan(
					&tmpl.ToColumn, &tmpl.Subject, &tmpl.Reference, &tmpl.BodyHTML, &tmpl.FooterHTML,
				)
			}

			if err != nil {
				http.Error(w, fmt.Sprintf("Template not found for claim_type '%s': %v", oo.ClaimType, err), http.StatusInternalServerError)
				return
			}
		}

		// --------------------------
		// Format VisitFrom / VisitTo to dd-mm-yyyy (Option B)
		// --------------------------
		oo.VisitFrom = formatDateValue(oo.VisitFrom)
		oo.VisitTo = formatDateValue(oo.VisitTo)
		oo.Orderdate = formatDateValue(oo.Orderdate)

		// Convert sql.NullString → string safely
		orderNo := ""
		if oo.OrderNo.Valid {
			orderNo = oo.OrderNo.String
		}
		originalOrderNo := ""
		if oo.Original_orderno.Valid {
			originalOrderNo = oo.Original_orderno.String
		}

		// If Original order no is empty, use OrderNo
		if strings.TrimSpace(originalOrderNo) == "" {
			originalOrderNo = orderNo
		}

		// ✅ Prepare template data
		data := TemplateValues{
			FacultyName:      oo.FacultyName,
			IDNumber:         oo.EmployeeID,
			Department:       oo.Department,
			Designation:      oo.Designation,
			Destination:      oo.Country,
			ActivityDetails:  oo.NatureOfParticipation,
			Place:            oo.CityTown,
			Country:          oo.Country,
			StartDate:        oo.VisitFrom,
			EndDate:          oo.VisitTo,
			OrderNo:          orderNo,
			Orderdate:        oo.Orderdate,
			Original_orderno: originalOrderNo,
			Leaves:           []LeaveInfo{},
		}

		// --- Parse leave details and extract TravelExpenseClaimAmount ---
		if oo.LeaveDetails != "" {
			var rawLeaves []map[string]interface{}
			if err := json.Unmarshal([]byte(oo.LeaveDetails), &rawLeaves); err == nil {
				for _, l := range rawLeaves {
					// Check for TravelExpenseClaimAmount
					if claim, ok := l["TravelExpenseClaimAmount"].(float64); ok {
						data.PerDiemAmount = claim
						continue
					}

					li := LeaveInfo{}
					if lv, ok := l["leavetype"].(string); ok {
						li.LeaveName = lv
					}
					if sd, ok := l["startdate"].(string); ok {
						li.StartDate = formatDateValue(sd)
					}
					if ed, ok := l["enddate"].(string); ok {
						li.EndDate = formatDateValue(ed)
					}
					data.Leaves = append(data.Leaves, li)
				}
			}
		}

		// ✅ Execute templates
		var buf = map[string]*bytes.Buffer{
			"Subject":   new(bytes.Buffer),
			"Reference": new(bytes.Buffer),
			"Body":      new(bytes.Buffer),
			"ToColumn":  new(bytes.Buffer),
			"Footer":    new(bytes.Buffer),
		}
		parse := func(name string, tmplStr sql.NullString) {
			if tmplStr.Valid {
				if t, err := template.New(name).Parse(tmplStr.String); err == nil {
					_ = t.Execute(buf[name], data)
				}
			}
		}
		parse("Subject", tmpl.Subject)
		var referenceRendered string
		if tmpl.Reference.Valid {
			if t, err := template.New("Reference").Parse(tmpl.Reference.String); err == nil {
				var tmpBuf bytes.Buffer
				if err := t.Execute(&tmpBuf, data); err == nil {
					referenceRendered = tmpBuf.String()
				}
			}
		}
		referenceRendered = CleanHTMLToNumberedList(referenceRendered)
		parse("Body", tmpl.BodyHTML)
		parse("ToColumn", tmpl.ToColumn)
		if tmpl.FooterHTML.Valid {
			buf["Footer"].WriteString(tmpl.FooterHTML.String)
		}

		// --------------------------
		// Convert dates inside Body HTML to dd-mm-yyyy (Option B)
		// --------------------------
		bodyHTML := buf["Body"].String()
		bodyHTML = convertBodyDates(bodyHTML)
		footerHTML := buf["Footer"].String()
		footerHTML = convertBodyDates(footerHTML)
		referenceRendered = convertBodyDates(referenceRendered)
		toColumn := buf["ToColumn"].String()
		toColumn = convertBodyDates(toColumn)

		// ✅ Build response record
		record := map[string]interface{}{
			"Employeeid":          oo.EmployeeID,
			"Employeename":        oo.FacultyName,
			"Employeecorrectname": oo.Facultynameoriginal,
			"Department":          oo.Department,
			"Designation":         oo.Designation,
			"NatureOfVisit":       oo.NatureOfParticipation,
			"VisitFrom":           oo.VisitFrom,
			"VisitTo":             oo.VisitTo,
			"Country":             oo.Country,
			"CityTown":            oo.CityTown,
			"ClaimType":           oo.ClaimType,
			"OrderNo":             orderNo,
			"Orderdate":           oo.Orderdate,
			"Original_orderno":    originalOrderNo,
			"ToColumn":            toColumn,
			"Subject":             buf["Subject"].String(),
			"Reference":           referenceRendered,
			"Body":                bodyHTML,
			"Footer":              footerHTML,
			"Leaves":              data.Leaves,
			"PerDiemAmount":       data.PerDiemAmount,
		}

		dataBlock := map[string]interface{}{
			"No Of Records": 1,
			"Records":       []interface{}{record},
		}

		innerResponse := map[string]interface{}{
			"P_id":    pid,
			"Status":  200,
			"message": "Success",
			"Data":    dataBlock,
		}

		innerJSON, _ := json.Marshal(innerResponse)
		encryptedData, err := utils.EncryptAES(string(innerJSON), key)
		if err != nil {
			http.Error(w, "Response encryption failed", http.StatusInternalServerError)
			return
		}

		finalResp := map[string]string{
			"Data": fmt.Sprintf("%s||%s", pid, encryptedData),
		}

		auth.SaveResponseLog(r, finalResp, http.StatusOK, "application/json", len(innerJSON), string(body))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(finalResp)
	})).ServeHTTP(w, r)
}
