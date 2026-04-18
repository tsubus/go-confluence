package page

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/umats/go-confluence/internal/transport"
	"github.com/umats/go-confluence/models"
	"github.com/umats/go-confluence/params"
)

// Service provides page REST v2 operations.
type Service struct {
	v2 *transport.V2Client
}

func NewService(client *transport.Client) *Service {
	return &Service{v2: transport.NewV2Client(client)}
}

func (s *Service) Get(ctx context.Context, id int64, queryParams *params.PageGetParams) (*models.PageSingle, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build page get query: %w", err)
	}
	var response models.PageSingle
	path := fmt.Sprintf("pages/%d", id)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request page: %w", err)
	}
	return &response, nil
}

func (s *Service) List(
	ctx context.Context,
	queryParams *params.PageListParams,
) (*params.MultiEntityResultPageBulk, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build page list query: %w", err)
	}
	var response params.MultiEntityResultPageBulk
	err = s.v2.DoJSON(ctx, http.MethodGet, "pages", query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request pages: %w", err)
	}
	return &response, nil
}

func (s *Service) Create(
	ctx context.Context,
	queryParams *params.PageCreateParams,
	request models.PageCreateRequest,
) (*models.PageSingle, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build page create query: %w", err)
	}
	var response models.PageSingle
	err = s.v2.DoJSON(ctx, http.MethodPost, "pages", query, request, &response)
	if err != nil {
		return nil, fmt.Errorf("create page: %w", err)
	}
	return &response, nil
}

func (s *Service) Update(ctx context.Context, id int64, request models.PageUpdateRequest) (*models.PageSingle, error) {
	var response models.PageSingle
	path := fmt.Sprintf("pages/%d", id)
	err := s.v2.DoJSON(ctx, http.MethodPut, path, nil, request, &response)
	if err != nil {
		return nil, fmt.Errorf("update page: %w", err)
	}
	return &response, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	path := fmt.Sprintf("pages/%d", id)
	err := s.v2.DoJSON(ctx, http.MethodDelete, path, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("delete page: %w", err)
	}
	return nil
}

func (s *Service) UpdateTitle(
	ctx context.Context,
	id int64,
	request models.PageTitleUpdateRequest,
) (*models.PageSingle, error) {
	var response models.PageSingle
	path := fmt.Sprintf("pages/%d/title", id)
	err := s.v2.DoJSON(ctx, http.MethodPut, path, nil, request, &response)
	if err != nil {
		return nil, fmt.Errorf("update page title: %w", err)
	}
	return &response, nil
}

func (s *Service) GetAttachments(
	ctx context.Context,
	id int64,
	queryParams *params.PageAttachmentsParams,
) (*params.MultiEntityResultAttachmentBulk, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build page attachments query: %w", err)
	}
	var response params.MultiEntityResultAttachmentBulk
	path := fmt.Sprintf("pages/%d/attachments", id)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request page attachments: %w", err)
	}
	return &response, nil
}

func (s *Service) GetLabels(
	ctx context.Context,
	id int64,
	queryParams *params.PageLabelsParams,
) (*params.MultiEntityResultLabel, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build page labels query: %w", err)
	}
	var response params.MultiEntityResultLabel
	path := fmt.Sprintf("pages/%d/labels", id)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request page labels: %w", err)
	}
	return &response, nil
}

func (s *Service) GetOperations(ctx context.Context, id int64) (*models.PermittedOperationsResponse, error) {
	var response models.PermittedOperationsResponse
	path := fmt.Sprintf("pages/%d/operations", id)
	err := s.v2.DoJSON(ctx, http.MethodGet, path, nil, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request page operations: %w", err)
	}
	return &response, nil
}

func (s *Service) GetVersions(
	ctx context.Context,
	id int64,
	queryParams *params.PageVersionsParams,
) (*params.MultiEntityResultPageVersion, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build page versions query: %w", err)
	}
	var response params.MultiEntityResultPageVersion
	path := fmt.Sprintf("pages/%d/versions", id)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request page versions: %w", err)
	}
	return &response, nil
}

func (s *Service) GetVersionDetails(
	ctx context.Context,
	id int64,
	versionNumber int64,
) (*models.DetailedVersion, error) {
	var response models.DetailedVersion
	path := fmt.Sprintf("pages/%d/versions/%d", id, versionNumber)
	err := s.v2.DoJSON(ctx, http.MethodGet, path, nil, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request page version details: %w", err)
	}
	return &response, nil
}

