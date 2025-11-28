// Package cronjob handles HTTP APIs for officeorder sync from mssql to postgresql.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 15-09-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 15-09-2025
// SyncData from mssql to postgresql
package cronjob



import (
	"context"
	"database/sql"
	"encoding/json"

	//"fmt"
	"log"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/lib/pq"

	//  "github.com/robfig/cron/v3"

	credentials "Hrmodule/dbconfig"
)

// Structs for rows
type OfficeOrder struct {
	EmployeeID                 sql.NullString `json:"EmployeeID"`
	FacultyName                sql.NullString `json:"FacultyName"`
	Department                 sql.NullString `json:"Department"`
	Designation                sql.NullString `json:"Designation"`
	FacultyDetails             sql.NullString `json:"FacultyDetails"`
	VisitFrom                  sql.NullTime   `json:"VisitFrom"`
	VisitTo                    sql.NullTime   `json:"VisitTo"`
	Country                    sql.NullString `json:"Country"`
	CityTown                   sql.NullString `json:"CityTown"`
	CoverPageNo                sql.NullString `json:"CoverPageNo"`
	NatureOfParticipation      sql.NullString `json:"NatureOfParticipation"`
	NatureOfParticipationValue sql.NullString `json:"NatureOfParticipation_Value"`
	InitiatedOn                sql.NullTime   `json:"InitiatedOn"`
	ClaimType                  sql.NullString `json:"ClaimType"`
	LeaveDetails               sql.NullString `json:"LeaveDetails"`
}

func nullStringToInterface(s sql.NullString) interface{} {
	if s.Valid {
		return s.String
	}
	return nil
}
func nullTimeToInterface(t sql.NullTime) interface{} {
	if t.Valid {
		return t.Time
	}
	return nil
}

