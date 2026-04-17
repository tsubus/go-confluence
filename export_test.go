package confluence

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestDownloadFromRedirect(t *testing.T) {
	tests := []struct {
		name           string
		baseURL        string
		locationHeader string
		wantErr        error
		wantErrText    string
	}{
		{
			name:    "missing location",
			baseURL: "http://localhost:8080",
			wantErr: ErrMissingLocation,
		},
		{
			name:           "invalid base URL",
			baseURL:        "://bad-url",
			locationHeader: "/download/file.pdf",
			wantErrText:    "parse baseURL",
		},
		{
			name:           "invalid location header",
			baseURL:        "http://localhost:8080",
			locationHeader: "http://[::1",
			wantErrText:    "parse Location header",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{baseURL: tt.baseURL}
			resp := &http.Response{Header: http.Header{}}
			if tt.locationHeader != "" {
				resp.Header.Set("Location", tt.locationHeader)
			}

			err := client.downloadFromRedirect(context.Background(), resp, io.Discard)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}
			if tt.wantErrText != "" {
				if err == nil {
					t.Fatalf("expected error containing %q", tt.wantErrText)
				}
				if !strings.Contains(err.Error(), tt.wantErrText) {
					t.Fatalf("expected error containing %q, got %q", tt.wantErrText, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestDownloadPDF(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		body        io.ReadCloser
		writer      io.Writer
		wantErrText string
	}{
		{
			name:        "non-200 status",
			statusCode:  http.StatusInternalServerError,
			body:        io.NopCloser(strings.NewReader("boom")),
			writer:      io.Discard,
			wantErrText: "unexpected download status code 500",
		},
		{
			name:        "non-200 status body read error",
			statusCode:  http.StatusInternalServerError,
			body:        &errReadCloser{readErr: errors.New("read failed")},
			writer:      io.Discard,
			wantErrText: "read download error response",
		},
		{
			name:        "copy error",
			statusCode:  http.StatusOK,
			body:        &errReadCloser{readErr: errors.New("copy failed")},
			writer:      io.Discard,
			wantErrText: "stream pdf response body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{httpClient: &http.Client{Transport: roundTripperFunc(func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: tt.statusCode,
					Body:       tt.body,
					Header:     http.Header{},
				}, nil
			})}}

			err := client.downloadPDF(context.Background(), "http://example.com/file.pdf", tt.writer)
			if err == nil {
				t.Fatalf("expected error containing %q", tt.wantErrText)
			}
			if !strings.Contains(err.Error(), tt.wantErrText) {
				t.Fatalf("expected error containing %q, got %q", tt.wantErrText, err.Error())
			}
		})
	}
}

func TestHandleOKResponse_PDF(t *testing.T) {
	client := &Client{}

	resp := &http.Response{
		Header:     http.Header{"Content-Type": []string{"application/pdf"}},
		Body:       io.NopCloser(strings.NewReader("%PDF-1.4")),
		StatusCode: http.StatusOK,
	}

	var buffer bytes.Buffer
	err := client.handleOKResponse(context.Background(), resp, &buffer)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if buffer.String() != "%PDF-1.4" {
		t.Fatalf("unexpected pdf content: %q", buffer.String())
	}
}

func TestHandleOKResponse_HTMLMissingTaskID(t *testing.T) {
	client := &Client{}

	resp := &http.Response{
		Header:     http.Header{"Content-Type": []string{"text/html"}},
		Body:       io.NopCloser(strings.NewReader("<html>no task</html>")),
		StatusCode: http.StatusOK,
	}

	err := client.handleOKResponse(context.Background(), resp, io.Discard)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, ErrTaskIDNotFound) {
		t.Fatalf("expected ErrTaskIDNotFound, got %v", err)
	}
}

