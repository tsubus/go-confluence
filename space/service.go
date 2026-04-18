package space

import (
	"context"
	"fmt"
	"net/http"

	"github.com/umats/go-confluence/internal/transport"
	"github.com/umats/go-confluence/models"
	"github.com/umats/go-confluence/params"
)

// Service provides space REST v2 operations.
type Service struct {
	v2 *transport.V2Client
}

func NewService(client *transport.Client) *Service {
	return &Service{v2: transport.NewV2Client(client)}
}

func (s *Service) Get(ctx context.Context, id int64, queryParams *params.SpaceGetParams) (*models.SpaceSingle, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build space get query: %w", err)
	}
	var response models.SpaceSingle
	path := fmt.Sprintf("spaces/%d", id)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request space: %w", err)
	}
	return &response, nil
}

func (s *Service) List(
	ctx context.Context,
	queryParams *params.SpaceListParams,
) (*params.MultiEntityResultSpaceBulk, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build space list query: %w", err)
	}
	var response params.MultiEntityResultSpaceBulk
	err = s.v2.DoJSON(ctx, http.MethodGet, "spaces", query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request spaces: %w", err)
	}
	return &response, nil
}

func (s *Service) Create(ctx context.Context, request models.SpaceCreateRequest) (*models.SpaceBulk, error) {
	var response models.SpaceBulk
	err := s.v2.DoJSON(ctx, http.MethodPost, "spaces", nil, request, &response)
	if err != nil {
		return nil, fmt.Errorf("create space: %w", err)
	}
	return &response, nil
}

func (s *Service) GetPages(
	ctx context.Context,
	id int64,
	queryParams *params.SpacePagesParams,
) (*params.MultiEntityResultPageBulk, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build space pages query: %w", err)
	}
	var response params.MultiEntityResultPageBulk
	path := fmt.Sprintf("spaces/%d/pages", id)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request space pages: %w", err)
	}
	return &response, nil
}

func (s *Service) GetBlogPosts(
	ctx context.Context,
	id int64,
	queryParams *params.SpaceBlogPostsParams,
) (*params.MultiEntityResultBlogPostBulk, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build space blog posts query: %w", err)
	}
	var response params.MultiEntityResultBlogPostBulk
	path := fmt.Sprintf("spaces/%d/blogposts", id)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request space blog posts: %w", err)
	}
	return &response, nil
}

func (s *Service) GetLabels(
	ctx context.Context,
	id int64,
	queryParams *params.SpaceLabelsParams,
) (*params.MultiEntityResultLabel, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build space labels query: %w", err)
	}
	var response params.MultiEntityResultLabel
	path := fmt.Sprintf("spaces/%d/labels", id)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request space labels: %w", err)
	}
	return &response, nil
}

func (s *Service) GetContentLabels(
	ctx context.Context,
	id int64,
	queryParams *params.SpaceLabelsParams,
) (*params.MultiEntityResultLabel, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build space content labels query: %w", err)
	}
	var response params.MultiEntityResultLabel
	path := fmt.Sprintf("spaces/%d/content/labels", id)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request space content labels: %w", err)
	}
	return &response, nil
}

func (s *Service) GetCustomContent(
	ctx context.Context,
	id int64,
	queryParams *params.SpaceCustomContentParams,
) (*params.MultiEntityResultCustomContentBulk, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build space custom content query: %w", err)
	}
	var response params.MultiEntityResultCustomContentBulk
	path := fmt.Sprintf("spaces/%d/custom-content", id)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request space custom content: %w", err)
	}
	return &response, nil
}

func (s *Service) GetOperations(ctx context.Context, id int64) (*models.PermittedOperationsResponse, error) {
	var response models.PermittedOperationsResponse
	path := fmt.Sprintf("spaces/%d/operations", id)
	err := s.v2.DoJSON(ctx, http.MethodGet, path, nil, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request space operations: %w", err)
	}
	return &response, nil
}