// Sync function
func SyncData() {
	log.Println("🔄 Sync started...")

	mssqlConnStr := credentials.GetLivedatabase10()
	pgConnStr := credentials.GetdatabaseWF_officeorder()

	// Connect MSSQL
	mssqlDB, err := sql.Open("mssql", mssqlConnStr)
	if err != nil {
		log.Println("MSSQL connect error:", err)
		return
	}
	defer mssqlDB.Close()

	// Connect PostgreSQL
	pgDB, err := sql.Open("postgres", pgConnStr)
	if err != nil {
		log.Println("PostgreSQL connect error:", err)
		return
	}
	defer pgDB.Close()

	query := `
	WITH ClaimAggregates AS (
    SELECT
        DistinctClaims.InstanceId,
        STRING_AGG(DistinctClaims.ClaimType, ', ') AS ClaimTypes
    FROM (
        SELECT DISTINCT InstanceId, ClaimType
        FROM Proof..[Client_BAF56E7B-1074-4997-AEB7-C985D6AC770E]
    ) AS DistinctClaims
    GROUP BY DistinctClaims.InstanceId
),
TotalClaimAmounts AS (
    SELECT
        T3.InstanceId,
        ISNULL(SUM(T3.TravelExpenseClaimAmount), 0.0) AS TravelExpenseClaimAmount
    FROM
        Proof..[Client_BAF56E7B-1074-4997-AEB7-C985D6AC770E] T3
    GROUP BY T3.InstanceId
),
LeaveAggregates AS (
    SELECT
        T2.InstanceId,
        (SELECT
             T2a.LeaveTypeId_Value AS leavetype,
             T2a.DurationOfVisit AS duration,
             T2a.CurrentVisitFrom as startdate,
             T2a.CurrentVisitTo as enddate
          FROM Proof..[Client_925B9F28-1AA3-41EC-81D7-5EA5E3F66B19] T2a
          WHERE T2a.InstanceId = T2.InstanceId
            AND T2a.LeaveTypeId_Value IS NOT NULL
            AND T2a.LeaveTypeId_Value <> ''
            AND T2a.CurrentVisitFrom IS NOT NULL
            AND T2a.CurrentVisitTo IS NOT NULL
            AND ISNULL(T2a.DurationOfVisit, 0) > 0
          FOR JSON PATH) AS LeaveObjectsJson
    FROM Proof..[Client_925B9F28-1AA3-41EC-81D7-5EA5E3F66B19] T2
    GROUP BY T2.InstanceId
)
SELECT
    T1.EmployeeID,
    LEFT(T1.FacultyName, NULLIF(CHARINDEX(' / ', T1.FacultyName), 0) - 1) AS FacultyName,
    T1.Department,
    T1.Designation,
    LEFT(T1.FacultyName, NULLIF(CHARINDEX(' / ', T1.FacultyName), 0) - 1)
        + ' / ' + T1.Department + ' / ' + T1.Designation AS FacultyDetails,
    T1.LProposedVisitFrom,
    T1.LProposedVisitTo,
    T1.Country,
    T1.CityTown,
    T1.CoverPageNo,
    T1.NatureOfParticipation,
    T1.NatureOfParticipation_Value,
    T1.InitiatedOn,
    CA.ClaimTypes,
    (
        SELECT CASE 
            WHEN LA.LeaveObjectsJson IS NULL OR LA.LeaveObjectsJson = '[]'
                THEN NULL
            ELSE JSON_QUERY(
                LEFT(LA.LeaveObjectsJson, LEN(LA.LeaveObjectsJson) - 1)
                + ', {"TravelExpenseClaimAmount":' + CAST(TCA.TravelExpenseClaimAmount AS NVARCHAR(20)) + '}]'
            )
        END
    ) AS leaveDetails
FROM
    Proof..[Client_2F62B79B-D6C1-404A-88B4-627F58F1EE07] T1
LEFT JOIN ClaimAggregates CA
    ON CA.InstanceId = T1.InstanceId
LEFT JOIN TotalClaimAmounts TCA
    ON TCA.InstanceId = T1.InstanceId
LEFT JOIN LeaveAggregates LA
    ON LA.InstanceId = T1.InstanceId
WHERE
     T1.InitiatedOn >= '2025-09-01'
     AND T1.InitiatedOn <= GETDATE()
     AND T1.EmployeeID IS NOT NULL
     AND T1.FacultyName IS NOT NULL
     AND T1.Department IS NOT NULL
     AND T1.Designation IS NOT NULL
     AND T1.LProposedVisitFrom IS NOT NULL
     AND T1.LProposedVisitTo IS NOT NULL
     AND T1.Country IS NOT NULL
     AND T1.CityTown IS NOT NULL
     AND T1.CoverPageNo IS NOT NULL
     AND T1.NatureOfParticipation IS NOT NULL
     AND T1.NatureOfParticipation_Value IS NOT NULL
     AND CA.ClaimTypes IS NOT NULL
     AND TCA.TravelExpenseClaimAmount IS NOT NULL
     AND LA.LeaveObjectsJson IS NOT NULL
     AND LA.LeaveObjectsJson <> '[]'
ORDER BY
    EmployeeID;
  `

	rows, err := mssqlDB.QueryContext(context.Background(), query)
	if err != nil {
		log.Println("Query error:", err)
		return
	}
	defer rows.Close()

	count := 0

	for rows.Next() {
		var o OfficeOrder
		err := rows.Scan(
			&o.EmployeeID,
			&o.FacultyName,
			&o.Department,
			&o.Designation,
			&o.FacultyDetails,
			&o.VisitFrom,
			&o.VisitTo,
			&o.Country,
			&o.CityTown,
			&o.CoverPageNo,
			&o.NatureOfParticipation,
			&o.NatureOfParticipationValue,
			&o.InitiatedOn,
			&o.ClaimType,
			&o.LeaveDetails,
		)
		if err != nil {
			log.Println("Row scan error:", err)
			continue
		}

		var exists bool
		err = pgDB.QueryRow(`SELECT EXISTS (
            SELECT 1 FROM WF_officeorder WHERE CoverPageNo = $1 )`, o.CoverPageNo).Scan(&exists)
		if err != nil {
			log.Println("Existence check error:", err)
			continue
		}

		if exists {
			continue
		}

		_, err = pgDB.Exec(`
            INSERT INTO WF_officeorder (
                EmployeeID, FacultyName, Department, Designation, FacultyDetails,
                VisitFrom, VisitTo, Country, CityTown, CoverPageNo,
                NatureOfParticipation, NatureOfParticipation_Value,
                InitiatedOn, ClaimType, LeaveDetails,
                Status, UpdatedBy
            ) VALUES (
                $1,$2,$3,$4,$5,
                $6,$7,$8,$9,$10,
                $11,$12,
                $13,$14,$15,
                0, 'WF_Admin'
            )
`,
			nullStringToInterface(o.EmployeeID),
			nullStringToInterface(o.FacultyName),
			nullStringToInterface(o.Department),
			nullStringToInterface(o.Designation),
			nullStringToInterface(o.FacultyDetails),
			nullTimeToInterface(o.VisitFrom),
			nullTimeToInterface(o.VisitTo),
			nullStringToInterface(o.Country),
			nullStringToInterface(o.CityTown),
			nullStringToInterface(o.CoverPageNo),
			nullStringToInterface(o.NatureOfParticipation),
			nullStringToInterface(o.NatureOfParticipationValue),
			nullTimeToInterface(o.InitiatedOn),
			nullStringToInterface(o.ClaimType),
			nullStringToInterface(o.LeaveDetails),
		)

		if err != nil {
			log.Printf("Insert error for EmployeeID %s, CoverPageNo %s: %v\n", o.EmployeeID, o.CoverPageNo, err)
			continue
		}

		count++
	}

	log.Printf("✅ Sync complete. %d new records inserted.\n", count)
}

// HTTP handler for manual sync
func HandleSync(w http.ResponseWriter, r *http.Request) {
	go SyncData()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "Sync triggered"})
}