func TestFetchProgress(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		body        io.ReadCloser
		wantErrText string
	}{
		{
			name:        "non-200 status",
			statusCode:  http.StatusBadRequest,
			body:        io.NopCloser(strings.NewReader("bad request")),
			wantErrText: "unexpected poll status code 400",
		},
		{
			name:        "invalid json",
			statusCode:  http.StatusOK,
			body:        io.NopCloser(strings.NewReader("not-json")),
			wantErrText: "decode poll response",
		},
		{
			name:        "read error",
			statusCode:  http.StatusOK,
			body:        &errReadCloser{readErr: errors.New("read failed")},
			wantErrText: "read poll response body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{httpClient: &http.Client{Transport: roundTripperFunc(func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: tt.statusCode,
					Body:       tt.body,
					Header:     http.Header{},
				}, nil
			})}}

			_, err := client.fetchProgress(context.Background(), "http://example.com/progress")
			if err == nil {
				t.Fatalf("expected error containing %q", tt.wantErrText)
			}
			if !strings.Contains(err.Error(), tt.wantErrText) {
				t.Fatalf("expected error containing %q, got %q", tt.wantErrText, err.Error())
			}
		})
	}
}

func TestWaitForNextPoll_ContextCancelled(t *testing.T) {
	client := &Client{pollInterval: 10 * time.Millisecond}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := client.waitForNextPoll(ctx)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "context cancelled while polling") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPollTaskProgress_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	var once sync.Once
	client := &Client{
		baseURL:      "http://example.com",
		pollInterval: 5 * time.Millisecond,
		httpClient: &http.Client{Transport: roundTripperFunc(func(*http.Request) (*http.Response, error) {
			once.Do(cancel)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(
					`{"progress":50,"state":"IN_PROGRESS","result":""}`,
				)),
				Header: http.Header{},
			}, nil
		})},
	}

	_, err := client.pollTaskProgress(ctx, "task-1")
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "context cancelled while polling") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHandleOKResponse_TaskResultEmpty(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/api/v2/pdfexporttask/progress/task-1", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"progress":100,"state":"SUCCEEDED","result":""}`))
		if err != nil {
			return
		}
	})

	client := &Client{
		baseURL:      server.URL,
		httpClient:   server.Client(),
		pollInterval: 1 * time.Millisecond,
	}

	resp := &http.Response{
		Header:     http.Header{"Content-Type": []string{"text/html"}},
		Body:       io.NopCloser(strings.NewReader(`<meta name="ajs-taskId" content="task-1">`)),
		StatusCode: http.StatusOK,
	}

	err := client.handleOKResponse(context.Background(), resp, io.Discard)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, ErrTaskResultEmpty) {
		t.Fatalf("expected ErrTaskResultEmpty, got %v", err)
	}
}

func TestDownloadFromRedirect_Success(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/download/file.pdf", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/pdf")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("%PDF"))
		if err != nil {
			return
		}
	})

	client := &Client{baseURL: server.URL, httpClient: server.Client()}
	resp := &http.Response{Header: http.Header{"Location": []string{"/download/file.pdf"}}}

	var buffer bytes.Buffer
	err := client.downloadFromRedirect(context.Background(), resp, &buffer)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if buffer.String() != "%PDF" {
		t.Fatalf("unexpected pdf content: %q", buffer.String())
	}
}

func TestFetchProgress_Success(t *testing.T) {
	client := &Client{httpClient: &http.Client{Transport: roundTripperFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(strings.NewReader(
				`{"progress":100,"state":"SUCCEEDED","result":"/download"}`,
			)),
			Header: http.Header{},
		}, nil
	})}}

	pr, err := client.fetchProgress(context.Background(), "http://example.com/progress")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if pr.Progress != 100 {
		t.Fatalf("unexpected progress: %d", pr.Progress)
	}
	if pr.State != "SUCCEEDED" {
		t.Fatalf("unexpected state: %s", pr.State)
	}
	if pr.Result != "/download" {
		t.Fatalf("unexpected result: %s", pr.Result)
	}
}

func TestExportURL(t *testing.T) {
	client := &Client{baseURL: "http://localhost:8080/base"}
	url, err := client.exportURL("123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(url, "pageId=123") {
		t.Fatalf("expected pageId in query, got %q", url)
	}
	if !strings.Contains(url, "/spaces/flyingpdf/pdfpageexport.action") {
		t.Fatalf("unexpected export path: %q", url)
	}
}

