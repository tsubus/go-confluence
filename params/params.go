package params

import "github.com/umats/go-confluence/models"

// MultiEntityResult is a generic list wrapper returned by many endpoints.
type MultiEntityResult[T any] struct {
	Results []T                      `json:"results"`
	Links   *models.MultiEntityLinks `json:"_links"`
}

// MultiEntityResultPageBulk wraps page list responses.
type MultiEntityResultPageBulk = MultiEntityResult[models.PageBulk]

// MultiEntityResultAttachmentBulk wraps attachment list responses.
type MultiEntityResultAttachmentBulk = MultiEntityResult[models.AttachmentBulk]

// MultiEntityResultLabel wraps label list responses.
type MultiEntityResultLabel = MultiEntityResult[models.Label]

// MultiEntityResultPageVersion wraps page version list responses.
type MultiEntityResultPageVersion = MultiEntityResult[models.PageVersion]

// MultiEntityResultAttachmentVersion wraps attachment version list responses.
type MultiEntityResultAttachmentVersion = MultiEntityResult[models.AttachmentVersion]

// MultiEntityResultAncestor wraps ancestor list responses.
type MultiEntityResultAncestor = MultiEntityResult[models.Ancestor]

// MultiEntityResultPageCommentModel wraps page comment list responses.
type MultiEntityResultPageCommentModel = MultiEntityResult[models.PageCommentModel]

// MultiEntityResultPageInlineCommentModel wraps inline comment list responses.
type MultiEntityResultPageInlineCommentModel = MultiEntityResult[models.PageInlineCommentModel]

// MultiEntityResultContentProperty wraps content property list responses.
type MultiEntityResultContentProperty = MultiEntityResult[models.ContentProperty]

// MultiEntityResultCustomContentBulk wraps custom content list responses.
type MultiEntityResultCustomContentBulk = MultiEntityResult[models.CustomContentBulk]

// MultiEntityResultBlogPostBulk wraps blog post list responses.
type MultiEntityResultBlogPostBulk = MultiEntityResult[models.BlogPostBulk]

// MultiEntityResultSpaceBulk wraps space list responses.
type MultiEntityResultSpaceBulk = MultiEntityResult[models.SpaceBulk]

// MultiEntityResultSpacePermissionAssignment wraps space permission list responses.
type MultiEntityResultSpacePermissionAssignment = MultiEntityResult[models.SpacePermissionAssignment]

// MultiEntityResultSpaceProperty wraps space property list responses.
type MultiEntityResultSpaceProperty = MultiEntityResult[models.SpaceProperty]

// MultiEntityResultLike wraps like list responses.
type MultiEntityResultLike = MultiEntityResult[models.Like]

// MultiEntityResultAttachmentCommentModel wraps attachment comment list responses.
type MultiEntityResultAttachmentCommentModel = MultiEntityResult[models.AttachmentCommentModel]

// AttachmentListParams defines query parameters for listing attachments.
type AttachmentListParams struct {
	Sort      *models.AttachmentSortOrder `json:"sort"`
	Cursor    *string                     `json:"cursor"`
	Status    []string                    `json:"status"`
	MediaType *string                     `json:"mediaType"`
	Filename  *string                     `json:"filename"`
	Limit     *int                        `json:"limit"`
}

// AttachmentGetParams defines query parameters for fetching an attachment.
type AttachmentGetParams struct {
	Version              *int  `json:"version"`
	IncludeLabels        *bool `json:"include-labels"`
	IncludeProperties    *bool `json:"include-properties"`
	IncludeOperations    *bool `json:"include-operations"`
	IncludeVersions      *bool `json:"include-versions"`
	IncludeVersion       *bool `json:"include-version"`
	IncludeCollaborators *bool `json:"include-collaborators"`
}

// AttachmentLabelsParams defines query parameters for attachment labels.
type AttachmentLabelsParams struct {
	Prefix *string                `json:"prefix"`
	Sort   *models.LabelSortOrder `json:"sort"`
	Cursor *string                `json:"cursor"`
	Limit  *int                   `json:"limit"`
}

// AttachmentVersionsParams defines query parameters for attachment versions.
type AttachmentVersionsParams struct {
	Cursor *string                  `json:"cursor"`
	Limit  *int                     `json:"limit"`
	Sort   *models.VersionSortOrder `json:"sort"`
}

// AttachmentCommentsParams defines query parameters for attachment comments.
type AttachmentCommentsParams struct {
	BodyFormat *models.PrimaryBodyRepresentation `json:"body-format"`
	Cursor     *string                           `json:"cursor"`
	Limit      *int                              `json:"limit"`
	Sort       *models.CommentSortOrder          `json:"sort"`
	Version    *int64                            `json:"version"`
}

