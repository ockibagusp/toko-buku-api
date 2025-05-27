package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
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
	testCases := []struct {
		name          string
		authorsStruct []authors.Authors
		expectedJSON  string
	}{
		{
			name:          "Authors empty",
			authorsStruct: []authors.Authors{},
			expectedJSON:  `{"status":200,"message":"OK","data":[]}`,
		},
		{
			name: "Authors is not empty",
			authorsStruct: []authors.Authors{
				{
					ID:     1,
					Author: "test 1",
				},
				{
					ID:     2,
					Author: "test 2",
				},
			},
			expectedJSON: `{"status":200,"message":"OK","data":[{"ID":1,"Updated_At":"0001-01-01T00:00:00Z","Country_Id":0,"Author":"test 1","City":""},{"ID":2,"Updated_At":"0001-01-01T00:00:00Z","Country_Id":0,"Author":"test 2","City":""}]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/authors", nil)
			writer := httptest.NewRecorder()

			authorMock := v1Mock.AuthorUsecaseMock{}
			authorMock.On("GetAuthors", mock.Anything).Return(&tc.authorsStruct, nil)

			authorHandlerImpl := AuthorHandlerImpl{
				Usecase: &authorMock,
				Log:     authorLogTest,
			}

			authorHandlerImpl.GetAuthors(writer, request)

			response := writer.Result()
			defer response.Body.Close()

			// assertions
			if response.Status != "200 OK" {
				t.Errorf("Expected status OK but got %s", response.Status)
			}
			if response.StatusCode != http.StatusOK {
				t.Errorf("Expected status code 200 but got %d", response.StatusCode)
			}
			if response.Header.Get("Content-Type") != "application/json" {
				t.Errorf("Expected content type application/json but got %s", response.Header.Get("Content-Type"))
			}

			data, err := io.ReadAll(response.Body)
			if err != nil {
				t.Errorf("Error: %v", err)
			}

			if string(data) != tc.expectedJSON {
				t.Errorf("Expected %v but got %v", tc.expectedJSON, string(data))
			}
		})
	}
}

func TestGetAuthorByIdWithErrMock_failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/authors/fail" {
			t.Errorf("Expected to request '/authors/fail', got: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got: %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":400,"message":"Bad Request"}`))
	}))
	defer server.Close()

	client := server.Client()

	request, _ := http.NewRequest("GET", server.URL+"/authors/fail", nil)
	response, err := client.Do(request)

	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer response.Body.Close()

	// assertions
	if response.Status != "400 Bad Request" {
		t.Errorf("Expected status OK but got %s", response.Status)
	}
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid response status code: got %d, want 400", response.StatusCode)
	}
	if response.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("invalid content type: got %s, want application/json", response.Header.Get("Content-Type"))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	expectedBody := `{"status":400,"message":"Bad Request"}`
	if !bytes.Equal(body, []byte(expectedBody)) {
		t.Fatalf("invalid response body: got %s, want %s", body, expectedBody)
	}

	var responseData map[string]any
	json.Unmarshal(body, &responseData)
	if int(responseData["status"].(float64)) != http.StatusBadRequest {
		t.Fatalf("invalid response status: got %d, want %d", responseData["status"], http.StatusBadRequest)
	}
	if responseData["message"] != "Bad Request" {
		t.Fatalf("invalid response message: got %s, want %s", responseData["message"], "Bad Request")
	}
}