func (s *Service) GetAncestors(
	ctx context.Context,
	id int64,
	queryParams *params.PageAncestorsParams,
) (*params.MultiEntityResultAncestor, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build page ancestors query: %w", err)
	}
	var response params.MultiEntityResultAncestor
	path := fmt.Sprintf("pages/%d/ancestors", id)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request page ancestors: %w", err)
	}
	return &response, nil
}

func (s *Service) GetDescendants(
	ctx context.Context,
	id int64,
	queryParams *params.PageDescendantsParams,
) (*models.DescendantsResponse, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build page descendants query: %w", err)
	}
	var response models.DescendantsResponse
	path := fmt.Sprintf("pages/%d/descendants", id)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request page descendants: %w", err)
	}
	return &response, nil
}

func (s *Service) GetChildren(
	ctx context.Context,
	id int64,
	queryParams *params.PageChildrenParams,
) (*models.ChildrenResponse, error) {
	return s.getChildren(ctx, id, "children", queryParams)
}

func (s *Service) GetDirectChildren(
	ctx context.Context,
	id int64,
	queryParams *params.PageChildrenParams,
) (*models.ChildrenResponse, error) {
	return s.getChildren(ctx, id, "direct-children", queryParams)
}

func (s *Service) getChildren(
	ctx context.Context,
	id int64,
	suffix string,
	queryParams *params.PageChildrenParams,
) (*models.ChildrenResponse, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build page children query: %w", err)
	}
	var response models.ChildrenResponse
	path := fmt.Sprintf("pages/%d/%s", id, suffix)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request page children: %w", err)
	}
	return &response, nil
}

func (s *Service) GetFooterComments(
	ctx context.Context,
	id int64,
	queryParams *params.PageCommentsParams,
) (*params.MultiEntityResultPageCommentModel, error) {
	return s.getComments(ctx, id, "footer-comments", queryParams)
}

func (s *Service) GetInlineComments(
	ctx context.Context,
	id int64,
	queryParams *params.PageCommentsParams,
) (*params.MultiEntityResultPageInlineCommentModel, error) {
	return s.getInlineComments(ctx, id, "inline-comments", queryParams)
}

func (s *Service) getComments(
	ctx context.Context,
	id int64,
	suffix string,
	queryParams *params.PageCommentsParams,
) (*params.MultiEntityResultPageCommentModel, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build page comments query: %w", err)
	}
	var response params.MultiEntityResultPageCommentModel
	path := fmt.Sprintf("pages/%d/%s", id, suffix)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request page comments: %w", err)
	}
	return &response, nil
}

func (s *Service) getInlineComments(
	ctx context.Context,
	id int64,
	suffix string,
	queryParams *params.PageCommentsParams,
) (*params.MultiEntityResultPageInlineCommentModel, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build page inline comments query: %w", err)
	}
	var response params.MultiEntityResultPageInlineCommentModel
	path := fmt.Sprintf("pages/%d/%s", id, suffix)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request page inline comments: %w", err)
	}
	return &response, nil
}

func (s *Service) GetContentProperties(
	ctx context.Context,
	id int64,
	queryParams *params.ContentPropertiesParams,
) (*params.MultiEntityResultContentProperty, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build page properties query: %w", err)
	}
	var response params.MultiEntityResultContentProperty
	path := fmt.Sprintf("pages/%d/properties", id)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request page properties: %w", err)
	}
	return &response, nil
}

func (s *Service) GetContentPropertyByID(
	ctx context.Context,
	id int64,
	propertyID int64,
) (*models.ContentProperty, error) {
	var response models.ContentProperty
	path := fmt.Sprintf("pages/%d/properties/%d", id, propertyID)
	err := s.v2.DoJSON(ctx, http.MethodGet, path, nil, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request page property: %w", err)
	}
	return &response, nil
}

func (s *Service) CreateContentProperty(
	ctx context.Context,
	id int64,
	request models.ContentPropertyCreateRequest,
) (*models.ContentProperty, error) {
	var response models.ContentProperty
	path := fmt.Sprintf("pages/%d/properties", id)
	err := s.v2.DoJSON(ctx, http.MethodPost, path, nil, request, &response)
	if err != nil {
		return nil, fmt.Errorf("create page property: %w", err)
	}
	return &response, nil
}

func (s *Service) UpdateContentProperty(
	ctx context.Context,
	id int64,
	propertyID int64,
	request models.ContentPropertyUpdateRequest,
) (*models.ContentProperty, error) {
	var response models.ContentProperty
	path := fmt.Sprintf("pages/%d/properties/%d", id, propertyID)
	err := s.v2.DoJSON(ctx, http.MethodPut, path, nil, request, &response)
	if err != nil {
		return nil, fmt.Errorf("update page property: %w", err)
	}
	return &response, nil
}

