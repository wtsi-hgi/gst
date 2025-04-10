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
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/wtsi-hgi/gst/db"
)

//go:embed static/*.html
var staticFiles embed.FS

// Config holds configuration options for the Server.
type Config struct {
	// QueryProvider is used to fetch sample data.
	QueryProvider db.QueryProvider

	// Port is the port to listen on.
	Port int
}

// Server handles HTTP requests for the sample tracking dashboard.
type Server struct {
	config        Config
	queryProvider db.QueryProvider
	templates     *template.Template
	mux           *http.ServeMux
}

// ChartData represents the data structure used for the Chart.js visualization.
type ChartData struct {
	Labels         []string `json:"labels"`
	SampleIds      []string `json:"sampleIds"`
	LibraryTime    []int    `json:"libraryTime"`
	SequencingTime []int    `json:"sequencingTime"`
}

// New creates a new Server with the given configuration.
func New(config Config) (*Server, error) {
	// Set default port if not specified
	if config.Port == 0 {
		config.Port = 8080
	}

	// Load and parse templates
	tmpl, err := template.ParseFS(staticFiles, "static/*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	// Create server
	server := &Server{
		config:        config,
		queryProvider: config.QueryProvider,
		templates:     tmpl,
		mux:           http.NewServeMux(),
	}

	// Register routes
	server.mux.HandleFunc("/", server.handleIndex)
	server.mux.HandleFunc("/api/samples", server.handleSamples)
	server.mux.HandleFunc("/api/chart", server.handleChart)

	return server, nil
}

// ServeHTTP implements the http.Handler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// handleIndex serves the main dashboard HTML page.
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	// Since we no longer need to pass auth credentials to the template,
	// we can just execute it without any data
	err := s.templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
	}
}

// handleSamples serves the HTML table of sample data.
func (s *Server) handleSamples(w http.ResponseWriter, r *http.Request) {
	// Get sample data
	samplesData, err := s.queryProvider.Execute()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving sample data: %v", err), http.StatusInternalServerError)
		return
	}

	// Render table template
	err = s.templates.ExecuteTemplate(w, "samples_table.html", samplesData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
	}
}

// handleChart provides JSON data for the Chart.js visualization.
func (s *Server) handleChart(w http.ResponseWriter, r *http.Request) {
	// Get sample data
	samplesData, err := s.queryProvider.Execute()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving sample data: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare chart data
	chartData := prepareChartData(samplesData.Samples)

	// Return as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chartData); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
	}
}

// prepareChartData converts sample data into a format suitable for Chart.js.
func prepareChartData(samples []db.TrackedSample) ChartData {
	chartData := ChartData{
		Labels:         make([]string, 0, len(samples)),
		SampleIds:      make([]string, 0, len(samples)),
		LibraryTime:    make([]int, 0, len(samples)),
		SequencingTime: make([]int, 0, len(samples)),
	}

	for _, sample := range samples {
		chartData.Labels = append(chartData.Labels, sample.SangerSampleID)
		chartData.SampleIds = append(chartData.SampleIds, sample.SangerSampleID)

		// Handle potential nil values
		libTime := 0
		if sample.LibraryTime != nil {
			libTime = *sample.LibraryTime
		}
		chartData.LibraryTime = append(chartData.LibraryTime, libTime)

		seqTime := 0
		if sample.SequencingTime != nil {
			seqTime = *sample.SequencingTime
		}
		chartData.SequencingTime = append(chartData.SequencingTime, seqTime)
	}

	return chartData
}

// Start starts the HTTP server on the configured port.
func (s *Server) Start() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", s.config.Port), s)
}
