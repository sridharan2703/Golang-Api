// Package controllerscommon contains APIs for updating workflow master records.
//
// This API updates master records with badge, priority, and starred status based on task ID.
//
// Path: Workflow Management
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 26-08-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 27-10-2025
package controllerscommon

import (
	credentials "Hrmodule/dbconfig"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB
var EncryptionKeyforsignature string

// -------------------- INIT FUNCTION --------------------
func init() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ Warning: .env file not found, using system environment variables")
	}

	// Fetch encryption key from .env
	EncryptionKeyforsignature = os.Getenv("EncryptionKeyforsignature")
	if EncryptionKeyforsignature == "" {
		log.Println("❌ ERROR: EncryptionKeyforsignature missing in .env")
	}

	// Open DB connection
	connectionString := credentials.Getdatabasemeivan()

	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("❌ DB open error: %v", err)
	}

	// Test DB connection
	if err = db.Ping(); err != nil {
		log.Fatalf("❌ DB connection failed: %v", err)
	}

	log.Println("✅ DB connected, signature key loaded")
}

func DownloadSignatureHandler(w http.ResponseWriter, r *http.Request) {
	employeeid := r.URL.Query().Get("role")
	if employeeid == "" {
		http.Error(w, "employeeid is required", http.StatusBadRequest)
		return
	}

	var decrypted []byte
	query := `
		SELECT public.pgp_sym_decrypt_bytea(signature, $1)
		FROM meivan.employee_signatures_new
		WHERE role=$2
		ORDER BY createdon DESC
		LIMIT 1
	`

	err := db.QueryRow(query, EncryptionKeyforsignature, employeeid).Scan(&decrypted)
	if err != nil {
		http.Error(w, "Decryption failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=signature.png")
	w.Header().Set("Content-Type", "image/png")
	w.Write(decrypted)
}

