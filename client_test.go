package confluence_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	confluence "github.com/umats/go-confluence"
)

func TestExportPageTo(t *testing.T) {
	expectedPDF := []byte("%PDF-1.4 streamed pdf content")
	server := newTestExportServer(t, expectedPDF)
	defer server.Close()

	client, err := confluence.NewClient(server.URL)
	require.NoError(t, err)

	var buffer bytes.Buffer
	err = client.ExportPageTo(context.Background(), "12345", &buffer)
	require.NoError(t, err)
	require.Equal(t, expectedPDF, buffer.Bytes())
}

func TestExportPageTo_RequiresWriter(t *testing.T) {
	server := newTestExportServer(t, []byte("%PDF-1.4"))
	defer server.Close()

	client, err := confluence.NewClient(server.URL)
	require.NoError(t, err)

	err = client.ExportPageTo(context.Background(), "12345", nil)
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

	_, err = client.ExportPage(context.Background(), "no-location")
	require.Error(t, err)
	require.ErrorIs(t, err, confluence.ErrMissingLocation)
}

func FuzzExtractTaskID(f *testing.F) {
	f.Add(`<html><head><meta name="ajs-taskId" content="task-1"></head></html>`)
	f.Add(`<html><head></head><body>missing</body></html>`)
	f.Add(`<meta name="ajs-taskId" content="task-2">`)

	f.Fuzz(func(_ *testing.T, input string) {
		_, _ = confluence.ExtractTaskIDForTest(input)
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

			got, err := client.ExportPage(ctx, tt.pageID)
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
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		pageID := r.URL.Query().Get("pageId")
		if pageID == "" {
			http.Error(w, "missing pageId", http.StatusBadRequest)
			return
		}

		if pageID == "auth-required" {
			user, pass, ok := r.BasicAuth()
			if !ok || user != "admin" || pass != "secret" {
				w.Header().Set("WWW-Authenticate", `Basic realm="test"`)
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}

		if pageID == "no-location" {
			w.WriteHeader(http.StatusFound)
			return
		}

		if pageID == "server-error" {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Location", server.URL+"/download/export/pdfexport-20240101-000000-1/test.pdf")
		w.WriteHeader(http.StatusFound)
	})

	mux.HandleFunc("/download/export/pdfexport-20240101-000000-1/test.pdf", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(expectedPDF)
		if err != nil {
			return
		}
	})

	return server
}

func TestExportPage_DownloadFailure(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/spaces/flyingpdf/pdfpageexport.action", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Location", server.URL+"/download/export/missing.pdf")
		w.WriteHeader(http.StatusFound)
	})

	mux.HandleFunc("/download/export/missing.pdf", func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	})

	client, err := confluence.NewClient(server.URL)
	require.NoError(t, err)

	ctx := context.Background()
	_, err = client.ExportPage(ctx, "99999")
	require.Error(t, err)
	require.ErrorContains(t, err, "unexpected download status code 404")
}

func TestExportPage_ContextCancellation(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/spaces/flyingpdf/pdfpageexport.action", func(_ http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
		case <-time.After(5 * time.Second):
		}
	})

	client, err := confluence.NewClient(server.URL)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err = client.ExportPage(ctx, "12345")
	require.Error(t, err)
	require.ErrorContains(t, err, "context deadline exceeded")
}

func TestExportPage_CloudFlow(t *testing.T) {
	expectedPDF := []byte("%PDF-1.4 cloud pdf content")
	pollCount := 0

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/spaces/flyingpdf/pdfpageexport.action", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html;charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, htmlTaskResponse("task-abc-123"))
		if err != nil {
			return
		}
	})

	mux.HandleFunc("/api/v2/pdfexporttask/progress/task-abc-123", func(w http.ResponseWriter, _ *http.Request) {
		pollCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if pollCount < 2 {
			_, err := fmt.Fprint(w, progressJSON(50, "IN_PROGRESS", "", 3000, 1000))
			if err != nil {
				return
			}
			return
		}
		_, err := fmt.Fprint(w, progressJSON(100, "SUCCEEDED", server.URL+"/download/pdf", 0, 3000))
		if err != nil {
			return
		}
	})

	mux.HandleFunc("/download/pdf", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/pdf")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(expectedPDF)
		if err != nil {
			return
		}
	})

	client, err := confluence.NewClient(server.URL, confluence.WithPollInterval(100*time.Millisecond))
	require.NoError(t, err)

	ctx := context.Background()
	got, err := client.ExportPage(ctx, "cloud-page")
	require.NoError(t, err)
	require.Equal(t, expectedPDF, got)
	require.GreaterOrEqual(t, pollCount, 2)
}

