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

package db_test

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wtsi-hgi/gst/db"
)

type mockConnector struct {
	connectCalled bool
	closeCalled   bool
	mockDB        *sql.DB
	shouldFail    bool
}

func (m *mockConnector) Connect() (*sql.DB, error) {
	m.connectCalled = true
	if m.shouldFail {
		return nil, sql.ErrConnDone
	}
	return m.mockDB, nil
}

func (m *mockConnector) Close() error {
	m.closeCalled = true
	return nil
}

func TestQueryProvider(t *testing.T) {
	Convey("Given a MySQL query provider with a mock connector", t, func() {
		mockConn := &mockConnector{}
		provider, err := db.New(db.WithMySQLConnector(mockConn))

		So(err, ShouldBeNil)
		So(provider, ShouldNotBeNil)

		Convey("When executing a query", func() {
			// This will fail since we're not providing a real DB connection
			_, err := provider.Execute()

			Convey("Then the connector methods should be called", func() {
				So(mockConn.connectCalled, ShouldBeTrue)
				So(mockConn.closeCalled, ShouldBeTrue)
			})

			Convey("And an appropriate error about nil connection should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "nil")
			})
		})
	})

	Convey("Given a mock query provider", t, func() {
		tempDir := t.TempDir()

		mockPath := filepath.Join(tempDir, "mock_data.tsv")
		createMockTSVFile(mockPath)

		provider, err := db.New(db.WithMockData(mockPath))
		So(err, ShouldBeNil)
		So(provider, ShouldNotBeNil)

		Convey("When executing a query", func() {
			samples, err := provider.Execute()

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)
				So(samples, ShouldNotBeNil)
				So(len(samples.Samples), ShouldEqual, 1)
				So(samples.Samples[0].StudyID, ShouldEqual, "1234")
				So(samples.Samples[0].StudyName, ShouldEqual, "Test Study")
				So(samples.Samples[0].Platform, ShouldEqual, "Illumina")
			})
		})
	})

	Convey("Given no options", t, func() {
		provider, err := db.New()

		Convey("It should create a default MySQL provider", func() {
			So(err, ShouldBeNil)
			So(provider, ShouldNotBeNil)
			// We can't really test further without connecting to a real database
		})
	})
}

func createMockTSVFile(path string) {
	// Format: tab-delimited with header row matching the fields in TrackedSample
	data := `StudyID	StudyName	FacultySponsor	Programme	SangerSampleID	SupplierName	ManifestCreated	ManifestUploaded	LabwareReceived	Plate/Tube	OrderMade	LibraryStart	LibraryComplete	LibraryTime	RunID	Platform	Pipeline	SequencingRunStart	SequencingQCComplete	SequencingTime	QCPass
1234	Test Study	Test Sponsor	Test Programme	SANG123	Test Supplier	2025-01-01T12:00:00Z	2025-01-01T12:00:00Z	2025-01-01T12:00:00Z	PLATE001	2025-01-01T12:00:00Z	2025-01-01T12:00:00Z	2025-01-01T12:00:00Z	5	RUN001	Illumina	Pipeline1	2025-01-01T12:00:00Z	2025-01-01T12:00:00Z	10	1`

	err := os.WriteFile(path, []byte(data), 0644)
	if err != nil {
		panic(err)
	}
}
