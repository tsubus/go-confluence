package confluence

import (
	"context"
	"io"

	"github.com/umats/go-confluence/models"
	"github.com/umats/go-confluence/params"
)

// PageReader defines read-only page operations.
type PageReader interface {
	Get(
		ctx context.Context,
		id int64,
		params *params.PageGetParams,
	) (*models.PageSingle, error)
	List(
		ctx context.Context,
		params *params.PageListParams,
	) (*params.MultiEntityResultPageBulk, error)
	GetOperations(ctx context.Context, id int64) (*models.PermittedOperationsResponse, error)
	GetVersions(
		ctx context.Context,
		id int64,
		params *params.PageVersionsParams,
	) (*params.MultiEntityResultPageVersion, error)
	GetVersionDetails(
		ctx context.Context,
		id int64,
		versionNumber int64,
	) (*models.DetailedVersion, error)
	GetAncestors(
		ctx context.Context,
		id int64,
		params *params.PageAncestorsParams,
	) (*params.MultiEntityResultAncestor, error)
	GetDescendants(
		ctx context.Context,
		id int64,
		params *params.PageDescendantsParams,
	) (*models.DescendantsResponse, error)
	GetChildren(
		ctx context.Context,
		id int64,
		params *params.PageChildrenParams,
	) (*models.ChildrenResponse, error)
	GetDirectChildren(
		ctx context.Context,
		id int64,
		params *params.PageChildrenParams,
	) (*models.ChildrenResponse, error)
	GetLikeCount(ctx context.Context, id int64) (*models.Like, error)
	GetLikeUsers(
		ctx context.Context,
		id int64,
		params *params.PageLikeUsersParams,
	) (*params.MultiEntityResultLike, error)
	GetClassificationLevel(ctx context.Context, id int64) (*models.ClassificationLevel, error)
}

// PageWriter defines mutating page operations.
type PageWriter interface {
	Create(
		ctx context.Context,
		params *params.PageCreateParams,
		request models.PageCreateRequest,
	) (*models.PageSingle, error)
	Update(
		ctx context.Context,
		id int64,
		request models.PageUpdateRequest,
	) (*models.PageSingle, error)
	Delete(ctx context.Context, id int64) error
	UpdateTitle(
		ctx context.Context,
		id int64,
		request models.PageTitleUpdateRequest,
	) (*models.PageSingle, error)
	UpdateClassificationLevel(
		ctx context.Context,
		id int64,
		request models.ContentClassificationLevelUpdateRequest,
	) (*models.ClassificationLevel, error)
	ResetClassificationLevel(
		ctx context.Context,
		id int64,
		request models.ContentClassificationLevelDeleteRequest,
	) (*models.ClassificationLevel, error)
	Redact(
		ctx context.Context,
		id int64,
		request models.RedactionRequest,
	) (*models.RedactionResponse, error)
}

// PageAttachmentsReader defines attachment operations scoped to pages.
type PageAttachmentsReader interface {
	GetAttachments(
		ctx context.Context,
		id int64,
		params *params.PageAttachmentsParams,
	) (*params.MultiEntityResultAttachmentBulk, error)
}

// PageLabelsReader defines label operations scoped to pages.
type PageLabelsReader interface {
	GetLabels(
		ctx context.Context,
		id int64,
		params *params.PageLabelsParams,
	) (*params.MultiEntityResultLabel, error)
}

// PageCommentsReader defines comment operations scoped to pages.
type PageCommentsReader interface {
	GetFooterComments(
		ctx context.Context,
		id int64,
		params *params.PageCommentsParams,
	) (*params.MultiEntityResultPageCommentModel, error)
	GetInlineComments(
		ctx context.Context,
		id int64,
		params *params.PageCommentsParams,
	) (*params.MultiEntityResultPageInlineCommentModel, error)
}

// PagePropertiesService defines content property operations for pages.
type PagePropertiesService interface {
	GetContentProperties(
		ctx context.Context,
		id int64,
		params *params.ContentPropertiesParams,
	) (*params.MultiEntityResultContentProperty, error)
	GetContentPropertyByID(
		ctx context.Context,
		id int64,
		propertyID int64,
	) (*models.ContentProperty, error)
	CreateContentProperty(
		ctx context.Context,
		id int64,
		request models.ContentPropertyCreateRequest,
	) (*models.ContentProperty, error)
	UpdateContentProperty(
		ctx context.Context,
		id int64,
		propertyID int64,
		request models.ContentPropertyUpdateRequest,
	) (*models.ContentProperty, error)
	DeleteContentProperty(
		ctx context.Context,
		id int64,
		propertyID int64,
	) error
}

// PageCustomContentReader defines custom content operations for pages.
type PageCustomContentReader interface {
	GetCustomContent(
		ctx context.Context,
		id int64,
		params *params.PageCustomContentParams,
	) (*params.MultiEntityResultCustomContentBulk, error)
	GetCustomContentByURL(
		ctx context.Context,
		fullURL string,
	) (*params.MultiEntityResultCustomContentBulk, error)
}

// PageService defines page-related operations.
//
//go:generate go vet ./...
type PageService interface {
	PageReader
	PageWriter
	PageAttachmentsReader
	PageLabelsReader
	PageCommentsReader
	PagePropertiesService
	PageCustomContentReader
}

