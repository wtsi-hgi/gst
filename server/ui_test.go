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
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wtsi-hgi/gst/db"
)

func TestChartData(t *testing.T) {
	Convey("When preparing chart data", t, func() {
		// Create test samples
		samples := createTestSamples(5)

		Convey("Chart data should include supplier names", func() {
			chartData := prepareChartData(samples)

			// Verify labels use supplier names
			So(len(chartData.Labels), ShouldEqual, 5)
			So(chartData.Labels[0], ShouldEqual, "Test Supplier 1")
			So(chartData.Labels[1], ShouldEqual, "Test Supplier 2")

			// Verify sample IDs are still included for tooltips
			So(len(chartData.SampleIds), ShouldEqual, 5)
			So(chartData.SampleIds[0], ShouldEqual, "SANG001")
			So(chartData.SampleIds[1], ShouldEqual, "SANG002")
		})
	})
}

// Helper function to create test samples
func createTestSamples(count int) []db.TrackedSample {
	samples := make([]db.TrackedSample, count)
	for i := 0; i < count; i++ {
		libraryTime := 5 + i
		seqTime := 10 + i

		samples[i] = db.TrackedSample{
			SangerSampleID: fmt.Sprintf("SANG%03d", i+1),
			SupplierName:   fmt.Sprintf("Test Supplier %d", i+1),
			LibraryTime:    &libraryTime,
			SequencingTime: &seqTime,
		}
	}
	return samples
}
