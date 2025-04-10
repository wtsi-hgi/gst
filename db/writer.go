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
	"os"
	"strconv"
	"time"
)

// ToTSV writes the collection of samples to a TSV file at the specified path.
func (sc *TrackedSampleCollection) ToTSV(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	writer.Comma = '\t'
	defer writer.Flush()

	// Write header
	header := []string{
		"StudyID",
		"StudyName",
		"FacultySponsor",
		"Programme",
		"SangerSampleID",
		"SupplierName",
		"ManifestCreated",
		"ManifestUploaded",
		"LabwareReceived",
		"Plate/Tube",
		"OrderMade",
		"LibraryStart",
		"LibraryComplete",
		"LibraryTime",
		"RunID",
		"Platform",
		"Pipeline",
		"SequencingRunStart",
		"SequencingQCComplete",
		"SequencingTime",
		"QCPass",
	}

	if err := writer.Write(header); err != nil {
		return err
	}

	// Format and write each sample
	for _, sample := range sc.Samples {
		row := []string{
			sample.StudyID,
			sample.StudyName,
			sample.FacultySponsor,
			sample.Programme,
			sample.SangerSampleID,
			sample.SupplierName,
			formatTimePointer(sample.ManifestCreated),
			formatTimePointer(sample.ManifestUploaded),
			formatTimePointer(sample.LabwareReceived),
			sample.LabwareHumanBarcode,
			formatTimePointer(sample.OrderMade),
			formatTimePointer(sample.LibraryStart),
			formatTimePointer(sample.LibraryComplete),
			formatIntPointer(sample.LibraryTime),
			sample.RunID,
			sample.Platform,
			sample.Pipeline,
			formatTimePointer(sample.SequencingRunStart),
			formatTimePointer(sample.SequencingQCComplete),
			formatIntPointer(sample.SequencingTime),
			sample.QCPass,
		}

		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// Helper functions to format pointers for output.
func formatTimePointer(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

func formatIntPointer(i *int) string {
	if i == nil {
		return ""
	}
	return strconv.Itoa(*i)
}
