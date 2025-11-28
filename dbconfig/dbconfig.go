// Package credentials contains data structures and database access logic.
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On:30-07-2025
//
// Last Modified By: Sivabala
//
// Last Modified Date: 30-07-2025
//
// Path:Login Page
package credentials

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/denisenkom/go-mssqldb" // Add MSSQL driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Failed to load .env: " + err.Error())
	}
}

// getDBConnectionString constructs and verifies a database connection string
// for the given driver (e.g., "postgres", "mysql", or "mssql"). It opens and pings the
// database to ensure the connection is valid. It returns the connection string
// or panics on failure.
func getDBConnectionString(driver, server, user, password, database, port string) string {
	var connStr string

	switch driver {
	case "postgres":
		// ✅ Correct DSN format for lib/pq and gorm postgres driver
		connStr = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			server, user, password, database, port)
	case "mysql":
		connStr = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, server, port, database)
	case "mssql":
		// MSSQL connection string format
		connStr = fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;port=%s;encrypt=disable",
			server, user, password, database, port)
	case "sqlserver":
		// Alternative format for SQL Server
		connStr = fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s&encrypt=disable",
			user, password, server, port, database)
	default:
		panic("Unsupported DB driver: " + driver)
	}

	db, err := sql.Open(driver, connStr)
	if err != nil {
		panic("Failed to open DB connection: " + err.Error())
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic("Database connection failed: " + err.Error())
	}

	return connStr
}

// getPostgresConnectionString constructs and verifies a Postgres connection string
func getPostgresConnectionStringforpostgre(server, user, password, database, port, schema string) string {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable search_path=%s",
		server, user, password, database, port, schema)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic("Failed to open Postgres connection: " + err.Error())
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic("Postgres connection failed: " + err.Error())
	}

	return connStr
}

// logFullyMaskedConnection masks all sensitive information in connection string
func logFullyMaskedConnection(connStr, password, user, host, database, dbType string) {
	safeConnStr := connStr

	// Replace sensitive information with ****
	safeConnStr = strings.Replace(safeConnStr, password, "****", -1)
	safeConnStr = strings.Replace(safeConnStr, user, "****", -1)
	safeConnStr = strings.Replace(safeConnStr, host, "****", -1)
	safeConnStr = strings.Replace(safeConnStr, database, "****", -1)

}

// GetGnanaThalamSchemaConnection returns a Postgres connection string for a given schema
func GetGnanaThalamSchemaConnection(schema string) string {
	server := os.Getenv("server_gnana")
	user := os.Getenv("user_gnana")
	password := os.Getenv("password_gnana")
	database := os.Getenv("database_gnana")
	port := os.Getenv("port_gnana")

	connStr := getPostgresConnectionStringforpostgre(server, user, password, database, port, schema)
	logFullyMaskedConnection(connStr, password, user, server, database, "Postgres")

	return connStr
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Postgres database
// Getdatabasehr returns Postgres connection string
func Getdatabasehr() string {
	return GetGnanaThalamSchemaConnection("humanresources")
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Getdatabasemeivan returns Postgres connection string
func Getdatabasemeivan() string {
	return GetGnanaThalamSchemaConnection("Meivan")
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// GetdatabaseWF_officeorder returns Postgres connection string

func GetdatabaseWF_officeorder() string {
	return GetGnanaThalamSchemaConnection("WF_Integration")
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 17 phpmyadmin
// GetMySQLDatabase17 returns MySQL connection string
func GetMySQLDatabase17() string {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	connStr := getDBConnectionString("mysql", host, user, password, database, port)
	logFullyMaskedConnection(connStr, password, user, host, database, "MySQL")

	return connStr
}

// 17 phpmyadmin
// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// GetMySQLDatabase17HR returns MySQL HR connection string
func GetMySQLDatabase17HR() string {
	host := os.Getenv("DB_HOST_HR")
	user := os.Getenv("DB_USER_HR")
	password := os.Getenv("DB_PASSWORD_HR")
	database := os.Getenv("DB_NAME_HR")
	port := os.Getenv("DB_PORT_HR")

	connStr := getDBConnectionString("mysql", host, user, password, database, port)
	logFullyMaskedConnection(connStr, password, user, host, database, "MySQL HR")

	return connStr
}

// MSSQL
// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// GetTestdatabase15 returns a verified MSSQL connection string for
// the test credentials with IITMAcademics database from environment variables
func GetTestdatabase15() string {
	server := os.Getenv("server1")
	user := os.Getenv("userId1")
	password := os.Getenv("password1")
	database := os.Getenv("database1")
	port := os.Getenv("port1")

	connStr := getDBConnectionString("mssql", server, user, password, database, port)
	logFullyMaskedConnection(connStr, password, user, server, database, "MSSQL")

	return connStr
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// GetLivedatabase10 returns a verified MSSQL connection string for
// the test credentials with IITMAcademics database from environment variables
func GetLivedatabase10() string {
	server := os.Getenv("server2")
	user := os.Getenv("userId2")
	password := os.Getenv("password2")
	database := os.Getenv("database2")
	port := os.Getenv("port2")

	connStr := getDBConnectionString("mssql", server, user, password, database, port)
	logFullyMaskedConnection(connStr, password, user, server, database, "MSSQL")

	return connStr
}
