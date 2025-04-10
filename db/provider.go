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
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

// QueryProvider defines an interface for executing database queries
// and retrieving results.
type QueryProvider interface {
	Execute() (*TrackedSampleCollection, error)
}

// config holds configuration options for creating a QueryProvider.
type config struct {
	useMock   bool
	mockPath  string
	connector DBConnector
}

// Option is a function that configures a config.
type Option func(*config)

// WithMockData configures the provider to use mock data from a TSV file.
func WithMockData(path string) Option {
	return func(c *config) {
		c.useMock = true
		c.mockPath = path
	}
}

// WithMySQLConnector configures the provider to use a specific MySQL connector.
func WithMySQLConnector(connector DBConnector) Option {
	return func(c *config) {
		c.connector = connector
	}
}

// MySQLQueryProvider implements QueryProvider for MySQL databases.
type MySQLQueryProvider struct {
	connector DBConnector
}

// Execute executes the SQL query and returns the results.
func (p *MySQLQueryProvider) Execute() (*TrackedSampleCollection, error) {
	db, err := p.connector.Connect()
	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}
	defer p.connector.Close()

	// Check for nil db connection - this protects against mock tests
	// that don't configure a proper DB object
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	// Execute the embedded query
	rows, err := db.Query(GetEmbeddedSQL())
	if err != nil {
		return nil, fmt.Errorf("query execution error: %w", err)
	}
	defer rows.Close()

	return parseRows(rows)
}

// MockQueryProvider implements QueryProvider for testing with mock data.
type MockQueryProvider struct {
	tsvPath string
}

// Execute reads sample data from a TSV file instead of the database.
func (p *MockQueryProvider) Execute() (*TrackedSampleCollection, error) {
	file, err := os.Open(p.tsvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open mock data file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read mock data: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("mock data file should contain at least header and one data row")
	}

	return parseMockRecords(records[1:])
}

// parseMockRecords converts TSV records into TrackedSample objects.
func parseMockRecords(records [][]string) (*TrackedSampleCollection, error) {
	samples := make([]TrackedSample, 0, len(records))

	for i, record := range records {
		if len(record) < 21 {
			return nil, fmt.Errorf("row %d has insufficient columns: expected 21, got %d",
				i+1, len(record))
		}

		sample, err := recordToSample(record)
		if err != nil {
			return nil, fmt.Errorf("error parsing row %d: %w", i+1, err)
		}

		samples = append(samples, sample)
	}

	return &TrackedSampleCollection{Samples: samples}, nil
}

// recordToSample converts a single TSV record to a TrackedSample.
func recordToSample(record []string) (TrackedSample, error) {
	sample := TrackedSample{
		StudyID:             record[0],
		StudyName:           record[1],
		FacultySponsor:      record[2],
		Programme:           record[3],
		SangerSampleID:      record[4],
		SupplierName:        record[5],
		LabwareHumanBarcode: record[9],
		RunID:               record[14],
		Platform:            record[15],
		Pipeline:            record[16],
		QCPass:              record[20],
	}

	// Parse date fields
	sample.ManifestCreated = parseTime(record[6])
	sample.ManifestUploaded = parseTime(record[7])
	sample.LabwareReceived = parseTime(record[8])
	sample.OrderMade = parseTime(record[10])
	sample.LibraryStart = parseTime(record[11])
	sample.LibraryComplete = parseTime(record[12])
	sample.SequencingRunStart = parseTime(record[17])
	sample.SequencingQCComplete = parseTime(record[18])

	// Parse integer fields
	sample.LibraryTime = parseInt(record[13])
	sample.SequencingTime = parseInt(record[19])

	return sample, nil
}

// New creates a new QueryProvider with the appropriate implementation based on options.
func New(opts ...Option) (QueryProvider, error) {
	cfg := &config{
		useMock:   false,
		connector: &MySQLConnector{},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.useMock {
		if cfg.mockPath == "" {
			return nil, fmt.Errorf("mock data path not provided")
		}

		// Verify the file exists but don't try to open it yet
		if _, err := os.Stat(cfg.mockPath); err != nil {
			return nil, fmt.Errorf("mock data file error: %w", err)
		}

		return &MockQueryProvider{tsvPath: cfg.mockPath}, nil
	}

	return &MySQLQueryProvider{connector: cfg.connector}, nil
}

func parseTime(s string) *time.Time {
	if s == "" {
		return nil
	}

	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil
	}

	return &t
}

func parseInt(s string) *int {
	if s == "" {
		return nil
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}

	return &i
}
