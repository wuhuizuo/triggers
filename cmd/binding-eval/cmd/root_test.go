/*
Copyright 2020 The Tekton Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEvalBinding(t *testing.T) {
	out := new(bytes.Buffer)
	if err := evalBinding(out, "../testdata/triggerbinding.yaml", "../testdata/http.txt"); err != nil {
		t.Fatalf("evalBinding: %v", err)
	}

	want := `[
  {
    "name": "bar",
    "value": "tacocat"
  },
  {
    "name": "foo",
    "value": "body"
  }
]
`
	if diff := cmp.Diff(want, out.String()); diff != "" {
		t.Errorf("-want +got: %s", diff)
	}
}

func TestEvalBindingMissingContentLength(t *testing.T) {
	out := new(bytes.Buffer)
	if err := evalBinding(out, "../testdata/triggerbinding.yaml", "../testdata/http_no_content_length.txt"); err != nil {
		t.Fatalf("evalBinding with missing Content-Length: %v", err)
	}

	want := `[
  {
    "name": "bar",
    "value": "tacocat"
  },
  {
    "name": "foo",
    "value": "body"
  }
]
`
	if diff := cmp.Diff(want, out.String()); diff != "" {
		t.Errorf("-want +got: %s", diff)
	}
}

func TestReadHTTPAutoComputesContentLength(t *testing.T) {
	req, body, err := readHTTP("../testdata/http_no_content_length.txt")
	if err != nil {
		t.Fatalf("readHTTP: %v", err)
	}

	expectedBodyLength := len(`{"test": "body"}`)
	if req.ContentLength != int64(expectedBodyLength) {
		t.Errorf("Expected ContentLength to be auto-computed to %d, got %d", expectedBodyLength, req.ContentLength)
	}

	if req.Header.Get("Content-Length") != "16" {
		t.Errorf("Expected Content-Length header to be set to '16', got '%s'", req.Header.Get("Content-Length"))
	}

	if len(body) != expectedBodyLength {
		t.Errorf("Expected body length to be %d, got %d", expectedBodyLength, len(body))
	}

	expectedBody := `{"test": "body"}`
	if string(body) != expectedBody {
		t.Errorf("Expected body to be %q, got %q", expectedBody, string(body))
	}
}

func TestReadHTTPPreservesExistingContentLength(t *testing.T) {
	// Test with the original file that already has Content-Length: 16
	req, body, err := readHTTP("../testdata/http.txt")
	if err != nil {
		t.Fatalf("readHTTP: %v", err)
	}

	expectedBodyLength := len(`{"test": "body"}`)
	if req.ContentLength != int64(expectedBodyLength) {
		t.Errorf("Expected ContentLength to remain %d, got %d", expectedBodyLength, req.ContentLength)
	}

	if req.Header.Get("Content-Length") != "16" {
		t.Errorf("Expected Content-Length header to remain '16', got '%s'", req.Header.Get("Content-Length"))
	}

	if len(body) != expectedBodyLength {
		t.Errorf("Expected body length to be %d, got %d", expectedBodyLength, len(body))
	}

	expectedBody := `{"test": "body"}`
	if string(body) != expectedBody {
		t.Errorf("Expected body to be %q, got %q", expectedBody, string(body))
	}
}