func (s *Service) DeleteContentProperty(
	ctx context.Context,
	id int64,
	propertyID int64,
) error {
	path := fmt.Sprintf("pages/%d/properties/%d", id, propertyID)
	err := s.v2.DoJSON(ctx, http.MethodDelete, path, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("delete page property: %w", err)
	}
	return nil
}

func (s *Service) GetCustomContent(
	ctx context.Context,
	id int64,
	queryParams *params.PageCustomContentParams,
) (*params.MultiEntityResultCustomContentBulk, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build page custom content query: %w", err)
	}
	var response params.MultiEntityResultCustomContentBulk
	path := fmt.Sprintf("pages/%d/custom-content", id)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request page custom content: %w", err)
	}
	return &response, nil
}

func (s *Service) GetLikeCount(ctx context.Context, id int64) (*models.Like, error) {
	var response models.Like
	path := fmt.Sprintf("pages/%d/likes/count", id)
	err := s.v2.DoJSON(ctx, http.MethodGet, path, nil, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request page like count: %w", err)
	}
	return &response, nil
}

func (s *Service) GetLikeUsers(
	ctx context.Context,
	id int64,
	queryParams *params.PageLikeUsersParams,
) (*params.MultiEntityResultLike, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build page like users query: %w", err)
	}
	var response params.MultiEntityResultLike
	path := fmt.Sprintf("pages/%d/likes/users", id)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request page like users: %w", err)
	}
	return &response, nil
}

func (s *Service) GetClassificationLevel(ctx context.Context, id int64) (*models.ClassificationLevel, error) {
	var response models.ClassificationLevel
	path := fmt.Sprintf("pages/%d/classification-level", id)
	err := s.v2.DoJSON(ctx, http.MethodGet, path, nil, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request page classification level: %w", err)
	}
	return &response, nil
}

func (s *Service) UpdateClassificationLevel(
	ctx context.Context,
	id int64,
	request models.ContentClassificationLevelUpdateRequest,
) (*models.ClassificationLevel, error) {
	var response models.ClassificationLevel
	path := fmt.Sprintf("pages/%d/classification-level", id)
	err := s.v2.DoJSON(ctx, http.MethodPut, path, nil, request, &response)
	if err != nil {
		return nil, fmt.Errorf("update page classification level: %w", err)
	}
	return &response, nil
}

func (s *Service) ResetClassificationLevel(
	ctx context.Context,
	id int64,
	request models.ContentClassificationLevelDeleteRequest,
) (*models.ClassificationLevel, error) {
	var response models.ClassificationLevel
	path := fmt.Sprintf("pages/%d/classification-level/reset", id)
	err := s.v2.DoJSON(ctx, http.MethodPost, path, nil, request, &response)
	if err != nil {
		return nil, fmt.Errorf("reset page classification level: %w", err)
	}
	return &response, nil
}

func (s *Service) Redact(
	ctx context.Context,
	id int64,
	request models.RedactionRequest,
) (*models.RedactionResponse, error) {
	var response models.RedactionResponse
	path := fmt.Sprintf("pages/%d/redact", id)
	err := s.v2.DoJSON(ctx, http.MethodPost, path, nil, request, &response)
	if err != nil {
		return nil, fmt.Errorf("redact page: %w", err)
	}
	return &response, nil
}

func (s *Service) GetCustomContentByURL(
	ctx context.Context,
	fullURL string,
) (*params.MultiEntityResultCustomContentBulk, error) {
	endpoint, err := url.Parse(fullURL)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}
	err = s.validateCustomContentURL(endpoint)
	if err != nil {
		return nil, fmt.Errorf("validate custom content url: %w", err)
	}
	query := endpoint.Query()
	path := strings.TrimPrefix(endpoint.Path, "/wiki/api/v2/")
	var response params.MultiEntityResultCustomContentBulk
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request custom content: %w", err)
	}
	return &response, nil
}

func (s *Service) validateCustomContentURL(endpoint *url.URL) error {
	if endpoint == nil {
		return errors.New("custom content URL is required")
	}
	if endpoint.Host == "" {
		return errors.New("custom content URL must include host")
	}
	client := s.v2.Client()
	if client == nil {
		return errors.New("transport client is required")
	}
	if !client.AllowCrossHostContentURL {
		base, err := url.Parse(client.BaseURL)
		if err != nil {
			return fmt.Errorf("parse baseURL: %w", err)
		}
		if !strings.EqualFold(base.Host, endpoint.Host) {
			return fmt.Errorf("custom content host %q does not match base host %q", endpoint.Host, base.Host)
		}
	}
	return nil
}

// Placeholder for future usage; keeps lints happy when unused.
var _ = url.Values{}
var _ = bytes.MinRead
