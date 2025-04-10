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

package server_test

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wtsi-hgi/gst/db"
	"github.com/wtsi-hgi/gst/server"
)

// mockQueryProvider implements db.QueryProvider for testing.
type mockQueryProvider struct {
	samples *db.TrackedSampleCollection
	err     error
}

func (m *mockQueryProvider) Execute() (*db.TrackedSampleCollection, error) {
	return m.samples, m.err
}

// getAvailablePort returns a random available port by asking the OS to assign one.
func getAvailablePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port, nil
}

func TestServer(t *testing.T) {
	Convey("Given a server with mock data", t, func() {
		// Find an available port
		port, err := getAvailablePort()
		So(err, ShouldBeNil)

		// Create mock sample data
		sampleTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
		libraryTime := 5
		seqTime := 10

		mockSamples := &db.TrackedSampleCollection{
			Samples: []db.TrackedSample{
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
				{
					StudyID:        "5678",
					StudyName:      "Another Study",
					FacultySponsor: "Another Sponsor",
					Programme:      "Test Programme",
					SangerSampleID: "SANG456",
					SupplierName:   "Test Supplier",
					LibraryTime:    &libraryTime,
					SequencingTime: &seqTime,
				},
			},
		}

		mockProvider := &mockQueryProvider{samples: mockSamples}

		// Create server with test config and random port
		srv, err := server.New(server.Config{
			QueryProvider: mockProvider,
			Port:          port,
		})

		Convey("It should initialize without error", func() {
			So(err, ShouldBeNil)
			So(srv, ShouldNotBeNil)
		})

		Convey("When requesting the index page", func() {
			req := httptest.NewRequest("GET", "/", nil)
			resp := httptest.NewRecorder()

			srv.ServeHTTP(resp, req)

			Convey("It should return 200 OK", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
			})

			Convey("It should contain HTML with HTMX", func() {
				body := resp.Body.String()
				So(body, ShouldContainSubstring, "<html")
				So(body, ShouldContainSubstring, "htmx")
				So(body, ShouldContainSubstring, "chart.js")
			})
		})

		Convey("When requesting the faculty sponsors for filtering", func() {
			req := httptest.NewRequest("GET", "/api/filters", nil)
			resp := httptest.NewRecorder()

			srv.ServeHTTP(resp, req)

			Convey("It should return 200 OK", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
			})

			Convey("It should return all unique sponsors as JSON", func() {
				var response struct {
					FacultySponsors []string `json:"facultySponsors"`
				}
				err := json.Unmarshal(resp.Body.Bytes(), &response)
				So(err, ShouldBeNil)

				So(len(response.FacultySponsors), ShouldEqual, 2)
				So(response.FacultySponsors, ShouldContain, "Test Sponsor")
				So(response.FacultySponsors, ShouldContain, "Another Sponsor")
			})
		})

		Convey("When requesting studies for a sponsor", func() {
			req := httptest.NewRequest("GET", "/api/studies?sponsor=Test+Sponsor", nil)
			resp := httptest.NewRecorder()

			srv.ServeHTTP(resp, req)

			Convey("It should return 200 OK", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
			})

			Convey("It should return studies for that sponsor as JSON", func() {
				var response struct {
					Studies []string `json:"studies"`
				}
				err := json.Unmarshal(resp.Body.Bytes(), &response)
				So(err, ShouldBeNil)

				So(len(response.Studies), ShouldEqual, 1)
				So(response.Studies, ShouldContain, "Test Study")
			})
		})

		Convey("When requesting the sample data API endpoint without required parameters", func() {
			req := httptest.NewRequest("GET", "/api/samples", nil)
			resp := httptest.NewRecorder()

			srv.ServeHTTP(resp, req)

			Convey("It should return 200 OK", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
			})

			Convey("It should indicate that sponsor and study selection is required", func() {
				body := resp.Body.String()
				So(body, ShouldContainSubstring, "Please select")
				So(body, ShouldNotContainSubstring, "Test Study")
			})
		})

		Convey("When requesting the sample data with sponsor and study parameters", func() {
			req := httptest.NewRequest("GET", "/api/samples?sponsor=Test+Sponsor&study=Test+Study", nil)
			resp := httptest.NewRecorder()

			srv.ServeHTTP(resp, req)

			Convey("It should return 200 OK", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
			})

			Convey("It should contain the filtered sample data", func() {
				body := resp.Body.String()
				So(body, ShouldContainSubstring, "SANG123")
				So(body, ShouldNotContainSubstring, "SANG456")
			})
		})

		Convey("When requesting the chart data API endpoint without required parameters", func() {
			req := httptest.NewRequest("GET", "/api/chart", nil)
			resp := httptest.NewRecorder()

			srv.ServeHTTP(resp, req)

			Convey("It should return 200 OK", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
			})

			Convey("It should return empty chart data", func() {
				var chartData struct {
					Labels []string `json:"labels"`
				}
				err := json.Unmarshal(resp.Body.Bytes(), &chartData)
				So(err, ShouldBeNil)
				So(len(chartData.Labels), ShouldEqual, 0)
			})
		})

		Convey("When requesting the chart data with sponsor and study parameters", func() {
			req := httptest.NewRequest("GET", "/api/chart?sponsor=Test+Sponsor&study=Test+Study", nil)
			resp := httptest.NewRecorder()

			srv.ServeHTTP(resp, req)

			Convey("It should return 200 OK", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
			})

			Convey("It should contain chart data for the filtered samples", func() {
				body := resp.Body.String()
				So(body, ShouldContainSubstring, "libraryTime")
				So(body, ShouldContainSubstring, "sequencingTime")
				So(body, ShouldContainSubstring, "SANG123")
				So(body, ShouldContainSubstring, "5")  // LibraryTime value
				So(body, ShouldContainSubstring, "10") // SequencingTime value
				So(body, ShouldNotContainSubstring, "SANG456")
			})
		})
	})
}
