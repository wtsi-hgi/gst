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
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wtsi-hgi/gst/db"
)

func TestWriter(t *testing.T) {
	Convey("Given a collection of TrackedSample records", t, func() {
		// Create test data
		sampleTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
		libraryTime := 5
		seqTime := 10
		samples := []db.TrackedSample{
			{
				StudyID:              "1234",
				StudyName:            "Test Study",
				FacultySponsor:       "Test Sponsor",
				Programme:            "Test Programme",
				SangerSampleID:       "SANG123",
				SupplierName:         "Test Supplier",
				ManifestCreated:      &sampleTime,
				ManifestUploaded:     &sampleTime,
				LabwareReceived:      &sampleTime,
				LabwareHumanBarcode:  "PLATE001",
				OrderMade:            &sampleTime,
				LibraryStart:         &sampleTime,
				LibraryComplete:      &sampleTime,
				LibraryTime:          &libraryTime,
				RunID:                "RUN001",
				Platform:             "Illumina",
				Pipeline:             "Pipeline1",
				SequencingRunStart:   &sampleTime,
				SequencingQCComplete: &sampleTime,
				SequencingTime:       &seqTime,
				QCPass:               "1",
			},
		}

		collection := db.TrackedSampleCollection{Samples: samples}

		Convey("When writing to a TSV file", func() {
			tmpDir, err := os.MkdirTemp("", "gst_test")
			So(err, ShouldBeNil)
			defer os.RemoveAll(tmpDir)

			outPath := filepath.Join(tmpDir, "output.tsv")
			err = collection.ToTSV(outPath)

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)

				Convey("And the file should exist with expected content", func() {
					content, err := os.ReadFile(outPath)
					So(err, ShouldBeNil)

					// Check header and at least one data row
					lines := strings.Split(string(content), "\n")
					So(len(lines), ShouldBeGreaterThanOrEqualTo, 2)

					// Check header has all expected columns
					header := lines[0]
					So(header, ShouldContainSubstring, "StudyID")
					So(header, ShouldContainSubstring, "StudyName")
					So(header, ShouldContainSubstring, "FacultySponsor")
					So(header, ShouldContainSubstring, "Platform")

					// Check data is present
					data := lines[1]
					So(data, ShouldContainSubstring, "1234")
					So(data, ShouldContainSubstring, "Test Study")
					So(data, ShouldContainSubstring, "SANG123")
					So(data, ShouldContainSubstring, "Illumina")
				})
			})
		})
	})
}
