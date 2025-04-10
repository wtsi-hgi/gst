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

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wtsi-hgi/gst/db"
)

func main() {
	outputPath := flag.String("output", "samples.tsv", "Path to output TSV file")
	useMock := flag.Bool("mock", false, "Use mock data instead of querying the database")
	mockPath := flag.String("mockPath", "db/sampledata.tsv", "Path to mock data TSV file")
	flag.Parse()

	var provider db.QueryProvider
	var err error

	if *useMock {
		provider, err = db.New(db.WithMockData(*mockPath))
	} else {
		fmt.Println("Executing database query. This may take several minutes...")
		provider, err = db.New()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating query provider: %v\n", err)
		os.Exit(1)
	}

	samples, err := provider.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing query: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Retrieved %d sample records\n", len(samples.Samples))

	// Ensure output directory exists
	outputDir := filepath.Dir(*outputPath)
	if outputDir != "" && outputDir != "." {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Write results to TSV
	if err := samples.ToTSV(*outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing TSV file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Results written to %s\n", *outputPath)
}
