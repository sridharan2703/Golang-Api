// Package controllerslogin provides generate and send otp to the mobile number
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On: 26-08-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 26-08-2025
package controllerslogin

import (
	"Hrmodule/auth"
	"Hrmodule/utils" // utils.Encrypt function
	"bytes"
	"crypto/aes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ---------------------
// AES Encryption/Decryption
// ---------------------

func pkcs5Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}

func pkcs5Unpad(data []byte) []byte {
	pad := int(data[len(data)-1])
	if pad > len(data) {
		return data
	}
	return data[:len(data)-pad]
}

func EncryptAESHex(plainText, key string) (string, error) {
	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %v", err)
	}
	data := pkcs5Pad([]byte(plainText), aes.BlockSize)
	encrypted := make([]byte, len(data))
	for i := 0; i < len(data); i += aes.BlockSize {
		block.Encrypt(encrypted[i:i+aes.BlockSize], data[i:i+aes.BlockSize])
	}
	return hex.EncodeToString(encrypted), nil
}

func DecryptAESHex(encryptedHex, key string) (string, error) {
	keyBytes := []byte(key)
	encryptedBytes, err := hex.DecodeString(encryptedHex)
	if err != nil {
		return "", fmt.Errorf("invalid hex encoding: %v", err)
	}
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %v", err)
	}
	if len(encryptedBytes)%aes.BlockSize != 0 {
		return "", fmt.Errorf("encrypted data is not a multiple of the block size")
	}
	decrypted := make([]byte, len(encryptedBytes))
	for i := 0; i < len(encryptedBytes); i += aes.BlockSize {
		block.Decrypt(decrypted[i:i+aes.BlockSize], encryptedBytes[i:i+aes.BlockSize])
	}
	decrypted = pkcs5Unpad(decrypted)
	return string(decrypted), nil
}

// ---------------------
// Struct for request
// ---------------------
type SendOTPRequest struct {
	Data string `json:"Data"`
}

// ---------------------
// Generate 6-digit OTP
// ---------------------
func generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(900000) + 100000 // 6-digit OTP
	return fmt.Sprintf("%06d", otp)
}

// ---------------------
// Handler: decrypt & send OTP
// ---------------------
func SendOTPHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}
	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	// Extract token
	var req SendOTPRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// 3️⃣ Split and decrypt
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

	token, ok := decryptedData["token"].(string)
	if !ok || token == "" {
		http.Error(w, "Token not found", http.StatusBadRequest)
		return
	}
	r.Header.Set("token", token)
	// Auth check
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	// Log + handler
	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := auth.IsValidIDFromRequest(r); err != nil {
			http.Error(w, "Invalid TOKEN provided", http.StatusBadRequest)
			return
		}
		// Read request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Parse JSON
		var req SendOTPRequest

		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		// Extract order_type_id from decrypted data
		SendTo, ok := decryptedData["send_to"].(string)
		if !ok || SendTo == "" {
			http.Error(w, "Missing 'send_to' in request data", http.StatusBadRequest)
			return
		}

		// Decrypt mobile number
		encryptionKey := "7xPz!qL3vNc#eRb9Wm@f2Zh8Kd$gYp1B"
		decryptedNumber, err := DecryptAESHex(SendTo, encryptionKey)
		if err != nil {
			http.Error(w, "Failed to decrypt number: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Generate OTP and message
		otp := generateOTP()
		message := fmt.Sprintf("The OTP to verify your mobile number on HMIS-Institute hospital is %s-IITMWF", otp)
		encodedMessage := url.QueryEscape(message)

		// Call Gupshup API with dynamic OTP
		apiURL := fmt.Sprintf(
			"https://enterprise.smsgupshup.com/GatewayAPI/rest?v=1.1&method=SendMessage&msg_type=TEXT&userid=2000230894&auth_scheme=plain&password=z6eucW@q%%20&format=text&msg=%s&send_to=%s",
			encodedMessage, decryptedNumber,
		)

		resp, err := http.Get(apiURL)
		if err != nil {
			http.Error(w, "Failed to call SMS API: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		respBody, _ := ioutil.ReadAll(resp.Body)

		// Prepare response object
		response := map[string]interface{}{
			"decrypted_send_to": decryptedNumber,
			"sms_api_response":  string(respBody),
			"otp_sent":          otp, // optional: include OTP for testing
		}

		// 6️⃣ Marshal & encrypt before sending
		responseJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Response marshal failed", http.StatusInternalServerError)
			return
		}

		encryptedResponse, err := utils.EncryptAES(string(responseJSON), key)
		if err != nil {
			http.Error(w, "Response encryption failed", http.StatusInternalServerError)
			return
		}

		finalResp := map[string]string{
			"Data": fmt.Sprintf("%s||%s", pid, encryptedResponse),
		}

		// ✅ Save exactly what is sent to client
		auth.SaveResponseLog(
			r,
			finalResp,          // only final response
			http.StatusOK,      // status code
			"application/json", // content type
			len(responseJSON),  // size
			string(body),       // original request
		)

		// ✅ Send to client
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(finalResp)
	})).ServeHTTP(w, r)
}
