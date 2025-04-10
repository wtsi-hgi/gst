/*******************************************************************************
 * Copyright (c) 2025 Genome Research Ltd.
 *
 * Author: Sendu Bala <sb10@sanger.ac.uk>
 *
 * Permission is hereby granted, free of charge, to any person obtaining
 * a copy of this software and associated documentation files (the
 * "Software"), to deal in the Software without restriction, including
 * without limitation the rights to use, copy, modify, merge, publish,
 * distribute, sublicense, and/or sell copies of the Software, and to
 * permit persons to whom the Software is furnished to do so, subject to
 * the following conditions:
 *
 * The above copyright notice and this permission notice shall be included
 * in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
 * IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
 * CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
 * TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
 * SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 ******************************************************************************/

package db

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

//go:embed query.sql
var sqlFiles embed.FS

// DBConnector provides an interface to connect to a database.
type DBConnector interface {
	Connect() (*sql.DB, error)
	Close() error
}

// MySQLConnector implements DBConnector for MySQL databases.
type MySQLConnector struct {
	db *sql.DB
}

// Connect establishes a connection to the MySQL database using environment
// vars.
func (c *MySQLConnector) Connect() (*sql.DB, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	user := os.Getenv("GST_AUTOMATION_SQL_USER")
	password := os.Getenv("GST_AUTOMATION_SQL_PASS")
	host := os.Getenv("GST_AUTOMATION_SQL_HOST")
	port := os.Getenv("GST_AUTOMATION_SQL_PORT")
	dbname := os.Getenv("GST_AUTOMATION_SQL_DB")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user, password, host, port, dbname)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	c.db = db
	return db, nil
}

// Close closes the database connection.
func (c *MySQLConnector) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// GetEmbeddedSQL retrieves the SQL query from the embedded file.
func GetEmbeddedSQL() string {
	data, err := sqlFiles.ReadFile("query.sql")
	if err != nil {
		// This should never happen as the file is embedded at compile time
		panic(fmt.Sprintf("failed to read embedded SQL file: %v", err))
	}
	return strings.TrimSpace(string(data))
}

// parseRows converts SQL rows to a TrackedSampleCollection.
func parseRows(rows *sql.Rows) (*TrackedSampleCollection, error) {
	var samples []TrackedSample

	for rows.Next() {
		// Use NullString for fields that might be NULL
		var s TrackedSample
		var runIDNull, platformNull, pipelineNull, qcPassNull sql.NullString

		err := rows.Scan(
			&s.StudyID,
			&s.StudyName,
			&s.FacultySponsor,
			&s.Programme,
			&s.SangerSampleID,
			&s.SupplierName,
			&s.ManifestCreated,
			&s.ManifestUploaded,
			&s.LabwareReceived,
			&s.LabwareHumanBarcode,
			&s.OrderMade,
			&s.LibraryStart,
			&s.LibraryComplete,
			&s.LibraryTime,
			&runIDNull,
			&platformNull,
			&pipelineNull,
			&s.SequencingRunStart,
			&s.SequencingQCComplete,
			&s.SequencingTime,
			&qcPassNull,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		// Convert NullString to string (empty string if NULL)
		s.RunID = getNullableString(runIDNull)
		s.Platform = getNullableString(platformNull)
		s.Pipeline = getNullableString(pipelineNull)
		s.QCPass = getNullableString(qcPassNull)

		samples = append(samples, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return &TrackedSampleCollection{Samples: samples}, nil
}

// getNullableString returns the string value of a sql.NullString,
// or empty string if the value is NULL.
func getNullableString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
