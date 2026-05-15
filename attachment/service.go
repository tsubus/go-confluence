package attachment

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"

	"github.com/umats/go-confluence/internal/transport"
	"github.com/umats/go-confluence/models"
	"github.com/umats/go-confluence/params"
)

// Service provides attachment REST v2 operations.
type Service struct {
	v2 *transport.V2Client
}

func NewService(client *transport.Client) *Service {
	return &Service{v2: transport.NewV2Client(client)}
}

func (s *Service) List(
	ctx context.Context,
	queryParams *params.AttachmentListParams,
) (*params.MultiEntityResultAttachmentBulk, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build attachment list query: %w", err)
	}
	var response params.MultiEntityResultAttachmentBulk
	err = s.v2.DoJSON(ctx, http.MethodGet, "attachments", query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request attachments: %w", err)
	}
	return &response, nil
}

func (s *Service) Get(
	ctx context.Context,
	id string,
	queryParams *params.AttachmentGetParams,
) (*models.AttachmentSingle, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build attachment get query: %w", err)
	}
	var response models.AttachmentSingle
	path := fmt.Sprintf("attachments/%s", id)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request attachment: %w", err)
	}
	return &response, nil
}

func (s *Service) Delete(ctx context.Context, id int64, purge bool) error {
	query := url.Values{}
	if purge {
		query.Set("purge", "true")
	}
	path := fmt.Sprintf("attachments/%d", id)
	err := s.v2.DoJSON(ctx, http.MethodDelete, path, query, nil, nil)
	if err != nil {
		return fmt.Errorf("delete attachment: %w", err)
	}
	return nil
}

func (s *Service) GetLabels(
	ctx context.Context,
	id string,
	queryParams *params.AttachmentLabelsParams,
) (*params.MultiEntityResultLabel, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build attachment labels query: %w", err)
	}
	var response params.MultiEntityResultLabel
	path := fmt.Sprintf("attachments/%s/labels", id)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request attachment labels: %w", err)
	}
	return &response, nil
}

func (s *Service) GetOperations(ctx context.Context, id string) (*models.PermittedOperationsResponse, error) {
	var response models.PermittedOperationsResponse
	path := fmt.Sprintf("attachments/%s/operations", id)
	err := s.v2.DoJSON(ctx, http.MethodGet, path, nil, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request attachment operations: %w", err)
	}
	return &response, nil
}

func (s *Service) GetVersions(
	ctx context.Context,
	id string,
	queryParams *params.AttachmentVersionsParams,
) (*params.MultiEntityResultAttachmentVersion, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build attachment versions query: %w", err)
	}
	var response params.MultiEntityResultAttachmentVersion
	path := fmt.Sprintf("attachments/%s/versions", id)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request attachment versions: %w", err)
	}
	return &response, nil
}

func (s *Service) GetVersionDetails(
	ctx context.Context,
	id string,
	versionNumber int64,
) (*models.DetailedVersion, error) {
	var response models.DetailedVersion
	path := fmt.Sprintf("attachments/%s/versions/%d", id, versionNumber)
	err := s.v2.DoJSON(ctx, http.MethodGet, path, nil, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request attachment version details: %w", err)
	}
	return &response, nil
}

func (s *Service) GetComments(
	ctx context.Context,
	id string,
	queryParams *params.AttachmentCommentsParams,
) (*params.MultiEntityResultAttachmentCommentModel, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build attachment comments query: %w", err)
	}
	var response params.MultiEntityResultAttachmentCommentModel
	path := fmt.Sprintf("attachments/%s/footer-comments", id)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request attachment comments: %w", err)
	}
	return &response, nil
}

func (s *Service) GetContentProperties(
	ctx context.Context,
	id string,
	queryParams *params.ContentPropertiesParams,
) (*params.MultiEntityResultContentProperty, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build attachment properties query: %w", err)
	}
	var response params.MultiEntityResultContentProperty
	path := fmt.Sprintf("attachments/%s/properties", id)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request attachment properties: %w", err)
	}
	return &response, nil
}

func (s *Service) GetContentPropertyByID(
	ctx context.Context,
	id string,
	propertyID int64,
) (*models.ContentProperty, error) {
	var response models.ContentProperty
	path := fmt.Sprintf("attachments/%s/properties/%d", id, propertyID)
	err := s.v2.DoJSON(ctx, http.MethodGet, path, nil, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request attachment property: %w", err)
	}
	return &response, nil
}

func (s *Service) CreateContentProperty(
	ctx context.Context,
	id string,
	request models.ContentPropertyCreateRequest,
) (*models.ContentProperty, error) {
	var response models.ContentProperty
	path := fmt.Sprintf("attachments/%s/properties", id)
	err := s.v2.DoJSON(ctx, http.MethodPost, path, nil, request, &response)
	if err != nil {
		return nil, fmt.Errorf("create attachment property: %w", err)
	}
	return &response, nil
}

