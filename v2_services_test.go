package confluence_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	confluence "github.com/umats/go-confluence"
)

type testQueryParams struct {
	Name  *string  `json:"name"`
	Tags  []string `json:"tags"`
	Limit *int     `json:"limit"`
}

type recordedRequest struct {
	method      string
	path        string
	rawQuery    string
	authUser    string
	authPass    string
	contentType string
	body        []byte
	handlerErr  error
}

type serverCase struct {
	method      string
	path        string
	rawQuery    string
	statusCode  int
	response    string
	contentType string
}

func newV2TestServer(t *testing.T, expected serverCase, record *recordedRequest) *httptest.Server {
	t.Helper()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		record.method = r.Method
		record.path = r.URL.Path
		record.rawQuery = r.URL.RawQuery
		record.contentType = r.Header.Get("Content-Type")
		record.authUser, record.authPass, _ = r.BasicAuth()

		body, err := ioReadAll(r)
		if err != nil {
			record.handlerErr = err
		}
		record.body = body

		if expected.contentType != "" {
			w.Header().Set("Content-Type", expected.contentType)
		}
		w.WriteHeader(expected.statusCode)
		_, writeErr := w.Write([]byte(expected.response))
		if writeErr != nil && record.handlerErr == nil {
			record.handlerErr = writeErr
		}
	})

	return httptest.NewServer(handler)
}

