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
	"time"

	"github.com/joho/godotenv"
	"github.com/wtsi-hgi/gst/db"
	"github.com/wtsi-hgi/gst/server"
)

func main() {
	// Load environment variables from .env file
	godotenv.Load()

	// Subcommands
	exportCmd := flag.NewFlagSet("export", flag.ExitOnError)
	serverCmd := flag.NewFlagSet("server", flag.ExitOnError)

	// Export command flags
	outputPath := exportCmd.String("output", "samples.tsv", "Path to output TSV file")

	// Server command flags
	serverPort := serverCmd.Int("port", 8080, "Port to run the server on")
	serverMockPath := serverCmd.String("mock", "samples.tsv", "Path to mock data TSV file")
	cacheTTL := serverCmd.Duration("cacheTTL", 5*time.Minute, "Duration to cache data before refreshing")

	// Check which subcommand is being used
	if len(os.Args) < 2 {
		fmt.Println("Expected 'export' or 'server' subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "export":
		exportCmd.Parse(os.Args[2:])
		runExport(outputPath)
	case "server":
		serverCmd.Parse(os.Args[2:])
		runServer(serverPort, serverMockPath, cacheTTL)
	default:
		fmt.Println("Expected 'export' or 'server' subcommand")
		os.Exit(1)
	}
}

func runExport(outputPath *string) {
	fmt.Println("Executing database query. This may take several minutes...")
	provider, err := db.New()
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

func runServer(port *int, mockPath *string, cacheTTL *time.Duration) {
	// Create query provider
	var provider db.QueryProvider
	var err error

	if *mockPath != "" {
		provider, err = db.New(db.WithMockData(*mockPath))
	} else {
		provider, err = db.New()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating query provider: %v\n", err)
		os.Exit(1)
	}

	// Create and start server
	srv, err := server.New(server.Config{
		QueryProvider: provider,
		Port:          *port,
		CacheTTL:      *cacheTTL,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating server: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Starting server on port %d...\n", *port)

	if err := srv.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting server: %v\n", err)
		os.Exit(1)
	}
}