func (s *Service) UpdateContentProperty(
	ctx context.Context,
	id string,
	propertyID int64,
	request models.ContentPropertyUpdateRequest,
) (*models.ContentProperty, error) {
	var response models.ContentProperty
	path := fmt.Sprintf("attachments/%s/properties/%d", id, propertyID)
	err := s.v2.DoJSON(ctx, http.MethodPut, path, nil, request, &response)
	if err != nil {
		return nil, fmt.Errorf("update attachment property: %w", err)
	}
	return &response, nil
}

func (s *Service) DeleteContentProperty(
	ctx context.Context,
	id string,
	propertyID int64,
) error {
	path := fmt.Sprintf("attachments/%s/properties/%d", id, propertyID)
	err := s.v2.DoJSON(ctx, http.MethodDelete, path, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("delete attachment property: %w", err)
	}
	return nil
}

// Download fetches attachment metadata and streams the file through
// Confluence's supported attachment download endpoint.
func (s *Service) Download(
	ctx context.Context,
	id string,
	queryParams *params.AttachmentGetParams,
	writer io.Writer,
) error {
	attachment, err := s.Get(ctx, id, queryParams)
	if err != nil {
		return fmt.Errorf("fetch attachment metadata: %w", err)
	}
	containerID := attachmentContainerID(attachment)
	if containerID == "" {
		return errors.New("attachment has no container id")
	}
	counting := countingWriter{writer: writer}
	err = s.downloadByAPI(ctx, containerID, attachment.Id, &counting)
	if err != nil {
		return err
	}
	return nil
}

// downloadByAPI streams an attachment through Confluence's supported download endpoint.
func (s *Service) downloadByAPI(
	ctx context.Context,
	containerID string,
	attachmentID *string,
	writer io.Writer,
) error {
	if attachmentID == nil || *attachmentID == "" {
		return errors.New("attachment has no id")
	}
	downloadURL, err := attachmentDownloadAPIURL(s.v2.Client().BaseURL, containerID, *attachmentID)
	if err != nil {
		return err
	}
	return s.downloadURL(ctx, downloadURL, writer, false)
}

type countingWriter struct {
	writer  io.Writer
	written int64
}

func (w *countingWriter) Write(p []byte) (int, error) {
	n, err := w.writer.Write(p)
	w.written += int64(n)
	if err != nil {
		return n, fmt.Errorf("write attachment bytes: %w", err)
	}
	return n, nil
}

// Deprecated: DownloadByURL streams an attachment from a direct download URL to writer.
// Prefer Download for Confluence attachment IDs.
func (s *Service) DownloadByURL(
	ctx context.Context,
	downloadURL string,
	writer io.Writer,
) error {
	return s.downloadURL(ctx, downloadURL, writer, true)
}

func (s *Service) downloadURL(
	ctx context.Context,
	downloadURL string,
	writer io.Writer,
	warnDeprecated bool,
) error {
	if warnDeprecated {
		warnDeprecatedDownloadByURL()
	}
	client := s.v2.Client()
	if client == nil {
		return errors.New("transport client is required")
	}
	parsedURL, err := url.Parse(downloadURL)
	if err != nil {
		return fmt.Errorf("parse download URL %q: %w", downloadURL, err)
	}
	if parsedURL.Host == "" {
		return fmt.Errorf("download URL %q must include a host", downloadURL)
	}
	if parsedURL.Host != "api.media.atlassian.com" {
		err = s.ensureRedirectHostAllowed(parsedURL.Host)
		if err != nil {
			return fmt.Errorf("download URL host validation failed for %q: %w", downloadURL, err)
		}
	}
	req, err := client.NewRequest(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return fmt.Errorf("create download request for %q: %w", downloadURL, err)
	}
	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute download request for %q: %w", downloadURL, err)
	}
	closeBody := func() error {
		if closeErr := resp.Body.Close(); closeErr != nil {
			return fmt.Errorf("close download response body for %q: %w", downloadURL, closeErr)
		}
		return nil
	}
	switch resp.StatusCode {
	case http.StatusFound:
		err = s.downloadFromRedirect(ctx, resp, writer)
		if closeErr := closeBody(); err == nil && closeErr != nil {
			return closeErr
		}
		return err
	case http.StatusOK:
		_, err = io.Copy(writer, resp.Body)
		if closeErr := closeBody(); err == nil && closeErr != nil {
			return closeErr
		}
		if err != nil {
			return fmt.Errorf("stream download response from %q: %w", downloadURL, err)
		}
		return nil
	default:
		body, readErr := io.ReadAll(resp.Body)
		if closeErr := closeBody(); readErr == nil && closeErr != nil {
			return closeErr
		}
		if readErr != nil {
			return fmt.Errorf("read download error response from %q: %w", downloadURL, readErr)
		}
		return fmt.Errorf("unexpected download status code %d for %q: %s", resp.StatusCode, downloadURL, string(body))
	}
}