func ioReadAll(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

func TestPageService_Methods(t *testing.T) {
	tests := []struct {
		name     string
		execute  func(client *confluence.Client) error
		want     serverCase
		wantBody bool
	}{
		{
			name: "list pages",
			execute: func(client *confluence.Client) error {
				limit := 10
				_, err := client.Page().List(context.Background(), &confluence.PageListParams{Limit: &limit})
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/pages",
				rawQuery:    "limit=10",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"id":"42"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "create page",
			execute: func(client *confluence.Client) error {
				rootLevel := true
				req := confluence.PageCreateRequest{SpaceId: "SPACE"}
				_, err := client.Page().Create(context.Background(), &confluence.PageCreateParams{RootLevel: &rootLevel}, req)
				return err
			},
			want: serverCase{
				method:      http.MethodPost,
				path:        "/wiki/api/v2/pages",
				rawQuery:    "root-level=true",
				statusCode:  http.StatusOK,
				response:    `{"id":"99"}`,
				contentType: "application/json",
			},
			wantBody: true,
		},
		{
			name: "delete page",
			execute: func(client *confluence.Client) error {
				return client.Page().Delete(context.Background(), 77)
			},
			want: serverCase{
				method:     http.MethodDelete,
				path:       "/wiki/api/v2/pages/77",
				statusCode: http.StatusNoContent,
				response:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var record recordedRequest
			expected := tt.want
			server := newV2TestServer(t, expected, &record)
			defer server.Close()

			client, err := confluence.NewClient(server.URL, confluence.WithBasicAuth("user", "pass"))
			require.NoError(t, err)

			err = tt.execute(client)
			require.NoError(t, err)
			require.NoError(t, record.handlerErr)
			require.Equal(t, expected.method, record.method)
			require.Equal(t, expected.path, record.path)
			require.Equal(t, expected.rawQuery, record.rawQuery)
			require.Equal(t, "user", record.authUser)
			require.Equal(t, "pass", record.authPass)
			if tt.wantBody {
				require.NotEmpty(t, record.body)
				if expected.contentType != "" {
					require.Equal(t, expected.contentType, record.contentType)
				}
			}
		})
	}
}

func TestSpaceService_Methods(t *testing.T) {
	tests := []struct {
		name    string
		execute func(client *confluence.Client) error
		want    serverCase
	}{
		{
			name: "list spaces",
			execute: func(client *confluence.Client) error {
				limit := 5
				_, err := client.Space().List(context.Background(), &confluence.SpaceListParams{Limit: &limit})
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/spaces",
				rawQuery:    "limit=5",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"id":"1"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "create space",
			execute: func(client *confluence.Client) error {
				req := confluence.SpaceCreateRequest{Name: "Engineering"}
				_, err := client.Space().Create(context.Background(), req)
				return err
			},
			want: serverCase{
				method:      http.MethodPost,
				path:        "/wiki/api/v2/spaces",
				statusCode:  http.StatusOK,
				response:    `{"id":"abc"}`,
				contentType: "application/json",
			},
		},
		{
			name: "delete space default classification",
			execute: func(client *confluence.Client) error {
				req := confluence.SpaceDefaultClassificationLevelUpdateRequest{}
				_, err := client.Space().DeleteDefaultClassificationLevel(context.Background(), 12, req)
				return err
			},
			want: serverCase{
				method:      http.MethodDelete,
				path:        "/wiki/api/v2/spaces/12/classification-level/default",
				statusCode:  http.StatusOK,
				response:    `{"id":"x"}`,
				contentType: "application/json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var record recordedRequest
			expected := tt.want
			server := newV2TestServer(t, expected, &record)
			defer server.Close()

			client, err := confluence.NewClient(server.URL, confluence.WithBasicAuth("user", "pass"))
			require.NoError(t, err)

			err = tt.execute(client)
			require.NoError(t, err)
			require.NoError(t, record.handlerErr)
			require.Equal(t, expected.method, record.method)
			require.Equal(t, expected.path, record.path)
			require.Equal(t, expected.rawQuery, record.rawQuery)
		})
	}
}

func TestAttachmentLabelServices(t *testing.T) {
	tests := []struct {
		name    string
		execute func(client *confluence.Client) error
		want    serverCase
	}{
		{
			name: "attachment delete purge",
			execute: func(client *confluence.Client) error {
				return client.Attachment().Delete(context.Background(), 999, true)
			},
			want: serverCase{
				method:     http.MethodDelete,
				path:       "/wiki/api/v2/attachments/999",
				rawQuery:   "purge=true",
				statusCode: http.StatusNoContent,
				response:   "",
			},
		},
		{
			name: "label list",
			execute: func(client *confluence.Client) error {
				limit := 2
				_, err := client.Label().List(context.Background(), &confluence.LabelListParams{Limit: &limit})
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/labels",
				rawQuery:    "limit=2",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"id":"l1"}]}`,
				contentType: "application/json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var record recordedRequest
			expected := tt.want
			server := newV2TestServer(t, expected, &record)
			defer server.Close()

			client, err := confluence.NewClient(server.URL, confluence.WithBasicAuth("user", "pass"))
			require.NoError(t, err)

			err = tt.execute(client)
			require.NoError(t, err)
			require.NoError(t, record.handlerErr)
			require.Equal(t, expected.method, record.method)
			require.Equal(t, expected.path, record.path)
			require.Equal(t, expected.rawQuery, record.rawQuery)
		})
	}
}

func TestV2DecodeResponse_EmptyBody(t *testing.T) {
	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("")),
	}

	client, err := confluence.NewClient("http://example.com")
	require.NoError(t, err)

	apiErr, err := confluence.DecodeResponseForTest(client, response, &map[string]any{})
	require.NoError(t, err)
	require.Empty(t, apiErr)
}

func TestBuildQuery_OmitsNilFields(t *testing.T) {
	params := &testQueryParams{
		Name:  nil,
		Tags:  nil,
		Limit: nil,
	}
	values, err := confluence.BuildQueryForTest(params)
	require.NoError(t, err)
	require.Equal(t, url.Values{}, values)
}

func TestV2ClientDoJSON_EncodeError(t *testing.T) {
	client, err := confluence.NewClient("http://example.com")
	require.NoError(t, err)

	bad := make(chan int)
	err = confluence.DoJSONForTest(context.Background(), client, http.MethodPost, "pages", nil, bad, &map[string]any{})
	require.Error(t, err)
	require.ErrorContains(t, err, "encode request")
}

func TestV2ClientDoJSON_QueryEncodingError(t *testing.T) {
	params := map[string]any{"bad": json.RawMessage("{")}
	_, err := confluence.BuildQueryForTest(params)
	require.Error(t, err)
	require.ErrorContains(t, err, "struct, map, or url.Values")
}