func TestExportPage_CloudFlow_FailedTask(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/spaces/flyingpdf/pdfpageexport.action", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html;charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, htmlTaskResponse("task-fail-456"))
		if err != nil {
			return
		}
	})

	mux.HandleFunc("/api/v2/pdfexporttask/progress/task-fail-456", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, progressJSON(100, "FAILED", "", 0, 1000))
		if err != nil {
			return
		}
	})

	client, err := confluence.NewClient(server.URL, confluence.WithPollInterval(100*time.Millisecond))
	require.NoError(t, err)

	ctx := context.Background()
	_, err = client.ExportPage(ctx, "cloud-page-fail")
	require.Error(t, err)
	require.ErrorContains(t, err, "pdf export task failed")
}

func TestWithHTTPClient(t *testing.T) {
	custom := &http.Client{Timeout: 5 * time.Second}
	client, err := confluence.NewClient("http://localhost:8090", confluence.WithHTTPClient(custom))
	require.NoError(t, err)
	require.NotNil(t, client)
}

func ExampleClient_ExportPage() {
	// This example demonstrates the API usage against a mock server.
	expectedPDF := []byte("%PDF-1.4")

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/spaces/flyingpdf/pdfpageexport.action", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Location", server.URL+"/download/export/test.pdf")
		w.WriteHeader(http.StatusFound)
	})
	mux.HandleFunc("/download/export/test.pdf", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/pdf")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(expectedPDF)
		if err != nil {
			return
		}
	})

	client, _ := confluence.NewClient(server.URL)
	ctx := context.Background()
	pdf, err := client.ExportPage(ctx, "12345")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("downloaded %d bytes\n", len(pdf))
	// Output: downloaded 8 bytes
}

func htmlTaskResponse(taskID string) string {
	return fmt.Sprintf(
		`<html><head><meta name="ajs-taskId" content="%s"></head><body>exporting...</body></html>`,
		taskID,
	)
}

func progressJSON(progress int, state, result string, estimatedTimeRemaining, timeElapsed int) string {
	return fmt.Sprintf(
		`{"progress":%d,"state":"%s","result":"%s","estimatedTimeRemaining":%d,"timeElapsed":%d}`,
		progress,
		state,
		result,
		estimatedTimeRemaining,
		timeElapsed,
	)
}

func TestEvaluateProgress(t *testing.T) {
	tests := []struct {
		name        string
		progress    int
		state       string
		result      string
		wantDone    bool
		wantErrText string
	}{
		{
			name:     "in progress",
			progress: 50,
			state:    "IN_PROGRESS",
			wantDone: false,
		},
		{
			name:        "failed in progress",
			progress:    80,
			state:       "FAILED",
			wantDone:    true,
			wantErrText: "pdf export task failed",
		},
		{
			name:     "succeeded with result",
			progress: 100,
			state:    "SUCCEEDED",
			result:   "https://example.com/result.pdf",
			wantDone: true,
		},
		{
			name:        "succeeded without result",
			progress:    100,
			state:       "SUCCEEDED",
			wantDone:    true,
			wantErrText: "result URL is empty",
		},
		{
			name:        "unexpected terminal state",
			progress:    100,
			state:       "UNKNOWN",
			wantDone:    true,
			wantErrText: "unexpected state",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := confluence.ProgressResponseForTest(tt.progress, tt.state, tt.result)
			result, done, err := confluence.EvaluateProgressForTest(pr)
			require.Equal(t, tt.wantDone, done)
			if tt.wantErrText != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.wantErrText)
				return
			}
			require.NoError(t, err)
			if done {
				require.Equal(t, tt.result, result)
			}
		})
	}
}