// LabelListParams defines query parameters for listing labels.
type LabelListParams struct {
	LabelID []int64                `json:"label-id"`
	Prefix  []string               `json:"prefix"`
	SpaceID []int64                `json:"space-id"`
	Sort    *models.LabelSortOrder `json:"sort"`
	Cursor  *string                `json:"cursor"`
	Limit   *int                   `json:"limit"`
}

// LabelAttachmentsParams defines query parameters for label attachments.
type LabelAttachmentsParams struct {
	Sort   *models.AttachmentSortOrder `json:"sort"`
	Cursor *string                     `json:"cursor"`
	Limit  *int                        `json:"limit"`
}

// LabelPagesParams defines query parameters for label pages.
type LabelPagesParams struct {
	SpaceID    []int64                           `json:"space-id"`
	BodyFormat *models.PrimaryBodyRepresentation `json:"body-format"`
	Sort       *models.PageSortOrder             `json:"sort"`
	Cursor     *string                           `json:"cursor"`
	Limit      *int                              `json:"limit"`
}

// LabelBlogPostsParams defines query parameters for label blog posts.
type LabelBlogPostsParams struct {
	SpaceID    []int64                           `json:"space-id"`
	BodyFormat *models.PrimaryBodyRepresentation `json:"body-format"`
	Sort       *models.BlogPostSortOrder         `json:"sort"`
	Cursor     *string                           `json:"cursor"`
	Limit      *int                              `json:"limit"`
}

// PageListParams defines query parameters for listing pages.
type PageListParams struct {
	ID         []int64                           `json:"id"`
	SpaceID    []int64                           `json:"space-id"`
	Sort       *models.PageSortOrder             `json:"sort"`
	Status     []string                          `json:"status"`
	Title      *string                           `json:"title"`
	BodyFormat *models.PrimaryBodyRepresentation `json:"body-format"`
	Subtype    *string                           `json:"subtype"`
	Cursor     *string                           `json:"cursor"`
	Limit      *int                              `json:"limit"`
}

// PageGetParams defines query parameters for fetching a page.
type PageGetParams struct {
	BodyFormat *models.PrimaryBodyRepresentationSingle `json:"body-format"`
	GetDraft   *bool                                   `json:"get-draft"`
	Status     []string                                `json:"status"`
	Version    *int                                    `json:"version"`

	IncludeLabels                       *bool `json:"include-labels"`
	IncludeProperties                   *bool `json:"include-properties"`
	IncludeOperations                   *bool `json:"include-operations"`
	IncludeLikes                        *bool `json:"include-likes"`
	IncludeVersions                     *bool `json:"include-versions"`
	IncludeVersion                      *bool `json:"include-version"`
	IncludeFavoritedByCurrentUserStatus *bool `json:"include-favorited-by-current-user-status"`
	IncludeWebresources                 *bool `json:"include-webresources"`
	IncludeDirectChildren               *bool `json:"include-direct-children"`
	Draft                               *bool `json:"draft"`
}

// PageCreateParams defines query parameters for creating a page.
type PageCreateParams struct {
	Embedded  *bool `json:"embedded"`
	Private   *bool `json:"private"`
	RootLevel *bool `json:"root-level"`
}

// PageAttachmentsParams defines query parameters for page attachments.
type PageAttachmentsParams struct {
	Sort      *models.AttachmentSortOrder `json:"sort"`
	Cursor    *string                     `json:"cursor"`
	Status    []string                    `json:"status"`
	MediaType *string                     `json:"mediaType"`
	Filename  *string                     `json:"filename"`
	Limit     *int                        `json:"limit"`
}

// PageLabelsParams defines query parameters for page labels.
type PageLabelsParams struct {
	Prefix *string                `json:"prefix"`
	Sort   *models.LabelSortOrder `json:"sort"`
	Cursor *string                `json:"cursor"`
	Limit  *int                   `json:"limit"`
}

// PageVersionsParams defines query parameters for page versions.
type PageVersionsParams struct {
	BodyFormat *models.PrimaryBodyRepresentation `json:"body-format"`
	Cursor     *string                           `json:"cursor"`
	Limit      *int                              `json:"limit"`
	Sort       *models.VersionSortOrder          `json:"sort"`
}

// PageAncestorsParams defines query parameters for page ancestors.
type PageAncestorsParams struct {
	Limit *int `json:"limit"`
}

// PageDescendantsParams defines query parameters for page descendants.
type PageDescendantsParams struct {
	Depth *int `json:"depth"`
}

// PageChildrenParams defines query parameters for page children.
type PageChildrenParams struct {
	Limit  *int                       `json:"limit"`
	Sort   *models.ChildPageSortOrder `json:"sort"`
	Cursor *string                    `json:"cursor"`
}

