package setup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
)

// TestRequest is a builder for making HTTP requests in tests.
// It simplifies sending requests with JSON bodies, auth cookies,
// and multipart form data.
type TestRequest struct {
	server *httptest.Server
	method string
	path   string
	body   io.Reader
	header http.Header
	cookie *http.Cookie
}

// JSON creates a new TestRequest with a JSON body.
func JSON(server *httptest.Server, method, path string, payload interface{}) *TestRequest {
	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			panic(fmt.Sprintf("failed to marshal JSON payload: %v", err))
		}
		body = bytes.NewReader(b)
	}

	tr := &TestRequest{
		server: server,
		method: method,
		path:   path,
		body:   body,
		header: make(http.Header),
	}
	tr.header.Set("Content-Type", "application/json")
	return tr
}

// GET creates a new GET TestRequest.
func GET(server *httptest.Server, path string) *TestRequest {
	return &TestRequest{
		server: server,
		method: http.MethodGet,
		path:   path,
		body:   nil,
		header: make(http.Header),
	}
}

// POST creates a new POST TestRequest.
func POST(server *httptest.Server, path string) *TestRequest {
	return &TestRequest{
		server: server,
		method: http.MethodPost,
		path:   path,
		body:   nil,
		header: make(http.Header),
	}
}

// PUT creates a new PUT TestRequest.
func PUT(server *httptest.Server, path string) *TestRequest {
	return &TestRequest{
		server: server,
		method: http.MethodPut,
		path:   path,
		body:   nil,
		header: make(http.Header),
	}
}

// DELETE creates a new DELETE TestRequest.
func DELETE(server *httptest.Server, path string) *TestRequest {
	return &TestRequest{
		server: server,
		method: http.MethodDelete,
		path:   path,
		body:   nil,
		header: make(http.Header),
	}
}

// WithAuth adds an authentication cookie to the request.
func (tr *TestRequest) WithAuth(sessionID string) *TestRequest {
	tr.cookie = AuthCookie(sessionID)
	return tr
}

// WithHeader adds a custom header to the request.
func (tr *TestRequest) WithHeader(key, value string) *TestRequest {
	tr.header.Set(key, value)
	return tr
}

// WithBody sets a raw body on the request.
func (tr *TestRequest) WithBody(body io.Reader) *TestRequest {
	tr.body = body
	return tr
}

// WithMultipart creates a multipart form data request.
// fields is a map of form field names to string values.
// files is a map of form field names to file content tuples (filename, content).
func (tr *TestRequest) WithMultipart(fields map[string]string, files map[string][]byte) *TestRequest {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	// Add text fields
	for key, value := range fields {
		if err := w.WriteField(key, value); err != nil {
			panic(fmt.Sprintf("failed to write field %s: %v", key, err))
		}
	}

	// Add files
	for fieldName, data := range files {
		part, err := w.CreateFormFile(fieldName, fieldName+"_test_file")
		if err != nil {
			panic(fmt.Sprintf("failed to create form file %s: %v", fieldName, err))
		}
		if _, err := part.Write(data); err != nil {
			panic(fmt.Sprintf("failed to write file data for %s: %v", fieldName, err))
		}
	}

	w.Close()
	tr.body = &buf
	tr.header.Set("Content-Type", w.FormDataContentType())
	return tr
}

// Do executes the request and returns the response.
func (tr *TestRequest) Do() (*http.Response, error) {
	url := tr.server.URL + tr.path

	var body io.Reader
	if tr.body != nil {
		// Read the body into a buffer so we can log it for debugging
		buf, err := io.ReadAll(tr.body)
		if err != nil {
			return nil, fmt.Errorf("reading request body: %w", err)
		}
		body = bytes.NewReader(buf)

		// Reset the original reader if it's seekable
		if seeker, ok := tr.body.(io.Seeker); ok {
			seeker.Seek(0, io.SeekStart)
		}
	}

	req, err := http.NewRequest(tr.method, url, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Apply headers
	for key, values := range tr.header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Apply auth cookie
	if tr.cookie != nil {
		req.AddCookie(tr.cookie)
	}

	return tr.server.Client().Do(req)
}

// DoOK executes the request and expects a successful response (2xx).
// Returns the response body as bytes.
func (tr *TestRequest) DoOK() ([]byte, error) {
	resp, err := tr.Do()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return body, fmt.Errorf("expected 2xx status, got %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// DoStatus executes the request and checks for a specific status code.
func (tr *TestRequest) DoStatus(expectedStatus int) ([]byte, error) {
	resp, err := tr.Do()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != expectedStatus {
		return body, fmt.Errorf("expected status %d, got %d: %s", expectedStatus, resp.StatusCode, string(body))
	}

	return body, nil
}

// ParseResponse parses the JSON response into the given target.
func ParseResponse(body []byte, target interface{}) error {
	return json.Unmarshal(body, target)
}

// ResponseWrapper matches the standard response structure from utilities.WriteJSON.
type ResponseWrapper struct {
	StatusCode int             `json:"status_code"`
	Message    string          `json:"message"`
	Data       json.RawMessage `json:"data,omitempty"`
}

// ParseResponseWrapper parses the standard API response wrapper.
func ParseResponseWrapper(body []byte) (*ResponseWrapper, error) {
	var wrapper ResponseWrapper
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper, nil
}

// GetData parses the data field from the response wrapper into the given target.
func (w *ResponseWrapper) GetData(target interface{}) error {
	if len(w.Data) == 0 {
		return nil
	}
	return json.Unmarshal(w.Data, target)
}

// Contains checks if the response body contains a given substring.
func Contains(body []byte, substr string) bool {
	return strings.Contains(string(body), substr)
}
