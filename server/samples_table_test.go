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

package server

import (
	"bytes"
	"html/template"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wtsi-hgi/gst/db"
)

func TestSamplesTableTemplate(t *testing.T) {
	Convey("Given a samples table template", t, func() {
		tmpl, err := template.ParseFiles("static/samples_table.html")
		So(err, ShouldBeNil)

		Convey("When rendering with sample data", func() {
			// Create a sample with all fields populated
			now := time.Now()
			libraryTime := 5
			seqTime := 10
			sample := db.TrackedSample{
				StudyID:              "1234",
				StudyName:            "Test Study",
				FacultySponsor:       "Test Sponsor",
				Programme:            "Test Programme",
				SangerSampleID:       "SANG123",
				SupplierName:         "Test Supplier",
				ManifestCreated:      &now,
				ManifestUploaded:     &now,
				LabwareReceived:      &now,
				LabwareHumanBarcode:  "PLATE001",
				OrderMade:            &now,
				LibraryStart:         &now,
				LibraryComplete:      &now,
				LibraryTime:          &libraryTime,
				RunID:                "RUN001",
				Platform:             "Illumina",
				Pipeline:             "Pipeline1",
				SequencingRunStart:   &now,
				SequencingQCComplete: &now,
				SequencingTime:       &seqTime,
				QCPass:               "1",
			}

			data := struct {
				HasData bool
				Samples []db.TrackedSample
			}{
				HasData: true,
				Samples: []db.TrackedSample{sample},
			}

			var buf bytes.Buffer
			err := tmpl.Execute(&buf, data)
			So(err, ShouldBeNil)

			output := buf.String()

			Convey("It should contain all required fields from TrackedSample", func() {
				// Check for existing fields
				So(output, ShouldContainSubstring, "Sanger Sample ID")
				So(output, ShouldContainSubstring, "Supplier Name")
				So(output, ShouldContainSubstring, "Manifest Created")
				So(output, ShouldContainSubstring, "Plate/Tube")
				So(output, ShouldContainSubstring, "Library Time")
				So(output, ShouldContainSubstring, "Run ID")
				So(output, ShouldContainSubstring, "Platform")
				So(output, ShouldContainSubstring, "Sequencing Time")
				So(output, ShouldContainSubstring, "QC Pass")

				// Check for new fields
				So(output, ShouldContainSubstring, "Manifest Uploaded")
				So(output, ShouldContainSubstring, "Labware Received")
				So(output, ShouldContainSubstring, "Order Made")
				So(output, ShouldContainSubstring, "Library Start")
				So(output, ShouldContainSubstring, "Library Complete")
				So(output, ShouldContainSubstring, "Pipeline")
				So(output, ShouldContainSubstring, "Sequencing Run Start")
				So(output, ShouldContainSubstring, "Sequencing QC Complete")

				// Check that data is displayed
				So(output, ShouldContainSubstring, "SANG123")
				So(output, ShouldContainSubstring, "Test Supplier")
				So(output, ShouldContainSubstring, "PLATE001")
				So(output, ShouldContainSubstring, "RUN001")
				So(output, ShouldContainSubstring, "Illumina")
				So(output, ShouldContainSubstring, "Pipeline1")
			})

			Convey("It should not contain the excluded fields", func() {
				So(output, ShouldNotContainSubstring, "<th>Study ID</th>")
				So(output, ShouldNotContainSubstring, "<th>Study Name</th>")
				So(output, ShouldNotContainSubstring, "<th>Faculty Sponsor</th>")
				So(output, ShouldNotContainSubstring, "<th>Programme</th>")
			})
		})
	})
}
