package client

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/stretchr/testify/assert"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/config"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

// mockRoundTripper struct
type mockRoundTripper struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

func TestRequestDo(t *testing.T) {
	// Create a mock HTTP client
	mockClient := &http.Client{}

	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write mock response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data": {"field": "value"}}`))
	}))
	defer mockServer.Close()

	// Create a mock request
	mockRequest, err := http.NewRequest("GET", mockServer.URL, nil)
	assert.NoError(t, err)

	// Create a mock data object
	expectedData := &struct {
		Field string `json:"field"`
	}{Field: "value"}

	// Create a mock slice for alldata
	var mockAllData []interface{}

	// Call the RequestDo function
	errorBool, diagnostics, hasPages, cursor := RequestDo(context.Background(), mockClient, mockRequest, nil, "resourceType", "operation", expectedData, &mockAllData)

	// Assertions
	assert.False(t, errorBool)
	assert.Empty(t, diagnostics)
	assert.False(t, hasPages)
	assert.Empty(t, cursor)
	assert.Equal(t, 1, len(mockAllData))
	assert.Equal(t, expectedData, mockAllData[0])
}
func TestExtractPageInfo(t *testing.T) {
	// NestedStruct struct
	type NestedStruct struct {
		PageInfo wiz.PageInfo `json:"pageInfo"`
	}
	// ParentStruct struct
	type ParentStruct struct {
		Foo NestedStruct `json:"foo"`
	}

	data := ParentStruct{
		Foo: NestedStruct{
			PageInfo: wiz.PageInfo{
				EndCursor:   "cursor123",
				HasNextPage: true,
			},
		},
	}

	// Call the function
	pageInfo, err := ExtractPageInfo(&data)

	// Check for errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check the extracted PageInfo values
	expected := wiz.PageInfo{
		EndCursor:   "cursor123",
		HasNextPage: true,
	}

	if !reflect.DeepEqual(pageInfo, expected) {
		t.Errorf("Unexpected PageInfo. Expected: %+v, but got: %+v", expected, pageInfo)
	}
}

func TestCreateRequest(t *testing.T) {
	ctx := context.TODO()

	// Create a sample configuration
	conf := &config.ProviderConf{

		Settings: &config.Settings{
			WizURL: "http://example.com", // This value won't be used as this test just validates the http.request object
		},
		UserAgent: "Test User Agent",
		TokenType: "Bearer",
		Token:     "testtoken",
	}

	// Create a sample request body
	body := bytes.NewBufferString("sample body")

	// Create sample diagnostics
	diags := diag.Diagnostics{}

	// Call the function
	request, isError, diagnostics := CreateRequest(ctx, conf, body, diags, "resourceType", "operation")

	// Check for errors
	if isError {
		t.Errorf("Unexpected error flag. Diagnostics: %+v", diagnostics)
	}

	// Check the request
	if request == nil {
		t.Errorf("Unexpected nil request")
	} else {
		// Check the request method
		expectedMethod := "POST"
		if request.Method != expectedMethod {
			t.Errorf("Unexpected request method. Expected: %s, but got: %s", expectedMethod, request.Method)
		}

		// Check the request URL
		expectedURL := "http://example.com"
		if request.URL.String() != expectedURL {
			t.Errorf("Unexpected request URL. Expected: %s, but got: %s", expectedURL, request.URL.String())
		}

		// Check the User-Agent header
		expectedUserAgent := "Test User Agent"
		if request.Header.Get("User-Agent") != expectedUserAgent {
			t.Errorf("Unexpected User-Agent header. Expected: %s, but got: %s", expectedUserAgent, request.Header.Get("User-Agent"))
		}

		// Check the Authorization header
		expectedAuthorization := "Bearer testtoken"
		if request.Header.Get("Authorization") != expectedAuthorization {
			t.Errorf("Unexpected Authorization header. Expected: %s, but got: %s", expectedAuthorization, request.Header.Get("Authorization"))
		}

		// Check the Content-Type header
		expectedContentType := "application/json"
		if request.Header.Get("Content-Type") != expectedContentType {
			t.Errorf("Unexpected Content-Type header. Expected: %s, but got: %s", expectedContentType, request.Header.Get("Content-Type"))
		}
	}
}

func TestProcessRequest(t *testing.T) {
	// Mock data
	mockVars := struct {
		ID   int
		Name string
	}{ID: 1, Name: "John"}

	mockData := struct {
		Field string `json:"field"`
	}{Field: "value"}

	mockQuery := "mock query"
	mockResourceType := "mock resource"
	mockOperation := "read"

	// Mock context
	ctx := context.TODO()

	// Create a custom RoundTripper with a RoundTrip implementation
	mockRoundTripper := &mockRoundTripper{
		// Implement the RoundTrip function to return a mock response
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			// Create a mock response with the desired data
			responseBody := []byte(`{"Data": {"field": "mock response"}}`)
			response := &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(responseBody)),
				Header:     make(http.Header),
			}
			return response, nil
		},
	}

	// Create a custom HTTP client using the mock RoundTripper
	mockClient := &http.Client{
		Transport: mockRoundTripper,
	}

	// Mock config
	mockProviderConf := &config.ProviderConf{
		HTTPClient: mockClient, // Use the mock HTTP client
		Settings: &config.Settings{
			WizURL: "http://example.com", // This value won't be used since we're using a mock client
		},
		UserAgent: "Test User Agent",
		TokenType: "Bearer",
		Token:     "testtoken",
	}

	// Call the function
	diags := ProcessRequest(ctx, mockProviderConf, mockVars, &mockData, mockQuery, mockResourceType, mockOperation)

	// Assertions
	assert.Empty(t, diags)
	// Add additional assertions as needed
}

func TestProcessPagedRequest(t *testing.T) {
	// Mock data
	mockVars := struct {
		ID   int
		Name string
	}{ID: 1, Name: "John"}

	mockData := struct {
		Field string `json:"field"`
	}{Field: "value"}

	mockQuery := "mock query"
	mockResourceType := "mock resource"
	mockOperation := "read"
	mockMaxPages := 3

	// Mock context
	ctx := context.TODO()

	// Create a custom RoundTripper with a RoundTrip implementation
	mockRoundTripper := &mockRoundTripper{
		// Implement the RoundTrip function to return a mock response
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			// Create a mock response with the desired data
			responseBody := []byte(`{"Data": {"field": "mock response"}}`)
			response := &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(responseBody)),
				Header:     make(http.Header),
			}
			return response, nil
		},
	}

	// Create a custom HTTP client using the mock RoundTripper
	mockClient := &http.Client{
		Transport: mockRoundTripper,
	}

	// Mock config
	mockProviderConf := &config.ProviderConf{
		HTTPClient: mockClient, // Use the mock HTTP client
		Settings: &config.Settings{
			WizURL: "http://example.com", // This value won't be used since we're using a mock client
		},
		UserAgent: "Test User Agent",
		TokenType: "Bearer",
		Token:     "testtoken",
	}

	// Call the function
	diags, allData := ProcessPagedRequest(ctx, mockProviderConf, mockVars, &mockData, mockQuery, mockResourceType, mockOperation, mockMaxPages)
	fmt.Printf("allData: %+v\n", allData)
	// Assertions
	assert.Empty(t, diags)
	assert.Equal(t, len(allData), 1)
}
