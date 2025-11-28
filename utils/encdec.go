// Package utils provides utility functions including AES-GCM
// encryption for secure response encoding.
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On:07-07-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 09-07-2025
package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

var db *sql.DB

// --- Connect to Postgres ---
func init() {
	var err error
	connStr := "host=10.24.2.18 port=5432 user=postgres password=NewStrongPassword dbname=GnanaThalam sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("DB ping failed:", err)
	}
	fmt.Println("✅ Connected to PostgreSQL")
}

// --- Fetch decrypt/encrypt key using P_id ---
func GetDecryptKey(pid string) (string, error) {
	var key string
	err := db.QueryRow(
		`SELECT decryptkey FROM meivan.session_data WHERE session_id=$1 AND is_active=1`,
		pid,
	).Scan(&key)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("invalid or inactive id")
		}
		return "", err
	}

	if len(key) != 32 {
		return "", fmt.Errorf("decrypt key must be 32 characters long (got %d)", len(key))
	}

	return key, nil
}

// --- AES-GCM Encrypt/Decrypt ---
func EncryptAES(plaintext, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptAES(ciphertextB64, key string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}
	nonce, encrypted := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// --- Structs ---
type EncryptResponse struct {
	Data string `json:"Data"`
}
type DecryptRequest struct {
	Data string `json:"Data"`
}

// --- /encryptData API ---
func EncryptHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request as generic JSON
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	pidVal, ok := req["P_id"]
	if !ok {
		http.Error(w, "Missing 'P_id' field", http.StatusBadRequest)
		return
	}

	pid := fmt.Sprintf("%v", pidVal)
	key, err := GetDecryptKey(pid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Key fetch failed: %v", err), http.StatusUnauthorized)
		return
	}

	// Remove P_id before encryption
	delete(req, "P_id")

	// Convert the remaining payload to JSON
	payloadBytes, _ := json.Marshal(req)

	// Encrypt payload using DB key
	encrypted, err := EncryptAES(string(payloadBytes), key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Encryption failed: %v", err), http.StatusInternalServerError)
		return
	}

	resp := EncryptResponse{
		Data: fmt.Sprintf("%s||%s", pid, encrypted),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// --- /decryptData API ---
func DecryptHandler(w http.ResponseWriter, r *http.Request) {
	var req DecryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	parts := strings.Split(req.Data, "||")
	if len(parts) != 2 {
		http.Error(w, "Invalid Data format", http.StatusBadRequest)
		return
	}

	pid := parts[0]
	encryptedPart := parts[1]

	key, err := GetDecryptKey(pid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Key fetch failed: %v", err), http.StatusUnauthorized)
		return
	}

	decryptedJSON, err := DecryptAES(encryptedPart, key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Decryption failed: %v", err), http.StatusInternalServerError)
		return
	}

	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(decryptedJSON), &payload); err != nil {
		http.Error(w, "Invalid decrypted JSON", http.StatusBadRequest)
		return
	}

	payload["P_id"] = pid

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}
