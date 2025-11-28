// Package controllersofficeorder handles HTTP APIs for history of officeorder pdf API.
//
// --- Creator's Info ---
// Creator: Sridharan
//
// Created On: 29-10-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 29-10-2025
package controllersofficeorder

import (
	"Hrmodule/auth"
	credentials "Hrmodule/dbconfig"
	"Hrmodule/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type OrderPDFRequest struct {
	Data string `json:"Data"`
}

type PDFItem struct {
	OrderNo  string
	FilePath string // This will be the path to the downloaded temp file
}

// ConnectSFTP establishes an SFTP connection using credentials from .env or Fallback
func ConnectSFTP() (*sftp.Client, *ssh.Client, error) {
	// 1. Try to get from .env
	host := os.Getenv("SFTP_HOST")
	port := os.Getenv("SFTP_PORT")
	user := os.Getenv("SFTP_USERNAME")
	pass := os.Getenv("SFTP_PASSWORD")

	// 2. Fallback if .env failed to load (Fixes "dial tcp :0" error)
	if host == "" {
		host = "wfvault.iitm.ac.in"
	}
	if port == "" {
		port = "22"
	}
	if user == "" {
		user = "NAS_Admin"
	}
	if pass == "" {
		pass = "@dm1n#N@&(102424)"
	}

	address := fmt.Sprintf("%s:%s", host, port)

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	// SSH connection
	sshConn, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return nil, nil, fmt.Errorf("SSH dial failed to %s: %v", address, err)
	}

	// Create SFTP
	client, err := sftp.NewClient(sshConn)
	if err != nil {
		sshConn.Close()
		return nil, nil, fmt.Errorf("SFTP client creation failed: %v", err)
	}

	return client, sshConn, nil
}

