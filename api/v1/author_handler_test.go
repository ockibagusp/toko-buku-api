package v1

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	. "toko-buku-api/api/v1/mock"
)

func TestGetAuthors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/authors" {
			t.Errorf("Expected to request '/authors', got: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got: %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":200,"message":"OK","data":[{...}]}`))
	}))
	defer server.Close()

	client := server.Client()

	request, _ := http.NewRequest("GET", server.URL+"/authors", nil)
	response, err := client.Do(request)

	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("invalid response status code: got %d, want 200", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("invalid response status code: got %d, want 200", response.StatusCode)
	}

	if response.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("invalid content type: got %s, want application/json", response.Header.Get("Content-Type"))
	}

	expectedBody := `{"status":200,"message":"OK","data":[{...}]}`
	if !bytes.Equal(body, []byte(expectedBody)) {
		t.Fatalf("invalid response body: got %s, want %s", body, expectedBody)
	}
}

func TestGetAuthorsWithMock(t *testing.T) {
	jsonBytes := []byte(`{"status":200,"message":"OK","data":[{...}]}`)

	client := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Status:     "200 OK",
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(jsonBytes)),
			}, nil
		},
	}

	request, _ := http.NewRequest("GET", "/authors", nil)

	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("invalid response status code: got %d, want 200", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	defer response.Body.Close()

	expectedBody := jsonBytes
	if !bytes.Equal(body, []byte(expectedBody)) {
		t.Fatalf("invalid response body: got %s, want %s", body, expectedBody)
	}
}

func TestGetAuthorById(t *testing.T) {
	jsonBytes := []byte(`{"status":200,"message":"OK","data":{"ID":1,...}}`)

	client := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Status:     "200 OK",
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(jsonBytes)),
			}, nil
		},
	}

	request, _ := http.NewRequest("GET", "/authors/1", nil)

	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("invalid response status code: got %d, want 200", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	defer response.Body.Close()

	if !bytes.Equal(body, jsonBytes) {
		t.Fatalf("invalid response body: got %s, want %s", body, jsonBytes)
	}
}