func TestGetAuthorByIdMock_failure(t *testing.T) {
	testCases := []struct {
		name               string
		id                 string
		authorStruct       authors.Authors
		authorErr          error
		expectedStatus     string
		expectedStatusCode int
		expectedMessage    string
		expectedJSON       string
	}{
		// failed test cases
		{
			name:               "Not Found by ID: 3",
			id:                 "3",
			authorStruct:       authors.Authors{},
			authorErr:          errors.New("not found"),
			expectedStatus:     "404 Not Found",
			expectedStatusCode: http.StatusNotFound,
			expectedMessage:    "Not Found",
			expectedJSON:       `{"status":404,"message":"Not Found"}`,
		},
		{
			name:               "Not Found by ID: -3",
			id:                 "-3",
			authorStruct:       authors.Authors{},
			authorErr:          errors.New("not found"),
			expectedStatus:     "404 Not Found",
			expectedStatusCode: http.StatusNotFound,
			expectedMessage:    "Not Found",
			expectedJSON:       `{"status":404,"message":"Not Found"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/authors/3", nil)
			request.SetPathValue("authorById", tc.id)

			writer := httptest.NewRecorder()
			writer.WriteHeader(http.StatusNotFound)

			authorMock := v1Mock.AuthorUsecaseMock{}

			numString := request.PathValue("authorById")
			numInt, err := strconv.Atoi(numString)
			if err == nil {
				uint16Id := uint16(numInt)
				id := uint16(uint16Id)
				authorMock.On("GetAuthorById", mock.Anything, id).Return(&tc.authorStruct, tc.authorErr)
			}

			GetAuthorsForAssertions(t, &authorMock, request, writer, tc)
		})
	}
}

func GetAuthorsForAssertions(t *testing.T, authorMock *v1Mock.AuthorUsecaseMock, request *http.Request, writer *httptest.ResponseRecorder, tc struct {
	name               string
	id                 string
	authorStruct       authors.Authors
	authorErr          error
	expectedStatus     string
	expectedStatusCode int
	expectedMessage    string
	expectedJSON       string
}) {
	authorHandlerImpl := AuthorHandlerImpl{
		Usecase: authorMock,
		Log:     authorLogTest,
	}

	authorHandlerImpl.GetAuthorById(writer, request)

	response := writer.Result()
	defer response.Body.Close()

	// assertions
	if response.Status != tc.expectedStatus {
		t.Errorf("Expected %s but got %s", tc.expectedStatus, response.Status)
	}
	if response.StatusCode != tc.expectedStatusCode {
		t.Errorf("Expected status code %d but got %d", tc.expectedStatusCode, response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if string(body) != tc.expectedJSON {
		t.Errorf("Expected %v but got %v", tc.expectedJSON, string(body))
	}

	var responseData map[string]any
	json.Unmarshal(body, &responseData)
	if int(responseData["status"].(float64)) != tc.expectedStatusCode {
		t.Fatalf("invalid response status: got %d, want %d", responseData["status"], tc.expectedStatusCode)
	}
	if responseData["message"] != tc.expectedMessage {
		t.Fatalf("invalid response message: got %s, want %s", responseData["message"], tc.expectedMessage)
	}
}

func TestGetAuthorByIdWithMock_success(t *testing.T) {
	testCases := []struct {
		name               string
		id                 string
		authorStruct       authors.Authors
		authorErr          error
		expectedStatus     string
		expectedStatusCode int
		expectedMessage    string
		expectedJSON       string
	}{
		// success test cases
		{
			name: "Success by ID: 1",
			id:   "1",
			authorStruct: authors.Authors{
				ID:     1,
				Author: "test author 1",
			},
			authorErr:          nil,
			expectedStatus:     "200 OK",
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "OK",
			expectedJSON:       `{"status":200,"message":"OK","data":{"ID":1,"Updated_At":"0001-01-01T00:00:00Z","Country_Id":0,"Author":"test author 1","City":""}}`,
		},
		{
			name: "Success by ID: 2",
			id:   "2",
			authorStruct: authors.Authors{
				ID:     2,
				Author: "test author 2",
			},
			authorErr:          nil,
			expectedStatus:     "200 OK",
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "OK",
			expectedJSON:       `{"status":200,"message":"OK","data":{"ID":2,"Updated_At":"0001-01-01T00:00:00Z","Country_Id":0,"Author":"test author 2","City":""}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/authors/3", nil)
			request.SetPathValue("authorById", tc.id)

			writer := httptest.NewRecorder()
			writer.WriteHeader(http.StatusOK)

			authorMock := v1Mock.AuthorUsecaseMock{}

			numString := request.PathValue("authorById")
			numInt, _ := strconv.Atoi(numString)

			uint16Id := uint16(numInt)
			id := uint16(uint16Id)
			authorMock.On("GetAuthorById", mock.Anything, id).Return(&tc.authorStruct, nil)

			GetAuthorsForAssertions(t, &authorMock, request, writer, tc)
		})
	}
}
