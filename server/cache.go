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
	"sort"
	"sync"
	"time"

	"github.com/wtsi-hgi/gst/db"
)

// Cache provides a time-based caching mechanism for sample data.
type Cache struct {
	provider    db.QueryProvider
	ttl         time.Duration
	samples     *db.TrackedSampleCollection
	lastFetched time.Time
	mu          sync.RWMutex
}

// NewCache creates a new cache with the specified provider and TTL.
func NewCache(provider db.QueryProvider, ttl time.Duration) *Cache {
	return &Cache{
		provider: provider,
		ttl:      ttl,
	}
}

// GetSamples returns sample data, either from the cache if it's still valid
// or by fetching fresh data from the provider.
func (c *Cache) GetSamples() (*db.TrackedSampleCollection, error) {
	c.mu.RLock()
	if c.samples != nil && time.Since(c.lastFetched) < c.ttl {
		defer c.mu.RUnlock()
		return c.samples, nil
	}
	c.mu.RUnlock()

	// Need to refresh the cache
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check again in case another goroutine refreshed while we were waiting
	if c.samples != nil && time.Since(c.lastFetched) < c.ttl {
		return c.samples, nil
	}

	// Fetch fresh data
	samples, err := c.provider.Execute()
	if err != nil {
		return nil, err
	}

	c.samples = samples
	c.lastFetched = time.Now()
	return c.samples, nil
}

// GetUniqueFacultySponsors returns a sorted list of unique faculty sponsors.
func GetUniqueFacultySponsors(samples []db.TrackedSample) []string {
	sponsorMap := make(map[string]struct{})

	for _, sample := range samples {
		if sample.FacultySponsor != "" {
			sponsorMap[sample.FacultySponsor] = struct{}{}
		}
	}

	sponsors := make([]string, 0, len(sponsorMap))
	for sponsor := range sponsorMap {
		sponsors = append(sponsors, sponsor)
	}

	sort.Strings(sponsors)
	return sponsors
}

// GetStudiesForSponsor returns a sorted list of study names for a given sponsor.
func GetStudiesForSponsor(samples []db.TrackedSample, sponsor string) []string {
	studyMap := make(map[string]struct{})

	for _, sample := range samples {
		if sample.FacultySponsor == sponsor && sample.StudyName != "" {
			studyMap[sample.StudyName] = struct{}{}
		}
	}

	studies := make([]string, 0, len(studyMap))
	for study := range studyMap {
		studies = append(studies, study)
	}

	sort.Strings(studies)
	return studies
}

// FilterSamples filters samples by faculty sponsor and optionally by study name.
func FilterSamples(samples []db.TrackedSample, sponsor, study string) []db.TrackedSample {
	if sponsor == "" {
		return samples
	}

	var filtered []db.TrackedSample

	for _, sample := range samples {
		if sample.FacultySponsor != sponsor {
			continue
		}

		if study != "" && sample.StudyName != study {
			continue
		}

		filtered = append(filtered, sample)
	}

	return filtered
}
