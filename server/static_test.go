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
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wtsi-hgi/gst/db"
)

func TestStaticFilesExist(t *testing.T) {
	Convey("Given the embedded static files", t, func() {
		staticDir, err := fs.Sub(staticFiles, "static")

		Convey("We should be able to access the static directory", func() {
			So(err, ShouldBeNil)
			So(staticDir, ShouldNotBeNil)
		})

		Convey("The CSS file should exist", func() {
			_, err := fs.Stat(staticDir, "styles.css")
			So(err, ShouldBeNil)
		})

		Convey("The JS file should exist", func() {
			_, err := fs.Stat(staticDir, "script.js")
			So(err, ShouldBeNil)
		})
	})
}

func TestStaticFileServing(t *testing.T) {
	Convey("Given a server with static files", t, func() {
		mockProvider := &mockQueryProvider{
			samples: &db.TrackedSampleCollection{
				Samples: []db.TrackedSample{},
			},
		}

		srv, err := New(Config{
			QueryProvider: mockProvider,
			Port:          8080,
		})

		So(err, ShouldBeNil)
		So(srv, ShouldNotBeNil)

		Convey("When requesting the CSS file", func() {
			req := httptest.NewRequest("GET", "/static/styles.css", nil)
			resp := httptest.NewRecorder()

			srv.ServeHTTP(resp, req)

			Convey("It should return 200 OK", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
			})

			Convey("It should have the correct content type", func() {
				contentType := resp.Header().Get("Content-Type")
				So(contentType, ShouldEqual, "text/css; charset=utf-8")
			})

			Convey("It should contain CSS content", func() {
				body := resp.Body.String()
				So(body, ShouldContainSubstring, "body")
				So(body, ShouldContainSubstring, "font-family")
			})
		})

		Convey("When requesting the JavaScript file", func() {
			req := httptest.NewRequest("GET", "/static/script.js", nil)
			resp := httptest.NewRecorder()

			srv.ServeHTTP(resp, req)

			Convey("It should return 200 OK", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
			})

			Convey("It should have the correct content type", func() {
				contentType := resp.Header().Get("Content-Type")
				So(contentType, ShouldEqual, "application/javascript; charset=utf-8")
			})

			Convey("It should contain JavaScript content", func() {
				body := resp.Body.String()
				So(body, ShouldContainSubstring, "function")
				So(body, ShouldContainSubstring, "fetch")
			})
		})

		Convey("When requesting the index page", func() {
			req := httptest.NewRequest("GET", "/", nil)
			resp := httptest.NewRecorder()

			srv.ServeHTTP(resp, req)

			Convey("It should link to the external CSS and JS files", func() {
				body := resp.Body.String()
				So(body, ShouldContainSubstring, `<link rel="stylesheet" href="/static/styles.css">`)
				So(body, ShouldContainSubstring, `<script src="/static/script.js"></script>`)
				So(body, ShouldNotContainSubstring, `<style>`)
				So(body, ShouldNotContainSubstring, `function updateChart`)
			})
		})
	})
}