func attachmentContainerID(attachment *models.AttachmentSingle) string {
	if attachment == nil {
		return ""
	}
	if attachment.PageId != nil && *attachment.PageId != "" {
		return *attachment.PageId
	}
	if attachment.BlogPostId != nil && *attachment.BlogPostId != "" {
		return *attachment.BlogPostId
	}
	if attachment.CustomContentId != nil && *attachment.CustomContentId != "" {
		return *attachment.CustomContentId
	}
	return ""
}

func attachmentDownloadAPIURL(baseURL, pageID, attachmentID string) (string, error) {
	parsedBase, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("parse baseURL %q: %w", baseURL, err)
	}
	parsedBase.Path = strings.TrimSuffix(parsedBase.Path, "/")
	if !strings.HasSuffix(parsedBase.Path, "/wiki") {
		parsedBase.Path = path.Join(parsedBase.Path, "wiki")
	}
	parsedBase.Path = path.Join(
		parsedBase.Path,
		"rest",
		"api",
		"content",
		pageID,
		"child",
		"attachment",
		attachmentID,
		"download",
	)
	return parsedBase.String(), nil
}

func (s *Service) downloadFromRedirect(
	ctx context.Context,
	resp *http.Response,
	writer io.Writer,
) error {
	location := resp.Header.Get("Location")
	if location == "" {
		return fmt.Errorf("download response from %q missing Location header", responseURL(resp, s.v2.Client().BaseURL))
	}
	parsedBase, err := url.Parse(s.v2.Client().BaseURL)
	if err != nil {
		return fmt.Errorf("parse baseURL %q: %w", s.v2.Client().BaseURL, err)
	}
	downloadURL, err := parsedBase.Parse(location)
	if err != nil {
		return fmt.Errorf("parse Location header %q relative to %q: %w", location, parsedBase.String(), err)
	}
	resolved := downloadURL.String()
	if downloadURL.Host != parsedBase.Host && downloadURL.Host != "api.media.atlassian.com" {
		err = s.ensureRedirectHostAllowed(downloadURL.Host)
		if err != nil {
			return fmt.Errorf("redirect download URL host validation failed for %q: %w", resolved, err)
		}
	}
	err = s.downloadURL(ctx, resolved, writer, false)
	if err != nil {
		return fmt.Errorf("redirected download failed for %q: %w", resolved, err)
	}
	return nil
}

func (s *Service) ensureRedirectHostAllowed(host string) error {
	if host == "" {
		return errors.New("redirect host is empty")
	}
	client := s.v2.Client()
	if len(client.AllowedRedirectHosts) == 0 {
		return errors.New("allowed redirect host list is empty")
	}
	_, ok := client.AllowedRedirectHosts[host]
	if !ok {
		return fmt.Errorf("redirect host %q is not allowed", host)
	}
	return nil
}

func responseURL(resp *http.Response, fallback string) string {
	if resp != nil && resp.Request != nil && resp.Request.URL != nil {
		return resp.Request.URL.String()
	}
	return fallback
}

type deprecationWarning struct {
	logger *slog.Logger
	m      sync.Mutex
	once   func()
}

func newDeprecationWarning(logger *slog.Logger) *deprecationWarning {
	warning := &deprecationWarning{logger: logger}
	warning.resetLocked()

	return warning
}

func (w *deprecationWarning) Warn() {
	w.m.Lock()
	once := w.once
	w.m.Unlock()
	once()
}

func (w *deprecationWarning) ResetForTest() {
	w.m.Lock()
	defer w.m.Unlock()

	w.resetLocked()
}

func (w *deprecationWarning) resetLocked() {
	w.once = sync.OnceFunc(func() {
		w.logger.Warn(
			"confluence: AttachmentService.DownloadByURL is deprecated and will be removed in a future release",
			"replacement",
			"AttachmentService.Download",
		)
	})
}

func warnDeprecatedDownloadByURL() {
	defaultDownloadByURLDeprecationWarning().Warn()
}

func defaultDownloadByURLDeprecationWarning() *deprecationWarning {
	return defaultDownloadByURLDeprecationWarningState().warning
}

type downloadByURLDeprecationWarningState struct {
	warning *deprecationWarning
}

func defaultDownloadByURLDeprecationWarningState() *downloadByURLDeprecationWarningState {
	return &downloadByURLDeprecationWarningState{
		warning: newDeprecationWarning(slog.Default()),
	}
}

// ResetDeprecatedDownloadByURLWarningForTest resets the once-only deprecation warning for tests.
func ResetDeprecatedDownloadByURLWarningForTest() {
	defaultDownloadByURLDeprecationWarning().ResetForTest()
}

// NewDownloadByURLDeprecationWarningForTest exposes warning construction for tests.
func NewDownloadByURLDeprecationWarningForTest() interface {
	ResetForTest()
	Warn()
} {
	return newDeprecationWarning(slog.Default())
}
