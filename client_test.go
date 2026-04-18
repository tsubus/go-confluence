package confluence_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	confluence "github.com/umats/go-confluence"
	"github.com/umats/go-confluence/export"
)

func TestExportPageTo(t *testing.T) {
	expectedPDF := []byte("%PDF-1.4 streamed pdf content")
	server := newTestExportServer(t, expectedPDF)
	defer server.Close()

	client, err := confluence.NewClient(server.URL)
	require.NoError(t, err)

	var buffer bytes.Buffer
	err = client.Export().PageTo(context.Background(), "12345", &buffer)
	require.NoError(t, err)
	require.Equal(t, expectedPDF, buffer.Bytes())
}

func TestExportPageTo_RequiresWriter(t *testing.T) {
	server := newTestExportServer(t, []byte("%PDF-1.4"))
	defer server.Close()

	client, err := confluence.NewClient(server.URL)
	require.NoError(t, err)

	err = client.Export().PageTo(context.Background(), "12345", nil)
	require.Error(t, err)
	require.ErrorContains(t, err, "writer is required")
}

func TestWithTimeout(t *testing.T) {
	_, err := confluence.NewClient("http://localhost:8090", confluence.WithTimeout(-1*time.Second))
	require.Error(t, err)
}

func TestWithPollTimeout(t *testing.T) {
	_, err := confluence.NewClient("http://localhost:8090", confluence.WithPollTimeout(-1*time.Second))
	require.Error(t, err)
}

func TestWithRequireHTTPS(t *testing.T) {
	_, err := confluence.NewClient("http://localhost:8090", confluence.WithRequireHTTPS())
	require.Error(t, err)
}

func TestExportPage_UsesSentinelErrors(t *testing.T) {
	expectedPDF := []byte("%PDF-1.4 dummy pdf content")
	server := newTestExportServer(t, expectedPDF)
	defer server.Close()

	client, err := confluence.NewClient(server.URL)
	require.NoError(t, err)

	_, err = client.Export().Page(context.Background(), "no-location")
	require.Error(t, err)
	require.ErrorIs(t, err, confluence.ErrMissingLocation)
}

func FuzzExtractTaskID(f *testing.F) {
	f.Add(`<html><head><meta name="ajs-taskId" content="task-1"></head></html>`)
	f.Add(`<html><head></head><body>missing</body></html>`)
	f.Add(`<meta name="ajs-taskId" content="task-2">`)

	f.Fuzz(func(_ *testing.T, input string) {
		_, _ = export.ExtractTaskIDForTest(input)
	})
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		opts    []confluence.Option
		wantErr bool
	}{
		{
			name:    "valid url with basic auth",
			baseURL: "http://localhost:8090",
			opts:    []confluence.Option{confluence.WithBasicAuth("admin", "admin")},
			wantErr: false,
		},
		{
			name:    "empty baseURL",
			baseURL: "",
			wantErr: true,
		},
		{
			name:    "invalid url",
			baseURL: "://missing-scheme",
			wantErr: true,
		},
		{
			name:    "url without host",
			baseURL: "/just/a/path",
			wantErr: true,
		},
		{
			name:    "nil http client option",
			baseURL: "http://localhost:8090",
			opts:    []confluence.Option{confluence.WithHTTPClient(nil)},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := confluence.NewClient(tt.baseURL, tt.opts...)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, client)
		})
	}
}

func TestExportPage(t *testing.T) {
	expectedPDF := []byte("%PDF-1.4 dummy pdf content")
	server := newTestExportServer(t, expectedPDF)
	defer server.Close()

	ctx := context.Background()

	tests := []struct {
		name    string
		pageID  string
		opts    []confluence.Option
		want    []byte
		wantErr bool
		errMsg  string
	}{
		{
			name:    "successful export",
			pageID:  "12345",
			want:    expectedPDF,
			wantErr: false,
		},
		{
			name:    "successful export with auth",
			pageID:  "auth-required",
			opts:    []confluence.Option{confluence.WithBasicAuth("admin", "secret")},
			want:    expectedPDF,
			wantErr: false,
		},
		{
			name:    "empty pageID",
			pageID:  "",
			wantErr: true,
			errMsg:  "pageID is required",
		},
		{
			name:    "server error on export",
			pageID:  "server-error",
			wantErr: true,
			errMsg:  "unexpected export status code 500",
		},
		{
			name:    "missing location header",
			pageID:  "no-location",
			wantErr: true,
			errMsg:  "export response missing Location header",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := confluence.NewClient(server.URL, tt.opts...)
			require.NoError(t, err)

			got, err := client.Export().Page(ctx, tt.pageID)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					require.ErrorContains(t, err, tt.errMsg)
				}
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func newTestExportServer(t *testing.T, expectedPDF []byte) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/spaces/flyingpdf/pdfpageexport.action", func(w http.ResponseWriter, r *http.Request) {
		pageID := r.URL.Query().Get("pageId")
		if pageID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("X-Accel-Redirect", "")
		w.Header().Set("Content-Type", "text/html")
		switch pageID {
		case "server-error":
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("server error"))
			return
		case "no-location":
			w.WriteHeader(http.StatusFound)
			return
		case "auth-required":
			user, pass, ok := r.BasicAuth()
			if !ok || user != "admin" || pass != "secret" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		w.Header().Set("Location", "/download/file.pdf")
		w.WriteHeader(http.StatusFound)
	})

	mux.HandleFunc("/download/file.pdf", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/pdf")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(expectedPDF)
		if err != nil {
			return
		}
	})

	mux.HandleFunc("/api/v2/pdfexporttask/progress/task-1", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"progress":100,"state":"SUCCEEDED","result":"/download/file.pdf"}`))
		if err != nil {
			return
		}
	})

	return server
}
