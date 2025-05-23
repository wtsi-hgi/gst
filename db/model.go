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
	"time"
)

// TrackedSample represents a row in the query results.
type TrackedSample struct {
	StudyID              string
	StudyName            string
	FacultySponsor       string
	Programme            string
	SangerSampleID       string
	SupplierName         string
	ManifestCreated      *time.Time
	ManifestUploaded     *time.Time
	LabwareReceived      *time.Time
	LabwareHumanBarcode  string // "Plate/Tube"
	OrderMade            *time.Time
	LibraryStart         *time.Time
	LibraryComplete      *time.Time
	LibraryTime          *int // DATEDIFF result
	RunID                string
	Platform             string
	Pipeline             string
	SequencingRunStart   *time.Time
	SequencingQCComplete *time.Time
	SequencingTime       *int // DATEDIFF result
	QCPass               string
}

// TrackedSampleCollection represents a collection of query results.
type TrackedSampleCollection struct {
	Samples []TrackedSample
}
