package attachment

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

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

// Download fetches attachment metadata and streams the file to writer.
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
	if attachment.DownloadLink == nil || *attachment.DownloadLink == "" {
		return errors.New("attachment has no download link")
	}
	return s.DownloadByURL(ctx, *attachment.DownloadLink, writer)
}

// DownloadByURL streams an attachment from a direct download URL to writer.
func (s *Service) DownloadByURL(
	ctx context.Context,
	downloadURL string,
	writer io.Writer,
) error {
	client := s.v2.Client()
	if client == nil {
		return errors.New("transport client is required")
	}
	req, err := client.NewRequest(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return fmt.Errorf("create download request: %w", err)
	}
	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute download request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	switch resp.StatusCode {
	case http.StatusFound:
		return s.downloadFromRedirect(ctx, resp, writer)
	case http.StatusOK:
		_, err = io.Copy(writer, resp.Body)
		if err != nil {
			return fmt.Errorf("stream download response: %w", err)
		}
		return nil
	default:
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return fmt.Errorf("read download error response: %w", readErr)
		}
		return fmt.Errorf("unexpected download status code %d: %s", resp.StatusCode, string(body))
	}
}

func (s *Service) downloadFromRedirect(
	ctx context.Context,
	resp *http.Response,
	writer io.Writer,
) error {
	location := resp.Header.Get("Location")
	if location == "" {
		return errors.New("download response missing Location header")
	}
	parsedBase, err := url.Parse(s.v2.Client().BaseURL)
	if err != nil {
		return fmt.Errorf("parse baseURL: %w", err)
	}
	downloadURL, err := parsedBase.Parse(location)
	if err != nil {
		return fmt.Errorf("parse Location header %q: %w", location, err)
	}
	err = s.ensureRedirectHostAllowed(downloadURL.Host)
	if err != nil {
		return err
	}
	return s.DownloadByURL(ctx, downloadURL.String(), writer)
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