func FetchOrderHistoryPDF(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Use POST", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req OrderPDFRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	parts := strings.Split(req.Data, "||")
	if len(parts) != 2 {
		http.Error(w, "Invalid Data format", http.StatusBadRequest)
		return
	}

	pid := parts[0]
	encrypted := parts[1]

	key, err := utils.GetDecryptKey(pid)
	if err != nil {
		http.Error(w, "Failed to get decrypt key", http.StatusUnauthorized)
		return
	}

	decryptedJSON, err := utils.DecryptAES(encrypted, key)
	if err != nil {
		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}

	var decrypted map[string]interface{}
	if err := json.Unmarshal([]byte(decryptedJSON), &decrypted); err != nil {
		http.Error(w, "Invalid decrypted JSON", http.StatusBadRequest)
		return
	}

	token, ok := decrypted["token"].(string)
	if !ok || token == "" {
		http.Error(w, "Missing or invalid token", http.StatusBadRequest)
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

		orderNo, ok := decrypted["officeorderhistorypdf"].(string)
		if !ok || orderNo == "" {
			http.Error(w, "Missing 'officeorderhistorypdf'", http.StatusBadRequest)
			return
		}

		// Extract templateType
		var templateType string
		found := false
		for k, v := range decrypted {
			if strings.EqualFold(k, "templateType") {
				if valStr, ok := v.(string); ok && strings.TrimSpace(valStr) != "" {
					templateType = strings.ToLower(strings.TrimSpace(valStr))
					found = true
				}
				break
			}
		}

		if !found {
			http.Error(w, "'templateType' missing in decrypted payload", http.StatusBadRequest)
			return
		}

		// 1. Connect to DB to get related Order Numbers
		connectionString := credentials.Getdatabasemeivan()
		db, err := sql.Open("postgres", connectionString)
		if err != nil {
			http.Error(w, "DB connection failed", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		// Fetch all related orders (Original + Amendments + Cancellations)
		query := `
		WITH base AS (
			SELECT regexp_replace($1, '/A[0-9]+(/CAN)?$', '') AS base_order
		)
		SELECT m.order_no
		FROM meivan.pcr_m m, base b
		WHERE m.order_no LIKE b.base_order || '%'
		AND m.task_status_id = 3
		ORDER BY length(m.order_no), m.order_no`

		rows, err := db.Query(query, orderNo)
		if err != nil {
			http.Error(w, "DB query failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var allOrders []string
		for rows.Next() {
			var o string
			if err := rows.Scan(&o); err != nil {
				continue
			}
			allOrders = append(allOrders, o)
		}

		if len(allOrders) == 0 {
			http.Error(w, "No related order numbers found", http.StatusNotFound)
			return
		}

		// 2. Identify Target Filenames
		// We create a map of "filename" -> "orderNo" to quickly check during the walk
		targetFiles := make(map[string]string)
		prefix := "officecopy_"
		if templateType == "usercopy" {
			prefix = "usercopy_"
		}

		for _, ord := range allOrders {
			// Sanitize: Replace "/" with "_" (e.g., No.F.Admn.I/PCR/2025/000001 -> No.F.Admn.I_PCR_2025_000001)
			sanitized := strings.ReplaceAll(ord, "/", "_")
			fileName := fmt.Sprintf("%s%s.pdf", prefix, sanitized)
			targetFiles[fileName] = ord
		}

		// 3. Connect to SFTP
		sftpClient, sshClient, err := ConnectSFTP()
		if err != nil {
			http.Error(w, "SFTP Connection Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer sshClient.Close()
		defer sftpClient.Close()

		// Create Temp Dir
		tempDir, err := os.MkdirTemp("", "sftp_search_")
		if err != nil {
			http.Error(w, "Failed to create temp dir", http.StatusInternalServerError)
			return
		}
		defer os.RemoveAll(tempDir)

		var results []PDFItem
		baseSearchPath := "/HR_Test/meivan/off_order"

		fmt.Printf("\n=== STARTING SFTP SEARCH (%s) ===\n", strings.ToUpper(templateType))
		fmt.Printf("Searching in: %s\n", baseSearchPath)
		fmt.Printf("Looking for %d files...\n", len(targetFiles))

		// 4. Walk the SFTP Directory
		walker := sftpClient.Walk(baseSearchPath)

		for walker.Step() {
			if walker.Err() != nil {
				log.Printf("SFTP Walk Error: %v", walker.Err())
				continue
			}

			remotePath := walker.Path()
			fileInfo := walker.Stat()

			if fileInfo.IsDir() {
				continue
			}

			fileName := fileInfo.Name()

			// Check if this file is in our target list
			if orderNo, exists := targetFiles[fileName]; exists {
				fmt.Printf("✓ FOUND: %s at %s\n", fileName, remotePath)

				// Download File
				remoteFile, err := sftpClient.Open(remotePath)
				if err != nil {
					log.Printf("Failed to open remote file %s: %v", remotePath, err)
					continue
				}

				localPath := filepath.Join(tempDir, fileName)
				localFile, err := os.Create(localPath)
				if err != nil {
					remoteFile.Close()
					log.Printf("Failed to create local file: %v", err)
					continue
				}

				_, err = io.Copy(localFile, remoteFile)
				localFile.Close()
				remoteFile.Close()

				if err != nil {
					log.Printf("Failed to download file: %v", err)
					continue
				}

				// Add to results
				results = append(results, PDFItem{
					OrderNo:  orderNo,
					FilePath: localPath,
				})

				// Remove from targets to optimize
				delete(targetFiles, fileName)

				// Stop walking if we found everything
				if len(targetFiles) == 0 {
					fmt.Println("All files found. Stopping search.")
					break
				}
			}
		}

		if len(results) == 0 {
			http.Error(w, "No files found on SFTP server", http.StatusNotFound)
			return
		}

		// 5. Sort Results (Amendments Descending)
		parseOrder := func(order string) (amendNum int, isCancelled bool) {
			isCancelled = strings.HasSuffix(order, "/CAN")
			base := order
			if isCancelled {
				base = strings.TrimSuffix(order, "/CAN")
			}
			re := regexp.MustCompile(`/A([0-9]+)$`)
			match := re.FindStringSubmatch(base)
			if len(match) == 2 {
				fmt.Sscanf(match[1], "%d", &amendNum)
			} else {
				amendNum = 0
			}
			return
		}

		sort.Slice(results, func(i, j int) bool {
			amendI, cancelI := parseOrder(results[i].OrderNo)
			amendJ, cancelJ := parseOrder(results[j].OrderNo)

			if amendI != amendJ {
				return amendI > amendJ
			}
			if cancelI != cancelJ {
				return cancelI
			}
			return false
		})

		// 6. Merge PDFs
		var inputFiles []string
		for _, item := range results {
			inputFiles = append(inputFiles, item.FilePath)
		}

		baseOrderNo := regexp.MustCompile(`/A[0-9]+(/CAN)?$`).ReplaceAllString(allOrders[0], "")
		cleanOrderNo := strings.ReplaceAll(baseOrderNo, "/", "_")
		mergedFileName := cleanOrderNo + "_History_" + templateType + ".pdf"
		mergedPath := filepath.Join(tempDir, mergedFileName)

		if err := api.MergeCreateFile(inputFiles, mergedPath, false, nil); err != nil {
			http.Error(w, "Failed merging PDF: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 7. Send Response
		mergedData, err := os.ReadFile(mergedPath)
		if err != nil {
			http.Error(w, "Failed reading merged PDF: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", mergedFileName))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(mergedData)))
		w.Write(mergedData)

		fmt.Println("✓ PDF sent to client successfully")

	})).ServeHTTP(w, r)
}
