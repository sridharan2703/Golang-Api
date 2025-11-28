// package databaseofficeorder handles DB access for Approval page
//
// --- Creator's Info ---
// Creator: Ramya
//
// Created On: 09-10-2025
// Last Modified By:
// Last Modified Date:
package databaseofficeorder

import (
	// Assuming the package path structure based on your sample code
	credentials "Hrmodule/dbconfig"
	modelsofficeorder "Hrmodule/models/officeorder"
	"database/sql"
	"fmt"
	"net/http"
	"strconv" // Need this for string-to-int conversion
)

func GetOfficeOrderMasterFromDB(w http.ResponseWriter, r *http.Request, taskStatusID string, employeeID string, coverPageNo string) ([]modelsofficeorder.OfficeOrderMasterStructure, int, error) {

	statusID, err := strconv.Atoi(taskStatusID)
	if err != nil {
		errMsg := fmt.Sprintf("Invalid taskStatusID provided: %s. Must be an integer.", taskStatusID)
		http.Error(w, errMsg, http.StatusBadRequest)
		return nil, 0, fmt.Errorf(errMsg)
	}
	connectionString := credentials.Getdatabasemeivan()

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		http.Error(w, fmt.Sprintf("DB open error for office order master: %v", err), http.StatusInternalServerError)
		return nil, 0, err
	}
	defer db.Close()

	data, err := modelsofficeorder.GetOfficeOrderMasters(db, statusID, employeeID, coverPageNo)
	if err != nil {

		return nil, 0, fmt.Errorf("retrieving office order master failed: %v", err)
	}

	return data, len(data), nil
}
