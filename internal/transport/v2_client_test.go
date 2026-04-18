package transport_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/umats/go-confluence/internal/transport"
)

type testQueryParams struct {
	Name  *string  `json:"name"`
	Tags  []string `json:"tags"`
	Limit *int     `json:"limit"`
}

type testRequest struct {
	Title string `json:"title"`
}

type testResponse struct {
	ID string `json:"id"`
}

func TestBuildQuery(t *testing.T) {
	name := "alpha"
	limit := 10
	params := &testQueryParams{
		Name:  &name,
		Tags:  []string{"a", "b"},
		Limit: &limit,
	}
	valuesInput := url.Values{"name": {"alpha"}}
	stringMap := map[string]string{"title": "hello"}
	sliceMap := map[string][]string{"tags": {"a", "b"}}

	testCases := []struct {
		name      string
		input     any
		expected  url.Values
		expectErr bool
	}{
		{
			name:     "nil",
			input:    nil,
			expected: url.Values{},
		},
		{
			name:  "struct",
			input: params,
			expected: url.Values{
				"name":  {"alpha"},
				"tags":  {"a", "b"},
				"limit": {"10"},
			},
		},
		{
			name:     "url-values",
			input:    valuesInput,
			expected: url.Values{"name": {"alpha"}},
		},
		{
			name:     "url-values-pointer",
			input:    &valuesInput,
			expected: url.Values{"name": {"alpha"}},
		},
		{
			name:     "string-map",
			input:    stringMap,
			expected: url.Values{"title": {"hello"}},
		},
		{
			name:     "slice-map",
			input:    sliceMap,
			expected: url.Values{"tags": {"a", "b"}},
		},
		{
			name:      "unsupported",
			input:     123,
			expectErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			values, err := transport.BuildQuery(testCase.input)
			if testCase.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, testCase.expected, values)
		})
	}
}

func TestV2ClientDoJSON(t *testing.T) {
	var handlerErr error
	var capturedMethod string
	var capturedPath string
	var capturedQuery url.Values
	var capturedAuthUser string
	var capturedAuthPass string
	var capturedPayload testRequest
	var capturedContentType string

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method
		capturedPath = r.URL.Path
		capturedQuery = r.URL.Query()
		capturedContentType = r.Header.Get("Content-Type")

		user, pass, ok := r.BasicAuth()
		if !ok {
			handlerErr = errors.New("missing basic auth")
		}
		capturedAuthUser = user
		capturedAuthPass = pass

		body, readErr := io.ReadAll(r.Body)
		if readErr != nil {
			handlerErr = readErr
		}

		var payload testRequest
		unmarshalErr := json.Unmarshal(body, &payload)
		if unmarshalErr != nil {
			handlerErr = unmarshalErr
		}
		capturedPayload = payload

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, writeErr := w.Write([]byte(`{"id":"123"}`))
		if writeErr != nil {
			handlerErr = writeErr
		}
	}))
	defer testServer.Close()

	client := &transport.Client{
		BaseURL:    testServer.URL,
		HTTPClient: http.DefaultClient,
		Username:   "user",
		Password:   "pass",
	}
	v2 := transport.NewV2Client(client)

	query := url.Values{"title": {"hello"}}
	request := testRequest{Title: "My Page"}
	response := testResponse{}

	err := v2.DoJSON(context.Background(), http.MethodPost, "pages/123", query, request, &response)
	require.NoError(t, err)
	require.NoError(t, handlerErr)
	require.Equal(t, "123", response.ID)

	assert.Equal(t, http.MethodPost, capturedMethod)
	assert.Equal(t, "/wiki/api/v2/pages/123", capturedPath)
	assert.Equal(t, "hello", capturedQuery.Get("title"))
	assert.Equal(t, "user", capturedAuthUser)
	assert.Equal(t, "pass", capturedAuthPass)
	assert.Equal(t, "application/json", capturedContentType)
	assert.Equal(t, "My Page", capturedPayload.Title)
}

func TestV2ClientNewURL(t *testing.T) {
	client := &transport.Client{BaseURL: "http://example.com/base/"}
	v2 := transport.NewV2Client(client)

	endpoint, err := transport.NewURLForTest(v2, "pages")
	require.NoError(t, err)
	require.Equal(t, "http://example.com/base/wiki/api/v2/pages", endpoint)
}

func TestV2ClientNewURL_WithWikiSuffix(t *testing.T) {
	client := &transport.Client{BaseURL: "http://example.com/wiki"}
	v2 := transport.NewV2Client(client)

	endpoint, err := transport.NewURLForTest(v2, "pages")
	require.NoError(t, err)
	require.Equal(t, "http://example.com/wiki/api/v2/pages", endpoint)
}

func TestV2ClientDecodeResponseError(t *testing.T) {
	response := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(bytes.NewBufferString("bad request")),
	}

	client := &transport.Client{BaseURL: "http://example.com"}
	v2 := transport.NewV2Client(client)

	err := transport.DecodeResponseForTest(v2, response, nil)
	require.Error(t, err)

	var typedErr transport.APIError
	require.ErrorAs(t, err, &typedErr)
	require.Equal(t, http.StatusBadRequest, typedErr.StatusCode)
	require.Equal(t, "bad request", typedErr.Body)
}
