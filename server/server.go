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
	"time"

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

	// CacheTTL is how long to cache data before refreshing.
	CacheTTL time.Duration
}

// Server handles HTTP requests for the sample tracking dashboard.
type Server struct {
	config    Config
	cache     *Cache
	templates *template.Template
	mux       *http.ServeMux
}

// ChartData represents the data structure used for the Chart.js visualization.
type ChartData struct {
	Labels         []string `json:"labels"`
	SampleIds      []string `json:"sampleIds"`
	LibraryTime    []int    `json:"libraryTime"`
	SequencingTime []int    `json:"sequencingTime"`
}

// FilterResponse contains data for populating filter dropdowns.
type FilterResponse struct {
	FacultySponsors []string `json:"facultySponsors"`
	Studies         []string `json:"studies,omitempty"`
}

// New creates a new Server with the given configuration.
func New(config Config) (*Server, error) {
	// Set default port if not specified
	if config.Port == 0 {
		config.Port = 8080
	}

	// Set default cache TTL if not specified
	if config.CacheTTL == 0 {
		config.CacheTTL = 5 * time.Minute
	}

	// Load and parse templates
	tmpl, err := template.ParseFS(staticFiles, "static/*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	// Create cache
	cache := NewCache(config.QueryProvider, config.CacheTTL)

	// Create server
	server := &Server{
		config:    config,
		cache:     cache,
		templates: tmpl,
		mux:       http.NewServeMux(),
	}

	// Register routes
	server.mux.HandleFunc("/", server.handleIndex)
	server.mux.HandleFunc("/api/samples", server.handleSamples)
	server.mux.HandleFunc("/api/chart", server.handleChart)
	server.mux.HandleFunc("/api/filters", server.handleFilters)
	server.mux.HandleFunc("/api/studies", server.handleStudies)

	return server, nil
}

// ServeHTTP implements the http.Handler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// handleIndex serves the main dashboard HTML page.
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	err := s.templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering template: %v", err),
			http.StatusInternalServerError)
	}
}

// handleSamples serves the HTML table of sample data.
func (s *Server) handleSamples(w http.ResponseWriter, r *http.Request) {
	// Get required filter parameters
	sponsor := r.URL.Query().Get("sponsor")
	study := r.URL.Query().Get("study")

	// Ensure both filters are provided
	if sponsor == "" || study == "" {
		// Return template with HasData = false
		templateData := struct {
			HasData bool
			Samples []db.TrackedSample
		}{
			HasData: false,
			Samples: nil,
		}

		err := s.templates.ExecuteTemplate(w, "samples_table.html", templateData)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error rendering template: %v", err),
				http.StatusInternalServerError)
		}
		return
	}

	// Get sample data from cache
	samplesData, err := s.cache.GetSamples()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving sample data: %v", err),
			http.StatusInternalServerError)
		return
	}

	// Apply filters (now both are required)
	filteredSamples := FilterSamples(samplesData.Samples, sponsor, study)

	// Create template data
	templateData := struct {
		HasData bool
		Samples []db.TrackedSample
	}{
		HasData: true,
		Samples: filteredSamples,
	}

	// Render table template
	err = s.templates.ExecuteTemplate(w, "samples_table.html", templateData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering template: %v", err),
			http.StatusInternalServerError)
	}
}

// handleChart provides JSON data for the Chart.js visualization.
func (s *Server) handleChart(w http.ResponseWriter, r *http.Request) {
	// Get required filter parameters
	sponsor := r.URL.Query().Get("sponsor")
	study := r.URL.Query().Get("study")

	// Return empty chart data if filters not provided
	if sponsor == "" || study == "" {
		emptyChart := ChartData{
			Labels:         []string{},
			SampleIds:      []string{},
			LibraryTime:    []int{},
			SequencingTime: []int{},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(emptyChart); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err),
				http.StatusInternalServerError)
		}
		return
	}

	// Get sample data from cache
	samplesData, err := s.cache.GetSamples()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving sample data: %v", err),
			http.StatusInternalServerError)
		return
	}

	// Apply filters
	filteredSamples := FilterSamples(samplesData.Samples, sponsor, study)

	// Prepare chart data
	chartData := prepareChartData(filteredSamples)

	// Return as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chartData); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err),
			http.StatusInternalServerError)
	}
}

// handleFilters provides a list of faculty sponsors for filtering.
func (s *Server) handleFilters(w http.ResponseWriter, r *http.Request) {
	// Get sample data from cache
	samplesData, err := s.cache.GetSamples()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving sample data: %v", err),
			http.StatusInternalServerError)
		return
	}

	fmt.Printf("in handleFilters, got %d samples from cache\n", len(samplesData.Samples))

	// Get unique faculty sponsors
	sponsors := GetUniqueFacultySponsors(samplesData.Samples)

	fmt.Printf("in handleFilters, got %d unique sponsors\n", len(sponsors))

	// Create response with explicitly initialized array
	response := FilterResponse{
		FacultySponsors: make([]string, len(sponsors)),
	}
	copy(response.FacultySponsors, sponsors)

	// Set proper headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")

	// Marshal directly for better control
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err),
			http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

// handleStudies provides a list of studies for a given faculty sponsor.
func (s *Server) handleStudies(w http.ResponseWriter, r *http.Request) {
	// Get sponsor parameter
	sponsor := r.URL.Query().Get("sponsor")
	if sponsor == "" {
		http.Error(w, "Sponsor parameter is required", http.StatusBadRequest)
		return
	}

	// Get sample data from cache
	samplesData, err := s.cache.GetSamples()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving sample data: %v", err),
			http.StatusInternalServerError)
		return
	}

	fmt.Printf("in handleStudies, got %d samples from cache\n", len(samplesData.Samples))

	// Get studies for this sponsor
	studies := GetStudiesForSponsor(samplesData.Samples, sponsor)

	// Create response with explicitly initialized array
	response := FilterResponse{
		Studies: make([]string, len(studies)),
	}
	copy(response.Studies, studies)

	// Set proper headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")

	// Marshal directly for better control
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err),
			http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
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
