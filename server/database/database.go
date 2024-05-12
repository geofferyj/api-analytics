package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)

var queries string = `
-- Ensure necessary extension for UUID generation is available
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Creating the 'users' table
CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    api_key UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    last_accessed TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

-- Creating the 'user_agents' table
CREATE TABLE IF NOT EXISTS user_agents (
    id SERIAL PRIMARY KEY,
    user_agent VARCHAR(255) UNIQUE NOT NULL
);

-- Creating the 'requests' table
CREATE TABLE IF NOT EXISTS requests (
    request_id SERIAL PRIMARY KEY,
    api_key UUID NOT NULL,
    ip_address CIDR,
    path VARCHAR(255),
    hostname VARCHAR(255),
    user_agent_id INT,
    method INT,
    response_time INT,
    status INT,
    location VARCHAR(255),
    user_id UUID,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    FOREIGN KEY (api_key) REFERENCES users(api_key),
    FOREIGN KEY (user_agent_id) REFERENCES user_agents(id)
);

-- Creating the 'monitor' table
CREATE TABLE IF NOT EXISTS monitor (
    monitor_id SERIAL PRIMARY KEY,
    api_key UUID NOT NULL,
    url VARCHAR(255) NOT NULL,
    secure BOOLEAN NOT NULL,
    ping BOOLEAN NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (api_key) REFERENCES users(api_key)
);

-- Creating the 'pings' table
CREATE TABLE IF NOT EXISTS pings (
    ping_id SERIAL PRIMARY KEY,
    api_key UUID NOT NULL,
    url VARCHAR(255) NOT NULL,
    response_time INT,
    status INT,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    FOREIGN KEY (api_key) REFERENCES users(api_key)
);

-- Indices for the 'users' table
CREATE INDEX IF NOT EXISTS idx_users_api_key ON users (api_key);

-- Indices for the 'requests' table
CREATE INDEX IF NOT EXISTS idx_requests_api_key ON requests (api_key);
CREATE INDEX IF NOT EXISTS idx_requests_created_at ON requests (created_at);
CREATE INDEX IF NOT EXISTS idx_requests_user_agent_id ON requests (user_agent_id);

-- Indices for the 'monitor' table
CREATE INDEX IF NOT EXISTS idx_monitor_api_key ON monitor (api_key);
CREATE UNIQUE INDEX IF NOT EXISTS idx_monitor_api_key_url ON monitor (api_key, url);

-- Indices for the 'pings' table
CREATE INDEX IF NOT EXISTS idx_pings_api_key ON pings (api_key);
CREATE INDEX IF NOT EXISTS idx_pings_api_key_url ON pings (api_key, url);


`

var dbUsername string
var dbPassword string
var dbName string

type UserRow struct {
	UserID     string    `json:"user_id"`
	APIKey     string    `json:"api_key"`
	CreatedAt  time.Time `json:"created_at"`
	LastAccess time.Time `json:"last_accessed"`
}

type RequestRow struct {
	RequestID    int            `json:"request_id"`
	APIKey       string         `json:"api_key"`
	Path         string         `json:"path"`
	Hostname     sql.NullString `json:"hostname"`
	IPAddress    sql.NullString `json:"ip_address"`
	Location     sql.NullString `json:"location"`
	UserAgent    sql.NullString `json:"user_agent"`
	Method       int16          `json:"method"`
	Status       int16          `json:"status"`
	ResponseTime int16          `json:"response_time"`
	Framework    int16          `json:"framework"`
	CreatedAt    time.Time      `json:"created_at"`
}

type MonitorRow struct {
	APIKey    string    `json:"api_key"`
	URL       string    `json:"url"`
	Secure    bool      `json:"secure"`
	Ping      bool      `json:"ping"`
	CreatedAt time.Time `json:"created_at"`
}

