package v1

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	v1Mock "toko-buku-api/api/v1/mock"
	"toko-buku-api/internal/authors"
	"toko-buku-api/pkg/logger"

	"github.com/stretchr/testify/mock"
)

var authorLogTest = logger.NewService("AUTHOR_TEST")

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
	expected := `{"status":200,"message":"OK","data":[{"ID":1,"Updated_At":"0001-01-01T00:00:00Z","Country_Id":0,"Author":"test 1","City":""},{"ID":2,"Updated_At":"0001-01-01T00:00:00Z","Country_Id":0,"Author":"test 2","City":""}]}`

	request := httptest.NewRequest(http.MethodGet, "/authors", nil)
	writer := httptest.NewRecorder()

	authorMock := v1Mock.AuthorUsecaseMock{}
	authorMock.On("GetAuthors", mock.Anything).Return(&[]authors.Authors{
		{
			ID:     1,
			Author: "test 1",
		},
		{
			ID:     2,
			Author: "test 2",
		},
	}, nil)

	a := AuthorHandlerImpl{
		Usecase: &authorMock,
		Log:     authorLogTest,
	}

	a.GetAuthors(writer, request)

	res := writer.Result()
	defer res.Body.Close()

	// assertions
	if res.Status != "200 OK" {
		t.Errorf("Expected status OK but got %s", res.Status)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200 but got %d", res.StatusCode)
	}
	if res.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected content type application/json but got %s", res.Header.Get("Content-Type"))
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if string(data) != expected {
		t.Errorf("Expected %v but got %v", expected, string(data))
	}
}

// func TestGetAuthorByIdWithMock_fail(t *testing.T) {
// 	testCases := []struct {
// 		name            string
// 		body            string
// 		expectedStatus  int
// 		expectedMessage string
// 	}{
// 		{
// 			name:            "Not Found",
// 			body:            `{"status":404,"message":"Not Found"}`,
// 			expectedStatus:  http.StatusNotFound,
// 			expectedMessage: "Not Found",
// 		},
// 		{
// 			name:            "Internal Server Error",
// 			body:            `{"status":500,"message":"Internal Server Error"}`,
// 			expectedStatus:  http.StatusInternalServerError,
// 			expectedMessage: "Internal Server Error",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			jsonBytes := []byte(tc.body)
// 			client := &MockClient{
// 				DoFunc: func(req *http.Request) (*http.Response, error) {
// 					return &http.Response{
// 						Status:     "200 OK",
// 						Header:     http.Header{"Content-Type": []string{"application/json"}},
// 						StatusCode: http.StatusOK,
// 						Body:       io.NopCloser(bytes.NewReader(jsonBytes)),
// 					}, nil
// 				},
// 			}

// 			request, _ := http.NewRequest("GET", "/authors/3", nil)
// 			response, err := client.Do(request)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			defer response.Body.Close()
// 			if response.Header.Get("Content-Type") != "application/json" {
// 				t.Fatalf("invalid content type: got %s, want application/json", response.Header.Get("Content-Type"))
// 			}
// 			if response.StatusCode != http.StatusOK {
// 				t.Fatalf("invalid response status code: got %d, want 200", response.StatusCode)
// 			}

// 			body, _ := io.ReadAll(response.Body)
// 			if !bytes.Equal(body, jsonBytes) {
// 				t.Fatalf("invalid response body: got %s, want %s", body, jsonBytes)
// 			}

// 			var responseBody map[string]any
// 			json.Unmarshal(body, &responseBody)
// 			if int(responseBody["status"].(float64)) != tc.expectedStatus {
// 				t.Fatalf("invalid response status: got %d, want %d", responseBody["status"], tc.expectedStatus)
// 			}
// 			if responseBody["message"] != tc.expectedMessage {
// 				t.Fatalf("invalid response message: got %s, want %s", responseBody["message"], tc.expectedMessage)
// 			}
// 		})
// 	}
// }

// func TestGetAuthorByIdWithMock_success(t *testing.T) {
// 	jsonBytes := []byte(`{"status":200,"message":"OK","data":{"ID":1,...}}`)

// 	client := &MockClient{
// 		DoFunc: func(req *http.Request) (*http.Response, error) {
// 			return &http.Response{
// 				Status:     "200 OK",
// 				Header:     http.Header{"Content-Type": []string{"application/json"}},
// 				StatusCode: http.StatusOK,
// 				Body:       io.NopCloser(bytes.NewReader(jsonBytes)),
// 			}, nil
// 		},
// 	}

// 	request, _ := http.NewRequest("GET", "/authors/1", nil)

// 	response, err := client.Do(request)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer response.Body.Close()

// 	if response.StatusCode != http.StatusOK {
// 		t.Fatalf("invalid response status code: got %d, want 200", response.StatusCode)
// 	}

// 	body, err := io.ReadAll(response.Body)
// 	if err != nil {
// 		t.Fatalf("failed to read response body: %v", err)
// 	}
// 	defer response.Body.Close()

// 	if !bytes.Equal(body, jsonBytes) {
// 		t.Fatalf("invalid response body: got %s, want %s", body, jsonBytes)
// 	}
// }
