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
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wtsi-hgi/gst/db"
)

func TestFilter(t *testing.T) {
	Convey("Given a filter with sample data", t, func() {
		samples := []db.TrackedSample{
			{
				StudyID:        "1234",
				StudyName:      "Study A",
				FacultySponsor: "Sponsor 1",
				Programme:      "Programme X",
				SangerSampleID: "SANG123",
			},
			{
				StudyID:        "5678",
				StudyName:      "Study B",
				FacultySponsor: "Sponsor 1",
				Programme:      "Programme Y",
				SangerSampleID: "SANG456",
			},
			{
				StudyID:        "9012",
				StudyName:      "Study C",
				FacultySponsor: "Sponsor 2",
				Programme:      "Programme Z",
				SangerSampleID: "SANG789",
			},
		}

		Convey("When getting unique faculty sponsors", func() {
			sponsors := GetUniqueFacultySponsors(samples)

			Convey("It should return all unique values", func() {
				So(len(sponsors), ShouldEqual, 2)
				So(sponsors, ShouldContain, "Sponsor 1")
				So(sponsors, ShouldContain, "Sponsor 2")
			})
		})

		Convey("When getting study names for a sponsor", func() {
			studies := GetStudiesForSponsor(samples, "Sponsor 1")

			Convey("It should return only studies for that sponsor", func() {
				So(len(studies), ShouldEqual, 2)
				So(studies, ShouldContain, "Study A")
				So(studies, ShouldContain, "Study B")
				So(studies, ShouldNotContain, "Study C")
			})
		})

		Convey("When filtering samples by sponsor and study", func() {
			filtered := FilterSamples(samples, "Sponsor 1", "Study B")

			Convey("It should return only matching samples", func() {
				So(len(filtered), ShouldEqual, 1)
				So(filtered[0].StudyName, ShouldEqual, "Study B")
				So(filtered[0].SangerSampleID, ShouldEqual, "SANG456")
			})

			Convey("When filtering by sponsor only", func() {
				filtered := FilterSamples(samples, "Sponsor 2", "")

				Convey("It should return all samples for that sponsor", func() {
					So(len(filtered), ShouldEqual, 1)
					So(filtered[0].StudyName, ShouldEqual, "Study C")
					So(filtered[0].SangerSampleID, ShouldEqual, "SANG789")
				})
			})
		})
	})
}

func TestFilterFunctions(t *testing.T) {
	Convey("Given a set of sample data with multiple sponsors and studies", t, func() {
		samples := []db.TrackedSample{
			{FacultySponsor: "Sponsor A", StudyName: "Study 1"},
			{FacultySponsor: "Sponsor A", StudyName: "Study 2"},
			{FacultySponsor: "Sponsor B", StudyName: "Study 3"},
			{FacultySponsor: "Sponsor C", StudyName: "Study 4"},
			{FacultySponsor: "Sponsor A", StudyName: "Study 1"}, // Duplicate
		}

		Convey("GetUniqueFacultySponsors should return all unique sponsors", func() {
			sponsors := GetUniqueFacultySponsors(samples)
			So(len(sponsors), ShouldEqual, 3)
			So(sponsors, ShouldContain, "Sponsor A")
			So(sponsors, ShouldContain, "Sponsor B")
			So(sponsors, ShouldContain, "Sponsor C")
		})

		Convey("GetStudiesForSponsor should return studies for the specified sponsor", func() {
			studies := GetStudiesForSponsor(samples, "Sponsor A")
			So(len(studies), ShouldEqual, 2)
			So(studies, ShouldContain, "Study 1")
			So(studies, ShouldContain, "Study 2")
		})
	})
}