func (s *Service) GetPermissions(
	ctx context.Context,
	id int64,
	queryParams *params.SpacePermissionsParams,
) (*params.MultiEntityResultSpacePermissionAssignment, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build space permissions query: %w", err)
	}
	var response params.MultiEntityResultSpacePermissionAssignment
	path := fmt.Sprintf("spaces/%d/permissions", id)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request space permissions: %w", err)
	}
	return &response, nil
}

func (s *Service) GetProperties(
	ctx context.Context,
	id int64,
	queryParams *params.SpacePropertiesParams,
) (*params.MultiEntityResultSpaceProperty, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build space properties query: %w", err)
	}
	var response params.MultiEntityResultSpaceProperty
	path := fmt.Sprintf("spaces/%d/properties", id)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request space properties: %w", err)
	}
	return &response, nil
}

func (s *Service) GetPropertyByID(ctx context.Context, id int64, propertyID int64) (*models.SpaceProperty, error) {
	var response models.SpaceProperty
	path := fmt.Sprintf("spaces/%d/properties/%d", id, propertyID)
	err := s.v2.DoJSON(ctx, http.MethodGet, path, nil, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request space property: %w", err)
	}
	return &response, nil
}

func (s *Service) CreateProperty(
	ctx context.Context,
	id int64,
	request models.SpacePropertyCreateRequest,
) (*models.SpaceProperty, error) {
	var response models.SpaceProperty
	path := fmt.Sprintf("spaces/%d/properties", id)
	err := s.v2.DoJSON(ctx, http.MethodPost, path, nil, request, &response)
	if err != nil {
		return nil, fmt.Errorf("create space property: %w", err)
	}
	return &response, nil
}

func (s *Service) UpdateProperty(
	ctx context.Context,
	id int64,
	propertyID int64,
	request models.SpacePropertyUpdateRequest,
) (*models.SpaceProperty, error) {
	var response models.SpaceProperty
	path := fmt.Sprintf("spaces/%d/properties/%d", id, propertyID)
	err := s.v2.DoJSON(ctx, http.MethodPut, path, nil, request, &response)
	if err != nil {
		return nil, fmt.Errorf("update space property: %w", err)
	}
	return &response, nil
}

func (s *Service) DeleteProperty(ctx context.Context, id int64, propertyID int64) error {
	path := fmt.Sprintf("spaces/%d/properties/%d", id, propertyID)
	err := s.v2.DoJSON(ctx, http.MethodDelete, path, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("delete space property: %w", err)
	}
	return nil
}

func (s *Service) GetDefaultClassificationLevel(ctx context.Context, id int64) (*models.ClassificationLevel, error) {
	var response models.ClassificationLevel
	path := fmt.Sprintf("spaces/%d/classification-level/default", id)
	err := s.v2.DoJSON(ctx, http.MethodGet, path, nil, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request space default classification: %w", err)
	}
	return &response, nil
}

func (s *Service) UpdateDefaultClassificationLevel(
	ctx context.Context,
	id int64,
	request models.SpaceDefaultClassificationLevelUpdateRequest,
) (*models.ClassificationLevel, error) {
	var response models.ClassificationLevel
	path := fmt.Sprintf("spaces/%d/classification-level/default", id)
	err := s.v2.DoJSON(ctx, http.MethodPut, path, nil, request, &response)
	if err != nil {
		return nil, fmt.Errorf("update space default classification: %w", err)
	}
	return &response, nil
}

func (s *Service) DeleteDefaultClassificationLevel(
	ctx context.Context,
	id int64,
	request models.SpaceDefaultClassificationLevelUpdateRequest,
) (*models.ClassificationLevel, error) {
	var response models.ClassificationLevel
	path := fmt.Sprintf("spaces/%d/classification-level/default", id)
	err := s.v2.DoJSON(ctx, http.MethodDelete, path, nil, request, &response)
	if err != nil {
		return nil, fmt.Errorf("delete space default classification: %w", err)
	}
	return &response, nil
}
