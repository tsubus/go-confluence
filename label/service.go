package label

import (
	"context"
	"fmt"
	"net/http"

	"github.com/umats/go-confluence/internal/transport"
	"github.com/umats/go-confluence/params"
)

// Service provides label REST v2 operations.
type Service struct {
	v2 *transport.V2Client
}

// NewService creates a new label service.
func NewService(client *transport.Client) *Service {
	return &Service{v2: transport.NewV2Client(client)}
}

// List returns labels matching the supplied query parameters.
func (s *Service) List(
	ctx context.Context,
	queryParams *params.LabelListParams,
) (*params.MultiEntityResultLabel, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build label list query: %w", err)
	}

	var response params.MultiEntityResultLabel
	err = s.v2.DoJSON(ctx, http.MethodGet, "labels", query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request labels: %w", err)
	}

	return &response, nil
}

// GetAttachments returns attachments associated with a label.
func (s *Service) GetAttachments(
	ctx context.Context,
	id int64,
	queryParams *params.LabelAttachmentsParams,
) (*params.MultiEntityResultAttachmentBulk, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build label attachments query: %w", err)
	}

	var response params.MultiEntityResultAttachmentBulk
	path := fmt.Sprintf("labels/%d/attachments", id)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request label attachments: %w", err)
	}

	return &response, nil
}

// GetPages returns pages associated with a label.
func (s *Service) GetPages(
	ctx context.Context,
	id int64,
	queryParams *params.LabelPagesParams,
) (*params.MultiEntityResultPageBulk, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build label pages query: %w", err)
	}

	var response params.MultiEntityResultPageBulk
	path := fmt.Sprintf("labels/%d/pages", id)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request label pages: %w", err)
	}

	return &response, nil
}

// GetBlogPosts returns blog posts associated with a label.
func (s *Service) GetBlogPosts(
	ctx context.Context,
	id int64,
	queryParams *params.LabelBlogPostsParams,
) (*params.MultiEntityResultBlogPostBulk, error) {
	query, err := transport.BuildQuery(queryParams)
	if err != nil {
		return nil, fmt.Errorf("build label blog posts query: %w", err)
	}

	var response params.MultiEntityResultBlogPostBulk
	path := fmt.Sprintf("labels/%d/blogposts", id)
	err = s.v2.DoJSON(ctx, http.MethodGet, path, query, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("request label blog posts: %w", err)
	}

	return &response, nil
}