// SpaceReader defines read-only space operations.
type SpaceReader interface {
	Get(
		ctx context.Context,
		id int64,
		params *params.SpaceGetParams,
	) (*models.SpaceSingle, error)
	List(
		ctx context.Context,
		params *params.SpaceListParams,
	) (*params.MultiEntityResultSpaceBulk, error)
	GetPages(
		ctx context.Context,
		id int64,
		params *params.SpacePagesParams,
	) (*params.MultiEntityResultPageBulk, error)
	GetBlogPosts(
		ctx context.Context,
		id int64,
		params *params.SpaceBlogPostsParams,
	) (*params.MultiEntityResultBlogPostBulk, error)
	GetLabels(
		ctx context.Context,
		id int64,
		params *params.SpaceLabelsParams,
	) (*params.MultiEntityResultLabel, error)
	GetContentLabels(
		ctx context.Context,
		id int64,
		params *params.SpaceLabelsParams,
	) (*params.MultiEntityResultLabel, error)
	GetCustomContent(
		ctx context.Context,
		id int64,
		params *params.SpaceCustomContentParams,
	) (*params.MultiEntityResultCustomContentBulk, error)
	GetOperations(ctx context.Context, id int64) (*models.PermittedOperationsResponse, error)
	GetPermissions(
		ctx context.Context,
		id int64,
		params *params.SpacePermissionsParams,
	) (*params.MultiEntityResultSpacePermissionAssignment, error)
	GetProperties(
		ctx context.Context,
		id int64,
		params *params.SpacePropertiesParams,
	) (*params.MultiEntityResultSpaceProperty, error)
	GetPropertyByID(
		ctx context.Context,
		id int64,
		propertyID int64,
	) (*models.SpaceProperty, error)
	GetDefaultClassificationLevel(ctx context.Context, id int64) (*models.ClassificationLevel, error)
}

// SpaceWriter defines mutating space operations.
type SpaceWriter interface {
	Create(ctx context.Context, request models.SpaceCreateRequest) (*models.SpaceBulk, error)
	CreateProperty(
		ctx context.Context,
		id int64,
		request models.SpacePropertyCreateRequest,
	) (*models.SpaceProperty, error)
	UpdateProperty(
		ctx context.Context,
		id int64,
		propertyID int64,
		request models.SpacePropertyUpdateRequest,
	) (*models.SpaceProperty, error)
	DeleteProperty(ctx context.Context, id int64, propertyID int64) error
	UpdateDefaultClassificationLevel(
		ctx context.Context,
		id int64,
		request models.SpaceDefaultClassificationLevelUpdateRequest,
	) (*models.ClassificationLevel, error)
	DeleteDefaultClassificationLevel(
		ctx context.Context,
		id int64,
		request models.SpaceDefaultClassificationLevelUpdateRequest,
	) (*models.ClassificationLevel, error)
}

// SpaceService defines space-related operations.
type SpaceService interface {
	SpaceReader
	SpaceWriter
}

// AttachmentReader defines read-only attachment operations.
type AttachmentReader interface {
	List(
		ctx context.Context,
		params *params.AttachmentListParams,
	) (*params.MultiEntityResultAttachmentBulk, error)
	Get(
		ctx context.Context,
		id string,
		params *params.AttachmentGetParams,
	) (*models.AttachmentSingle, error)
	GetLabels(
		ctx context.Context,
		id int64,
		params *params.AttachmentLabelsParams,
	) (*params.MultiEntityResultLabel, error)
	GetOperations(ctx context.Context, id string) (*models.PermittedOperationsResponse, error)
	GetVersions(
		ctx context.Context,
		id string,
		params *params.AttachmentVersionsParams,
	) (*params.MultiEntityResultAttachmentVersion, error)
	GetVersionDetails(
		ctx context.Context,
		id string,
		versionNumber int64,
	) (*models.DetailedVersion, error)
	GetComments(
		ctx context.Context,
		id string,
		params *params.AttachmentCommentsParams,
	) (*params.MultiEntityResultAttachmentCommentModel, error)
	GetContentProperties(
		ctx context.Context,
		id string,
		params *params.ContentPropertiesParams,
	) (*params.MultiEntityResultContentProperty, error)
	GetContentPropertyByID(
		ctx context.Context,
		id string,
		propertyID int64,
	) (*models.ContentProperty, error)
}

// AttachmentWriter defines mutating attachment operations.
type AttachmentWriter interface {
	Delete(ctx context.Context, id int64, purge bool) error
	CreateContentProperty(
		ctx context.Context,
		id string,
		request models.ContentPropertyCreateRequest,
	) (*models.ContentProperty, error)
	UpdateContentProperty(
		ctx context.Context,
		id string,
		propertyID int64,
		request models.ContentPropertyUpdateRequest,
	) (*models.ContentProperty, error)
	DeleteContentProperty(
		ctx context.Context,
		id string,
		propertyID int64,
	) error
}

// AttachmentService defines attachment-related operations.
type AttachmentService interface {
	AttachmentReader
	AttachmentWriter
}

// LabelService defines label-related operations.
type LabelService interface {
	List(
		ctx context.Context,
		params *params.LabelListParams,
	) (*params.MultiEntityResultLabel, error)
	GetAttachments(
		ctx context.Context,
		id int64,
		params *params.LabelAttachmentsParams,
	) (*params.MultiEntityResultAttachmentBulk, error)
	GetPages(
		ctx context.Context,
		id int64,
		params *params.LabelPagesParams,
	) (*params.MultiEntityResultPageBulk, error)
	GetBlogPosts(
		ctx context.Context,
		id int64,
		params *params.LabelBlogPostsParams,
	) (*params.MultiEntityResultBlogPostBulk, error)
}

// ExportService defines PDF export operations.
type ExportService interface {
	Page(ctx context.Context, pageID string) ([]byte, error)
	PageTo(
		ctx context.Context,
		pageID string,
		writer io.Writer,
	) error
}
