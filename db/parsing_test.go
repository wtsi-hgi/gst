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
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestParseRows(t *testing.T) {
	Convey("Given a set of SQL rows with NULL values", t, func() {
		// Create a mock database and connection
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		defer db.Close()

		// Define expected columns that match our SELECT query
		columns := []string{
			"study_id", "StudyName", "faculty_sponsor", "programme",
			"sanger_sample_id", "supplier_name", "manifest_created",
			"manifest_uploaded", "labware_received", "Plate/Tube",
			"order_made", "library_start", "library_complete", "LibraryTime",
			"RunID", "Platform", "Pipeline", "sequencing_run_start",
			"sequencing_qc_complete", "SequencingTime", "qcPass",
		}

		// Create mock rows with some NULL values
		mock.ExpectQuery("SELECT").WillReturnRows(
			mock.NewRows(columns).
				AddRow(
					"12345", "Test Study", "Test Sponsor", "Test Programme",
					"SANG123", "Test Supplier", nil, nil,
					nil, "PLATE001", nil, nil,
					nil, nil, // End of first 14 columns
					nil,                          // NULL for RunID (column 15)
					nil, nil, nil, nil, nil, nil, // Remaining columns
				),
		)

		// Execute query
		rows, err := db.Query("SELECT 1")
		So(err, ShouldBeNil)
		defer rows.Close()

		// Parse rows with our function
		result, err := parseRows(rows)

		Convey("It should handle NULL values without error", func() {
			So(err, ShouldBeNil)
			So(result, ShouldNotBeNil)
			So(len(result.Samples), ShouldEqual, 1)

			sample := result.Samples[0]
			So(sample.StudyID, ShouldEqual, "12345")
			So(sample.StudyName, ShouldEqual, "Test Study")
			So(sample.RunID, ShouldEqual, "")
		})
	})
}
