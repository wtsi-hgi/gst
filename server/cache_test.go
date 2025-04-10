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
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wtsi-hgi/gst/db"
)

// mockQueryProvider implements db.QueryProvider for testing.
type mockQueryProvider struct {
	samples      *db.TrackedSampleCollection
	err          error
	executeCalls int
}

func (m *mockQueryProvider) Execute() (*db.TrackedSampleCollection, error) {
	m.executeCalls++
	return m.samples, m.err
}

func TestCache(t *testing.T) {
	Convey("Given a cache with a mock provider", t, func() {
		mockSamples := &db.TrackedSampleCollection{
			Samples: []db.TrackedSample{
				{
					StudyID:        "1234",
					StudyName:      "Test Study",
					FacultySponsor: "Test Sponsor",
				},
			},
		}

		mockProvider := &mockQueryProvider{samples: mockSamples}

		// Create cache with a short TTL for testing (100ms)
		cache := NewCache(mockProvider, 100*time.Millisecond)

		Convey("When getting samples for the first time", func() {
			samples, err := cache.GetSamples()

			Convey("It should fetch from the provider", func() {
				So(err, ShouldBeNil)
				So(samples, ShouldNotBeNil)
				So(mockProvider.executeCalls, ShouldEqual, 1)
			})

			Convey("When getting samples again immediately", func() {
				samples2, err := cache.GetSamples()

				Convey("It should use the cache and not call the provider again", func() {
					So(err, ShouldBeNil)
					So(samples2, ShouldNotBeNil)
					So(mockProvider.executeCalls, ShouldEqual, 1) // Still 1
				})
			})

			Convey("When getting samples after TTL expires", func() {
				time.Sleep(150 * time.Millisecond) // Wait longer than TTL
				samples2, err := cache.GetSamples()

				Convey("It should fetch from the provider again", func() {
					So(err, ShouldBeNil)
					So(samples2, ShouldNotBeNil)
					So(mockProvider.executeCalls, ShouldEqual, 2) // Now 2
				})
			})
		})
	})
}
