package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	credentials "Hrmodule/dbconfig"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var dbAESKey []byte

// LoadKeyFromDB loads the AES key from the database using the given sessionID.
func LoadKeyFromDB(sessionID string) error {
	connectionString := credentials.Getdatabasemeivan()

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return fmt.Errorf("DB open error: %v", err)
	}
	defer db.Close()

	query := `
		SELECT decryptkey
		FROM session_data
		WHERE is_active = 1 AND session_id = $1
	`

	var key string
	err = db.QueryRow(query, sessionID).Scan(&key)
	if err != nil {
		return fmt.Errorf("failed to fetch decryptkey: %w", err)
	}
	// Print the fetched key (for debugging, remove in production)
	fmt.Printf("Fetched AES Key for session %s: %s\n", sessionID, key)

	if len(key) != 32 {
		return fmt.Errorf("decryptkey must be 32 bytes, got %d", len(key))
	}

	dbAESKey = []byte(key)
	return nil
}

// Encryptnew encrypts data using the loaded AES key (dbAESKey) with AES-GCM and returns base64 string.
func Encryptnew(plainText []byte) (string, error) {
	if dbAESKey == nil {
		return "", errors.New("AES key is not loaded")
	}

	block, err := aes.NewCipher(dbAESKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %v", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %v", err)
	}

	cipherText := aesGCM.Seal(nonce, nonce, plainText, nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}