func TestNewRequest_BasicAuth(t *testing.T) {
	client := &Client{username: "user", password: "pass"}
	req, err := client.newRequest(context.Background(), "http://example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	user, pass, ok := req.BasicAuth()
	if !ok {
		t.Fatalf("expected basic auth")
	}
	if user != "user" || pass != "pass" {
		t.Fatalf("unexpected basic auth: %s/%s", user, pass)
	}
}

func TestNewRequest_InvalidURL(t *testing.T) {
	client := &Client{}
	_, err := client.newRequest(context.Background(), "://bad-url")
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "create request") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHandleOKResponse_PDFCopyError(t *testing.T) {
	client := &Client{}

	resp := &http.Response{
		Header:     http.Header{"Content-Type": []string{"application/pdf"}},
		Body:       &errReadCloser{readErr: errors.New("copy failed")},
		StatusCode: http.StatusOK,
	}

	err := client.handleOKResponse(context.Background(), resp, io.Discard)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "stream pdf response body") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHandleOKResponse_ReadHTMLFailure(t *testing.T) {
	client := &Client{}

	resp := &http.Response{
		Header:     http.Header{"Content-Type": []string{"text/html"}},
		Body:       &errReadCloser{readErr: errors.New("read failed")},
		StatusCode: http.StatusOK,
	}

	err := client.handleOKResponse(context.Background(), resp, io.Discard)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "read export response body") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHandleOKResponse_TaskFailed(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/api/v2/pdfexporttask/progress/task-fail", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"progress":100,"state":"FAILED","result":""}`))
		if err != nil {
			return
		}
	})

	client := &Client{
		baseURL:      server.URL,
		httpClient:   server.Client(),
		pollInterval: 1 * time.Millisecond,
	}

	resp := &http.Response{
		Header:     http.Header{"Content-Type": []string{"text/html"}},
		Body:       io.NopCloser(strings.NewReader(`<meta name="ajs-taskId" content="task-fail">`)),
		StatusCode: http.StatusOK,
	}

	err := client.handleOKResponse(context.Background(), resp, io.Discard)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, ErrTaskFailed) {
		t.Fatalf("expected ErrTaskFailed, got %v", err)
	}
}

func TestHandleOKResponse_UnexpectedState(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/api/v2/pdfexporttask/progress/task-unknown", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"progress":100,"state":"UNKNOWN","result":""}`))
		if err != nil {
			return
		}
	})

	client := &Client{
		baseURL:      server.URL,
		httpClient:   server.Client(),
		pollInterval: 1 * time.Millisecond,
	}

	resp := &http.Response{
		Header:     http.Header{"Content-Type": []string{"text/html"}},
		Body:       io.NopCloser(strings.NewReader(`<meta name="ajs-taskId" content="task-unknown">`)),
		StatusCode: http.StatusOK,
	}

	err := client.handleOKResponse(context.Background(), resp, io.Discard)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "unexpected state") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExtractTaskID(t *testing.T) {
	tests := []struct {
		name        string
		html        string
		wantTaskID  string
		wantErrText string
	}{
		{
			name:       "valid task id",
			html:       `<meta name="ajs-taskId" content="task-123">`,
			wantTaskID: "task-123",
		},
		{
			name:        "missing task id",
			html:        "<html>missing</html>",
			wantErrText: ErrTaskIDNotFound.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractTaskID(tt.html)
			if tt.wantErrText != "" {
				if err == nil {
					t.Fatalf("expected error")
				}
				if !strings.Contains(err.Error(), tt.wantErrText) {
					t.Fatalf("expected error containing %q, got %q", tt.wantErrText, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if got != tt.wantTaskID {
				t.Fatalf("expected task id %q, got %q", tt.wantTaskID, got)
			}
		})
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

type errReadCloser struct {
	readErr error
}

func (rc *errReadCloser) Read(_ []byte) (int, error) {
	return 0, rc.readErr
}

func (rc *errReadCloser) Close() error {
	return nil
}