type PingsRow struct {
	APIKey       string    `json:"api_key"`
	URL          string    `json:"url"`
	ResponseTime int       `json:"response_time"`
	Status       int       `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

func getDBLogin() (string, string) {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	panic(err)
	// }

	username := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")
	return username, password
}

func openDBConnection(database string, username string, password string) *sql.DB {
	args := fmt.Sprintf("host=%s port=%d dbname=%s user='%s' password=%s sslmode=%s", "localhost", 5432, database, username, password, "disable")

	db, err := sql.Open("postgres", args)
	if err != nil {
		panic(err)
	}

	return db
}

func OpenDBConnection() *sql.DB {
	if dbUsername == "" || dbPassword == "" || dbName == "" {
		username, password := getDBLogin()
		database := os.Getenv("POSTGRES_DATABASE")
		dbUsername = username
		dbPassword = password
		dbName = database
		return openDBConnection(database, username, password)
	} else {
		return openDBConnection(dbName, dbUsername, dbPassword)
	}
}

func OpenDBConnectionNamed(database string) *sql.DB {
	if dbUsername == "" || dbPassword == "" {
		username, password := getDBLogin()
		return openDBConnection(database, username, password)
	} else {
		return openDBConnection(database, dbUsername, dbPassword)
	}
}

func CreateUsersTable(db *pgx.Conn) error {

	_, err := db.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS users (user_id UUID NOT NULL, api_key UUID NOT NULL, created_at TIMESTAMPTZ NOT NULL, last_accessed TIMESTAMPTZ, PRIMARY KEY (api_key));")
	return err
}

func CreateRequestsTable(db *pgx.Conn) error {

	_, err := db.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS requests (request_id INTEGER, api_key UUID NOT NULL, path VARCHAR(255) NOT NULL, hostname VARCHAR(255), ip_address CIDR, location CHAR(2), user_agent VARCHAR(255), method SMALLINT NOT NULL, status SMALLINT NOT NULL, response_time SMALLINT NOT NULL, framework SMALLINT NOT NULL, created_at TIMESTAMPTZ NOT NULL, PRIMARY KEY (api_key));")
	return err
}

func CreateMonitorTable(db *pgx.Conn) error {

	_, err := db.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS monitor (api_key UUID NOT NULL, url VARCHAR(255) NOT NULL, secure BOOLEAN, PING BOOLEAN, created_at TIMESTAMPTZ NOT NULL, PRIMARY KEY (api_key));")
	return err
}

func CreatePingsTable(db *pgx.Conn) error {
	_, err := db.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS pings (api_key UUID NOT NULL, url VARCHAR(255) NOT NULL, response_time INTEGER, status SMALLINT, created_at TIMESTAMPTZ NOT NULL, PRIMARY KEY (api_key));")
	return err
}

func CreateTables() error {
	db := NewConnection()
	defer db.Close(context.Background())
	_, err := db.Exec(context.Background(), queries)
	return err
}

func DeleteUser(apiKey string) error {
	db := OpenDBConnection()
	defer db.Close()

	query := fmt.Sprintf("DELETE FROM users WHERE api_key = '%s';", apiKey)
	_, err := db.Query(query)
	return err
}

func DeleteRequests(apiKey string) error {
	db := OpenDBConnection()
	defer db.Close()

	query := fmt.Sprintf("DELETE FROM requests WHERE api_key = '%s';", apiKey)
	_, err := db.Query(query)
	return err
}

func DeleteMonitors(apiKey string) error {
	db := OpenDBConnection()
	defer db.Close()

	query := fmt.Sprintf("DELETE FROM monitor WHERE api_key = '%s';", apiKey)
	_, err := db.Query(query)
	return err
}

func DeletePings(apiKey string) error {
	db := OpenDBConnection()
	defer db.Close()

	query := fmt.Sprintf("DELETE FROM pings WHERE api_key = '%s';", apiKey)
	_, err := db.Query(query)
	return err
}

func InsertUserData(db *sql.DB, rows []UserRow) error {
	if len(rows) == 0 {
		return nil
	}

	var query strings.Builder
	query.WriteString("INSERT INTO users (api_key, user_id, created_at) VALUES")
	for i, user := range rows {
		if i > 0 {
			query.WriteString(",")
		}
		query.WriteString(fmt.Sprintf(" ('%s', '%s', '%s')", user.APIKey, user.UserID, user.CreatedAt.UTC().Format(time.RFC3339)))
	}
	query.WriteString(";")

	_, err := db.Query(query.String())
	return err
}

func InsertRequestsData(db *sql.DB, rows []RequestRow) error {
	blockSize := 1000
	if len(rows) == 0 {
		return nil
	}

	var query strings.Builder
	for i, request := range rows {
		if i%blockSize == 0 {
			query = strings.Builder{}
			query.WriteString("INSERT INTO requests (api_key, method, created_at, path, status, response_time, framework, user_agent, hostname, ip_address, location) VALUES")
		}
		if i%blockSize > 0 {
			query.WriteString(",")
		}

		fmtIPAddress := nullInsertString(request.IPAddress)
		fmtHostname := nullInsertString(request.Hostname)
		fmtLocation := nullInsertString(request.Location)
		fmtUserAgent := nullInsertString(request.UserAgent)

		query.WriteString(fmt.Sprintf(" ('%s', %d, '%s', '%s', %d, %d, %d, %s, %s, %s)", request.APIKey, request.Method, request.CreatedAt.UTC().Format(time.RFC3339), request.Path, request.Status, request.ResponseTime, request.Framework, fmtUserAgent, fmtHostname, fmtIPAddress, fmtLocation))

		if (i+1)%blockSize == 0 || i == len(rows)-1 {
			query.WriteString(";")

			fmt.Println("Writing %d rows to requests database...", blockSize)

			db = OpenDBConnection()
			_, err := db.Query(query.String())
			db.Close()
			if err != nil {
				return err
			}
			time.Sleep(8 * time.Second) // Pause to avoid exceeding memory space
		}
	}
	return nil
}

func nullInsertString(value sql.NullString) string {
	if value.String == "" {
		return "NULL"
	} else {
		return fmt.Sprintf("'%s'", value.String)
	}
}

func InsertMonitorData(db *sql.DB, rows []MonitorRow) error {
	if len(rows) == 0 {
		return nil
	}

	var query strings.Builder
	query.WriteString("INSERT INTO monitor (api_key, url, secure, ping, created_at) VALUES")
	for i, monitor := range rows {
		if i > 0 {
			query.WriteString(",")
		}
		query.WriteString(fmt.Sprintf(" ('%s', '%s', %t, %t, '%s')", monitor.APIKey, monitor.URL, monitor.Secure, monitor.Ping, monitor.CreatedAt.UTC().Format(time.RFC3339)))
	}
	query.WriteString(";")

	_, err := db.Query(query.String())
	return err
}

func InsertPingsData(db *sql.DB, rows []PingsRow) error {
	if len(rows) == 0 {
		return nil
	}

	var query strings.Builder
	query.WriteString("INSERT INTO pings (api_key, url, response_time, status, created_at) VALUES")
	for i, monitor := range rows {
		if i > 0 {
			query.WriteString(",")
		}
		query.WriteString(fmt.Sprintf(" ('%s', '%s', %d, %d, '%s')", monitor.APIKey, monitor.URL, monitor.ResponseTime, monitor.Status, monitor.CreatedAt.UTC().Format(time.RFC3339)))
	}
	query.WriteString(";")

	_, err := db.Query(query.String())
	return err
}