// PageCommentsParams defines query parameters for page comments.
type PageCommentsParams struct {
	BodyFormat       *models.PrimaryBodyRepresentation `json:"body-format"`
	Cursor           *string                           `json:"cursor"`
	Limit            *int                              `json:"limit"`
	Sort             *models.CommentSortOrder          `json:"sort"`
	ResolutionStatus []string                          `json:"resolution-status"`
}

// ContentPropertiesParams defines query parameters for content properties.
type ContentPropertiesParams struct {
	Key    *string                          `json:"key"`
	Sort   *models.ContentPropertySortOrder `json:"sort"`
	Cursor *string                          `json:"cursor"`
	Limit  *int                             `json:"limit"`
}

// PageCustomContentParams defines query parameters for custom content in pages.
type PageCustomContentParams struct {
	Type       string                                  `json:"type"`
	Cursor     *string                                 `json:"cursor"`
	Limit      *int                                    `json:"limit"`
	BodyFormat *models.CustomContentBodyRepresentation `json:"body-format"`
}

// PageLikeUsersParams defines query parameters for page like users.
type PageLikeUsersParams struct {
	Cursor *string `json:"cursor"`
	Limit  *int    `json:"limit"`
}

// SpaceListParams defines query parameters for listing spaces.
type SpaceListParams struct {
	IDs               []int64                                    `json:"ids"`
	Keys              []string                                   `json:"keys"`
	Type              *models.SpaceType                          `json:"type"`
	Status            *models.SpaceStatus                        `json:"status"`
	Labels            []string                                   `json:"labels"`
	FavoritedBy       *string                                    `json:"favorited-by"`
	NotFavoritedBy    *string                                    `json:"not-favorited-by"`
	Sort              *models.SpaceSortOrder                     `json:"sort"`
	DescriptionFormat *models.SpaceDescriptionBodyRepresentation `json:"description-format"`
	IncludeIcon       *bool                                      `json:"include-icon"`
	Cursor            *string                                    `json:"cursor"`
	Limit             *int                                       `json:"limit"`
}

// SpaceGetParams defines query parameters for fetching a space.
type SpaceGetParams struct {
	DescriptionFormat  *models.SpaceDescriptionBodyRepresentation `json:"description-format"`
	IncludeIcon        *bool                                      `json:"include-icon"`
	IncludeOperations  *bool                                      `json:"include-operations"`
	IncludeProperties  *bool                                      `json:"include-properties"`
	IncludePermissions *bool                                      `json:"include-permissions"`
	IncludeLabels      *bool                                      `json:"include-labels"`
}

// SpacePagesParams defines query parameters for space pages.
type SpacePagesParams struct {
	Sort       *models.PageSortOrder             `json:"sort"`
	Status     []string                          `json:"status"`
	Title      *string                           `json:"title"`
	BodyFormat *models.PrimaryBodyRepresentation `json:"body-format"`
	Cursor     *string                           `json:"cursor"`
	Limit      *int                              `json:"limit"`
}

// SpaceBlogPostsParams defines query parameters for space blog posts.
type SpaceBlogPostsParams struct {
	Sort       *models.BlogPostSortOrder         `json:"sort"`
	Status     []string                          `json:"status"`
	Title      *string                           `json:"title"`
	BodyFormat *models.PrimaryBodyRepresentation `json:"body-format"`
	Cursor     *string                           `json:"cursor"`
	Limit      *int                              `json:"limit"`
}

// SpaceLabelsParams defines query parameters for space labels.
type SpaceLabelsParams struct {
	Prefix *string                `json:"prefix"`
	Sort   *models.LabelSortOrder `json:"sort"`
	Cursor *string                `json:"cursor"`
	Limit  *int                   `json:"limit"`
}

// SpaceCustomContentParams defines query parameters for space custom content.
type SpaceCustomContentParams struct {
	Type       string                                  `json:"type"`
	Cursor     *string                                 `json:"cursor"`
	Limit      *int                                    `json:"limit"`
	BodyFormat *models.CustomContentBodyRepresentation `json:"body-format"`
}

// SpacePermissionsParams defines query parameters for space permissions.
type SpacePermissionsParams struct {
	Cursor *string `json:"cursor"`
	Limit  *int    `json:"limit"`
}

// SpacePropertiesParams defines query parameters for space properties.
type SpacePropertiesParams struct {
	Key    *string                          `json:"key"`
	Sort   *models.ContentPropertySortOrder `json:"sort"`
	Cursor *string                          `json:"cursor"`
	Limit  *int                             `json:"limit"`
}

// SpaceClassificationLevelParams defines query parameters for space classification.
type SpaceClassificationLevelParams struct {
	Status *string `json:"status"`
}

// PageClassificationLevelParams defines query parameters for page classification.
type PageClassificationLevelParams struct {
	Status *string `json:"status"`
}
