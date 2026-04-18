package attachment

import (
	"context"
	"fmt"
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
	id int64,
	queryParams *params.AttachmentLabelsParams,
) (*params.MultiEntityResultLabel, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build attachment labels query: %w", err)
	}
	var response params.MultiEntityResultLabel
	path := fmt.Sprintf("attachments/%d/labels", id)
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
