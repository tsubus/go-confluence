// Code generated from spec/openapi_subset.json; DO NOT EDIT.
package models

import "time"

// The type of ancestor.
type AncestorType string
const (
	AncestorTypePage AncestorType = "page"
	AncestorTypeWhiteboard AncestorType = "whiteboard"
	AncestorTypeDatabase AncestorType = "database"
	AncestorTypeEmbed AncestorType = "embed"
	AncestorTypeFolder AncestorType = "folder"
)

// The sort fields for attachments. The default sort direction is ascending. To sort in descending order, append a `-` character before the sort field. For example, `fieldName` or `-fieldName`.
type AttachmentSortOrder string
const (
	AttachmentSortOrderCreatedDate AttachmentSortOrder = "created-date"
	AttachmentSortOrderDescCreatedDate AttachmentSortOrder = "-created-date"
	AttachmentSortOrderModifiedDate AttachmentSortOrder = "modified-date"
	AttachmentSortOrderDescModifiedDate AttachmentSortOrder = "-modified-date"
)

// The status of the content.
type BlogPostContentStatus string
const (
	BlogPostContentStatusCurrent BlogPostContentStatus = "current"
	BlogPostContentStatusDraft BlogPostContentStatus = "draft"
	BlogPostContentStatusHistorical BlogPostContentStatus = "historical"
	BlogPostContentStatusTrashed BlogPostContentStatus = "trashed"
	BlogPostContentStatusDeleted BlogPostContentStatus = "deleted"
	BlogPostContentStatusAny BlogPostContentStatus = "any"
)

// The sort fields for blog posts. The default sort direction is ascending. To sort in descending order, append a `-` character before the sort field. For example, `fieldName` or `-fieldName`.
type BlogPostSortOrder string
const (
	BlogPostSortOrderId BlogPostSortOrder = "id"
	BlogPostSortOrderDescId BlogPostSortOrder = "-id"
	BlogPostSortOrderCreatedDate BlogPostSortOrder = "created-date"
	BlogPostSortOrderDescCreatedDate BlogPostSortOrder = "-created-date"
	BlogPostSortOrderModifiedDate BlogPostSortOrder = "modified-date"
	BlogPostSortOrderDescModifiedDate BlogPostSortOrder = "-modified-date"
)

// The sort fields for child pages. The default sort direction is ascending by child-position. To sort in descending order, append a `-` character before the sort field. For example, `fieldName` or `-fieldName`.
type ChildPageSortOrder string
const (
	ChildPageSortOrderCreatedDate ChildPageSortOrder = "created-date"
	ChildPageSortOrderDescCreatedDate ChildPageSortOrder = "-created-date"
	ChildPageSortOrderId ChildPageSortOrder = "id"
	ChildPageSortOrderDescId ChildPageSortOrder = "-id"
	ChildPageSortOrderChildPosition ChildPageSortOrder = "child-position"
	ChildPageSortOrderDescChildPosition ChildPageSortOrder = "-child-position"
	ChildPageSortOrderModifiedDate ChildPageSortOrder = "modified-date"
	ChildPageSortOrderDescModifiedDate ChildPageSortOrder = "-modified-date"
)

type ClassificationLevelColor string
const (
	ClassificationLevelColorRED ClassificationLevelColor = "RED"
	ClassificationLevelColorREDBOLD ClassificationLevelColor = "RED_BOLD"
	ClassificationLevelColorORANGE ClassificationLevelColor = "ORANGE"
	ClassificationLevelColorYELLOW ClassificationLevelColor = "YELLOW"
	ClassificationLevelColorGREEN ClassificationLevelColor = "GREEN"
	ClassificationLevelColorBLUE ClassificationLevelColor = "BLUE"
	ClassificationLevelColorNAVY ClassificationLevelColor = "NAVY"
	ClassificationLevelColorTEAL ClassificationLevelColor = "TEAL"
	ClassificationLevelColorPURPLE ClassificationLevelColor = "PURPLE"
	ClassificationLevelColorGREY ClassificationLevelColor = "GREY"
	ClassificationLevelColorLIME ClassificationLevelColor = "LIME"
)

type ClassificationLevelStatus string
const (
	ClassificationLevelStatusDRAFT ClassificationLevelStatus = "DRAFT"
	ClassificationLevelStatusPUBLISHED ClassificationLevelStatus = "PUBLISHED"
	ClassificationLevelStatusARCHIVED ClassificationLevelStatus = "ARCHIVED"
)

// The sort fields for comments. The default sort direction is ascending. To sort in descending order, append a `-` character before the sort field. For example, `fieldName` or `-fieldName`.
type CommentSortOrder string
const (
	CommentSortOrderCreatedDate CommentSortOrder = "created-date"
	CommentSortOrderDescCreatedDate CommentSortOrder = "-created-date"
	CommentSortOrderModifiedDate CommentSortOrder = "modified-date"
	CommentSortOrderDescModifiedDate CommentSortOrder = "-modified-date"
)

// The sort fields for content properties. The default sort direction is ascending. To sort in descending order, append a `-` character before the sort field. For example, `fieldName` or `-fieldName`.
type ContentPropertySortOrder string
const (
	ContentPropertySortOrderKey ContentPropertySortOrder = "key"
	ContentPropertySortOrderDescKey ContentPropertySortOrder = "-key"
)

// The sort fields for hierarchical content types. The default sort direction is ascending. To sort in descending order, append a `-` character before the sort field. For example, `fieldName` or `-fieldName`.
type ContentSortOrder string
const (
	ContentSortOrderCreatedDate ContentSortOrder = "created-date"
	ContentSortOrderDescCreatedDate ContentSortOrder = "-created-date"
	ContentSortOrderId ContentSortOrder = "id"
	ContentSortOrderDescId ContentSortOrder = "-id"
	ContentSortOrderModifiedDate ContentSortOrder = "modified-date"
	ContentSortOrderDescModifiedDate ContentSortOrder = "-modified-date"
	ContentSortOrderChildPosition ContentSortOrder = "child-position"
	ContentSortOrderDescChildPosition ContentSortOrder = "-child-position"
	ContentSortOrderTitle ContentSortOrder = "title"
	ContentSortOrderDescTitle ContentSortOrder = "-title"
)

// The status of the content.
type ContentStatus string
const (
	ContentStatusCurrent ContentStatus = "current"
	ContentStatusDraft ContentStatus = "draft"
	ContentStatusArchived ContentStatus = "archived"
	ContentStatusHistorical ContentStatus = "historical"
	ContentStatusTrashed ContentStatus = "trashed"
	ContentStatusDeleted ContentStatus = "deleted"
	ContentStatusAny ContentStatus = "any"
)

// The formats a custom content body can be represented as. A subset of BodyRepresentation.
type CustomContentBodyRepresentation string
const (
	CustomContentBodyRepresentationRaw CustomContentBodyRepresentation = "raw"
	CustomContentBodyRepresentationStorage CustomContentBodyRepresentation = "storage"
	CustomContentBodyRepresentationAtlasDocFormat CustomContentBodyRepresentation = "atlas_doc_format"
)

// The sort fields for custom content. The default sort direction is ascending. To sort in descending order, append a `-` character before the sort field. For example, `fieldName` or `-fieldName`.
type CustomContentSortOrder string
const (
	CustomContentSortOrderId CustomContentSortOrder = "id"
	CustomContentSortOrderDescId CustomContentSortOrder = "-id"
	CustomContentSortOrderCreatedDate CustomContentSortOrder = "created-date"
	CustomContentSortOrderDescCreatedDate CustomContentSortOrder = "-created-date"
	CustomContentSortOrderModifiedDate CustomContentSortOrder = "modified-date"
	CustomContentSortOrderDescModifiedDate CustomContentSortOrder = "-modified-date"
	CustomContentSortOrderTitle CustomContentSortOrder = "title"
	CustomContentSortOrderDescTitle CustomContentSortOrder = "-title"
)

// Inline comment resolution status
type InlineCommentResolutionStatus string
const (
	InlineCommentResolutionStatusOpen InlineCommentResolutionStatus = "open"
	InlineCommentResolutionStatusReopened InlineCommentResolutionStatus = "reopened"
	InlineCommentResolutionStatusResolved InlineCommentResolutionStatus = "resolved"
	InlineCommentResolutionStatusDangling InlineCommentResolutionStatus = "dangling"
)

// The sort fields for labels. The default sort direction is ascending. To sort in descending order, append a `-` character before the sort field. For example, `fieldName` or `-fieldName`.
type LabelSortOrder string
const (
	LabelSortOrderCreatedDate LabelSortOrder = "created-date"
	LabelSortOrderDescCreatedDate LabelSortOrder = "-created-date"
	LabelSortOrderId LabelSortOrder = "id"
	LabelSortOrderDescId LabelSortOrder = "-id"
	LabelSortOrderName LabelSortOrder = "name"
	LabelSortOrderDescName LabelSortOrder = "-name"
)

// The status of the content.
type OnlyArchivedAndCurrentContentStatus string
const (
	OnlyArchivedAndCurrentContentStatusCurrent OnlyArchivedAndCurrentContentStatus = "current"
	OnlyArchivedAndCurrentContentStatusArchived OnlyArchivedAndCurrentContentStatus = "archived"
)

// The sort fields for pages. The default sort direction is ascending. To sort in descending order, append a `-` character before the sort field. For example, `fieldName` or `-fieldName`.
type PageSortOrder string
const (
	PageSortOrderId PageSortOrder = "id"
	PageSortOrderDescId PageSortOrder = "-id"
	PageSortOrderCreatedDate PageSortOrder = "created-date"
	PageSortOrderDescCreatedDate PageSortOrder = "-created-date"
	PageSortOrderModifiedDate PageSortOrder = "modified-date"
	PageSortOrderDescModifiedDate PageSortOrder = "-modified-date"
	PageSortOrderTitle PageSortOrder = "title"
	PageSortOrderDescTitle PageSortOrder = "-title"
)

// Content type of the parent, or null if there is no parent.
type ParentContentType string
const (
	ParentContentTypePage ParentContentType = "page"
	ParentContentTypeWhiteboard ParentContentType = "whiteboard"
	ParentContentTypeDatabase ParentContentType = "database"
	ParentContentTypeEmbed ParentContentType = "embed"
	ParentContentTypeFolder ParentContentType = "folder"
)

// The primary formats a body can be represented as. A subset of BodyRepresentation. These formats are the only allowed formats in certain use cases.
type PrimaryBodyRepresentation string
const (
	PrimaryBodyRepresentationStorage PrimaryBodyRepresentation = "storage"
	PrimaryBodyRepresentationAtlasDocFormat PrimaryBodyRepresentation = "atlas_doc_format"
)

// The primary formats a body can be represented as. A subset of BodyRepresentation. These formats are the only allowed formats in certain use cases.
type PrimaryBodyRepresentationSingle string
const (
	PrimaryBodyRepresentationSingleStorage PrimaryBodyRepresentationSingle = "storage"
	PrimaryBodyRepresentationSingleAtlasDocFormat PrimaryBodyRepresentationSingle = "atlas_doc_format"
	PrimaryBodyRepresentationSingleView PrimaryBodyRepresentationSingle = "view"
	PrimaryBodyRepresentationSingleExportView PrimaryBodyRepresentationSingle = "export_view"
	PrimaryBodyRepresentationSingleAnonymousExportView PrimaryBodyRepresentationSingle = "anonymous_export_view"
	PrimaryBodyRepresentationSingleStyledView PrimaryBodyRepresentationSingle = "styled_view"
	PrimaryBodyRepresentationSingleEditor PrimaryBodyRepresentationSingle = "editor"
)

// The principal type.
type PrincipalType string
const (
	PrincipalTypeUSER PrincipalType = "USER"
	PrincipalTypeGROUP PrincipalType = "GROUP"
	PrincipalTypeACCESSCLASS PrincipalType = "ACCESS_CLASS"
)

// The role type.
type RoleType string
const (
	RoleTypeSYSTEM RoleType = "SYSTEM"
	RoleTypeCUSTOM RoleType = "CUSTOM"
)

// The formats a space description can be represented as. A subset of BodyRepresentation.
type SpaceDescriptionBodyRepresentation string
const (
	SpaceDescriptionBodyRepresentationPlain SpaceDescriptionBodyRepresentation = "plain"
	SpaceDescriptionBodyRepresentationView SpaceDescriptionBodyRepresentation = "view"
)

// The sort fields for spaces. The default sort direction is ascending. To sort in descending order, append a `-` character before the sort field. For example, `fieldName` or `-fieldName`.
type SpaceSortOrder string
const (
	SpaceSortOrderId SpaceSortOrder = "id"
	SpaceSortOrderDescId SpaceSortOrder = "-id"
	SpaceSortOrderKey SpaceSortOrder = "key"
	SpaceSortOrderDescKey SpaceSortOrder = "-key"
	SpaceSortOrderName SpaceSortOrder = "name"
	SpaceSortOrderDescName SpaceSortOrder = "-name"
)

// The status of the space.
type SpaceStatus string
const (
	SpaceStatusCurrent SpaceStatus = "current"
	SpaceStatusArchived SpaceStatus = "archived"
)

// The type of space.
type SpaceType string
const (
	SpaceTypeGlobal SpaceType = "global"
	SpaceTypeCollaboration SpaceType = "collaboration"
	SpaceTypeKnowledgeBase SpaceType = "knowledge_base"
	SpaceTypePersonal SpaceType = "personal"
	SpaceTypeSystem SpaceType = "system"
	SpaceTypeOnboarding SpaceType = "onboarding"
	SpaceTypeXflowSampleSpace SpaceType = "xflow_sample_space"
)

// The sort fields for versions. The default sort direction is ascending. To sort in descending order, append a `-` character before the sort field. For example, `fieldName` or `-fieldName`.
type VersionSortOrder string
const (
	VersionSortOrderModifiedDate VersionSortOrder = "modified-date"
	VersionSortOrderDescModifiedDate VersionSortOrder = "-modified-date"
)

type AbstractPageLinks struct {
	Webui *string `json:"webui"`
	Editui *string `json:"editui"`
	Tinyui *string `json:"tinyui"`
}

type Ancestor struct {
	Id *string `json:"id"`
	Type *AncestorType `json:"type"`
}

type AttachmentBulk struct {
	Id *string `json:"id"`
	Status *ContentStatus `json:"status"`
	Title *string `json:"title"`
	CreatedAt *time.Time `json:"createdAt"`
	PageId *string `json:"pageId"`
	BlogPostId *string `json:"blogPostId"`
	CustomContentId *string `json:"customContentId"`
	MediaType *string `json:"mediaType"`
	MediaTypeDescription *string `json:"mediaTypeDescription"`
	Comment *string `json:"comment"`
	FileId *string `json:"fileId"`
	FileSize *int64 `json:"fileSize"`
	WebuiLink *string `json:"webuiLink"`
	DownloadLink *string `json:"downloadLink"`
	Version *Version `json:"version"`
	Links *AttachmentLinks `json:"_links"`
}

type AttachmentCommentModel struct {
	Id *string `json:"id"`
	Status *ContentStatus `json:"status"`
	Title *string `json:"title"`
	AttachmentId *string `json:"attachmentId"`
	Version *Version `json:"version"`
	Body *BodySingle `json:"body"`
	Links *CommentLinks `json:"_links"`
}

type AttachmentLinks struct {
	Webui *string `json:"webui"`
	Download *string `json:"download"`
}

type AttachmentSingle struct {
	Id *string `json:"id"`
	Status *ContentStatus `json:"status"`
	Title *string `json:"title"`
	CreatedAt *time.Time `json:"createdAt"`
	PageId *string `json:"pageId"`
	BlogPostId *string `json:"blogPostId"`
	CustomContentId *string `json:"customContentId"`
	MediaType *string `json:"mediaType"`
	MediaTypeDescription *string `json:"mediaTypeDescription"`
	Comment *string `json:"comment"`
	FileId *string `json:"fileId"`
	FileSize *int64 `json:"fileSize"`
	WebuiLink *string `json:"webuiLink"`
	DownloadLink *string `json:"downloadLink"`
	Version *Version `json:"version"`
	Labels map[string]interface{} `json:"labels"`
	Properties map[string]interface{} `json:"properties"`
	Operations map[string]interface{} `json:"operations"`
	Versions map[string]interface{} `json:"versions"`
	Links *AttachmentLinks `json:"_links"`
}

type AttachmentVersion struct {
	CreatedAt *time.Time `json:"createdAt"`
	Message *string `json:"message"`
	Number int `json:"number"`
	MinorEdit *bool `json:"minorEdit"`
	AuthorId *string `json:"authorId"`
	Attachment *VersionedEntity `json:"attachment"`
}

type BlogPostBulk struct {
	Id *string `json:"id"`
	Status *BlogPostContentStatus `json:"status"`
	Title *string `json:"title"`
	SpaceId *string `json:"spaceId"`
	AuthorId *string `json:"authorId"`
	CreatedAt *time.Time `json:"createdAt"`
	Version *Version `json:"version"`
	Body *BodyBulk `json:"body"`
	Links *AbstractPageLinks `json:"_links"`
}

// Contains fields for each representation type requested.
type BodyBulk struct {
	Storage *BodyType `json:"storage"`
	AtlasDocFormat *BodyType `json:"atlas_doc_format"`
}

// Contains fields for each representation type requested.
type BodySingle struct {
	Storage *BodyType `json:"storage"`
	AtlasDocFormat *BodyType `json:"atlas_doc_format"`
	View *BodyType `json:"view"`
}

type BodyType struct {
	Representation *string `json:"representation"`
	Value *string `json:"value"`
}

type ChildPage struct {
	Id *string `json:"id"`
	Status *OnlyArchivedAndCurrentContentStatus `json:"status"`
	Title *string `json:"title"`
	SpaceId *string `json:"spaceId"`
	ChildPosition int `json:"childPosition"`
}

type ChildrenResponse struct {
	Id *string `json:"id"`
	Status *OnlyArchivedAndCurrentContentStatus `json:"status"`
	Title *string `json:"title"`
	Type *string `json:"type"`
	SpaceId *string `json:"spaceId"`
	ChildPosition int `json:"childPosition"`
}

// A unit of [data classification](https://support.atlassian.com/security-and-access-policies/docs/what-is-data-classification/) defined by an organiation. A classification level may be associated with specific storage and handling requirements or expectations.
type ClassificationLevel struct {
	Id *string `json:"id"`
	Status *ClassificationLevelStatus `json:"status"`
	Order *float64 `json:"order"`
	Name *string `json:"name"`
	Description *string `json:"description"`
	Guideline *string `json:"guideline"`
	Color *ClassificationLevelColor `json:"color"`
}

type CommentLinks struct {
	Webui *string `json:"webui"`
}

type ContentProperty struct {
	Id *string `json:"id"`
	Key *string `json:"key"`
	Value interface{} `json:"value"`
	Version *Version `json:"version"`
}

type ContentPropertyCreateRequest struct {
	Key *string `json:"key"`
	Value interface{} `json:"value"`
}

type ContentPropertyUpdateRequest struct {
	Key *string `json:"key"`
	Value interface{} `json:"value"`
	Version map[string]interface{} `json:"version"`
}

// Contains fields for each representation type requested.
type CustomContentBodyBulk struct {
	Raw *BodyType `json:"raw"`
	Storage *BodyType `json:"storage"`
	AtlasDocFormat *BodyType `json:"atlas_doc_format"`
}

type CustomContentBulk struct {
	Id *string `json:"id"`
	Type *string `json:"type"`
	Status *ContentStatus `json:"status"`
	Title *string `json:"title"`
	SpaceId *string `json:"spaceId"`
	PageId *string `json:"pageId"`
	BlogPostId *string `json:"blogPostId"`
	CustomContentId *string `json:"customContentId"`
	AuthorId *string `json:"authorId"`
	CreatedAt *time.Time `json:"createdAt"`
	Version *Version `json:"version"`
	Body *CustomContentBodyBulk `json:"body"`
	Links *CustomContentLinks `json:"_links"`
}

type CustomContentLinks struct {
	Webui *string `json:"webui"`
}

type DescendantsResponse struct {
	Id *string `json:"id"`
	Status *OnlyArchivedAndCurrentContentStatus `json:"status"`
	Title *string `json:"title"`
	Type *string `json:"type"`
	ParentId *string `json:"parentId"`
	Depth int `json:"depth"`
	ChildPosition int `json:"childPosition"`
}

type DetailedVersion struct {
	Number int `json:"number"`
	AuthorId *string `json:"authorId"`
	Message *string `json:"message"`
	CreatedAt *time.Time `json:"createdAt"`
	MinorEdit *bool `json:"minorEdit"`
	ContentTypeModified *bool `json:"contentTypeModified"`
	Collaborators []string `json:"collaborators"`
	PrevVersion int `json:"prevVersion"`
	NextVersion int `json:"nextVersion"`
}

type InlineCommentProperties struct {
	InlineMarkerRef *string `json:"inlineMarkerRef"`
	InlineOriginalSelection *string `json:"inlineOriginalSelection"`
}

type Label struct {
	Id *string `json:"id"`
	Name *string `json:"name"`
	Prefix *string `json:"prefix"`
}

type Like struct {
	AccountId *string `json:"accountId"`
}

type MultiEntityLinks struct {
	Next *string `json:"next"`
	Base *string `json:"base"`
}

type Operation struct {
	Operation *string `json:"operation"`
	TargetType *string `json:"targetType"`
}

type OptionalFieldLinks struct {
	Self *string `json:"self"`
}

type OptionalFieldMeta struct {
	HasMore *bool `json:"hasMore"`
	Cursor *string `json:"cursor"`
}

type PageBulk struct {
	Id *string `json:"id"`
	Status *ContentStatus `json:"status"`
	Title *string `json:"title"`
	SpaceId *string `json:"spaceId"`
	ParentId *string `json:"parentId"`
	ParentType *ParentContentType `json:"parentType"`
	Position int `json:"position"`
	AuthorId *string `json:"authorId"`
	OwnerId *string `json:"ownerId"`
	LastOwnerId *string `json:"lastOwnerId"`
	Subtype *string `json:"subtype"`
	CreatedAt *time.Time `json:"createdAt"`
	Version *Version `json:"version"`
	Body *BodyBulk `json:"body"`
	Links *AbstractPageLinks `json:"_links"`
}

type PageCommentModel struct {
	Id *string `json:"id"`
	Status *ContentStatus `json:"status"`
	Title *string `json:"title"`
	PageId *string `json:"pageId"`
	Version *Version `json:"version"`
	Body *BodyBulk `json:"body"`
	Links *CommentLinks `json:"_links"`
}

type PageInlineCommentModel struct {
	Id *string `json:"id"`
	Status *ContentStatus `json:"status"`
	Title *string `json:"title"`
	PageId *string `json:"pageId"`
	Version *Version `json:"version"`
	Body *BodyBulk `json:"body"`
	ResolutionStatus *InlineCommentResolutionStatus `json:"resolutionStatus"`
	Properties *InlineCommentProperties `json:"properties"`
	Links *CommentLinks `json:"_links"`
}

type PageSingle struct {
	Id *string `json:"id"`
	Status *ContentStatus `json:"status"`
	Title *string `json:"title"`
	SpaceId *string `json:"spaceId"`
	ParentId *string `json:"parentId"`
	ParentType *ParentContentType `json:"parentType"`
	Position int `json:"position"`
	AuthorId *string `json:"authorId"`
	OwnerId *string `json:"ownerId"`
	LastOwnerId *string `json:"lastOwnerId"`
	CreatedAt *time.Time `json:"createdAt"`
	Version *Version `json:"version"`
	Body *BodySingle `json:"body"`
	Labels map[string]interface{} `json:"labels"`
	Properties map[string]interface{} `json:"properties"`
	Operations map[string]interface{} `json:"operations"`
	Likes map[string]interface{} `json:"likes"`
	Versions map[string]interface{} `json:"versions"`
	IsFavoritedByCurrentUser *bool `json:"isFavoritedByCurrentUser"`
	Links *AbstractPageLinks `json:"_links"`
}

type PageVersion struct {
	CreatedAt *time.Time `json:"createdAt"`
	Message *string `json:"message"`
	Number int `json:"number"`
	MinorEdit *bool `json:"minorEdit"`
	AuthorId *string `json:"authorId"`
	Page *VersionedEntity `json:"page"`
}

// The list of operations permitted on entity.
type PermittedOperationsResponse struct {
	Operations []Operation `json:"operations"`
}

// The principal of the role assignment.
type Principal struct {
	PrincipalType *PrincipalType `json:"principalType"`
	PrincipalId *string `json:"principalId"`
}

type RedactionPointerResponse struct {
	Pointer *string `json:"pointer"`
	From int `json:"from"`
	To int `json:"to"`
	Reason *string `json:"reason"`
	RedactionId *string `json:"redactionId"`
}

// Response containing details of all redactions that were applied to the content. Each redaction includes a unique ID for restoration, except that code block redactions cannot be restored.
type RedactionResponse struct {
	Body *RedactionSectionResponse `json:"body"`
	Title *RedactionSectionResponse `json:"title"`
}

type RedactionSectionResponse struct {
	Redactions []RedactionPointerResponse `json:"redactions"`
}

type SpaceBulk struct {
	Id *string `json:"id"`
	Key *string `json:"key"`
	Name *string `json:"name"`
	Type *SpaceType `json:"type"`
	Status *SpaceStatus `json:"status"`
	AuthorId *string `json:"authorId"`
	CurrentActiveAlias *string `json:"currentActiveAlias"`
	CreatedAt *time.Time `json:"createdAt"`
	HomepageId *string `json:"homepageId"`
	Description *SpaceDescription `json:"description"`
	Icon *SpaceIcon `json:"icon"`
	Links *SpaceLinks `json:"_links"`
}

// Contains fields for each representation type requested.
type SpaceDescription struct {
	Plain *BodyType `json:"plain"`
	View *BodyType `json:"view"`
}

// The icon of the space
type SpaceIcon struct {
	Path *string `json:"path"`
	ApiDownloadLink *string `json:"apiDownloadLink"`
}

type SpaceLinks struct {
	Webui *string `json:"webui"`
}

type SpacePermissionAssignment struct {
	Id *string `json:"id"`
	Principal map[string]interface{} `json:"principal"`
	Operation map[string]interface{} `json:"operation"`
}

type SpaceProperty struct {
	Id *string `json:"id"`
	Key *string `json:"key"`
	Value interface{} `json:"value"`
	CreatedAt *time.Time `json:"createdAt"`
	CreatedBy *string `json:"createdBy"`
	Version map[string]interface{} `json:"version"`
}

type SpacePropertyCreateRequest struct {
	Key *string `json:"key"`
	Value interface{} `json:"value"`
}

type SpacePropertyUpdateRequest struct {
	Key *string `json:"key"`
	Value interface{} `json:"value"`
	Version map[string]interface{} `json:"version"`
}

type SpaceSingle struct {
	Id *string `json:"id"`
	Key *string `json:"key"`
	Name *string `json:"name"`
	Type *SpaceType `json:"type"`
	Status *SpaceStatus `json:"status"`
	AuthorId *string `json:"authorId"`
	CreatedAt *time.Time `json:"createdAt"`
	HomepageId *string `json:"homepageId"`
	Description *SpaceDescription `json:"description"`
	Icon *SpaceIcon `json:"icon"`
	Labels map[string]interface{} `json:"labels"`
	Properties map[string]interface{} `json:"properties"`
	Operations map[string]interface{} `json:"operations"`
	Permissions map[string]interface{} `json:"permissions"`
	Links *SpaceLinks `json:"_links"`
}

type Version struct {
	CreatedAt *time.Time `json:"createdAt"`
	Message *string `json:"message"`
	Number int `json:"number"`
	MinorEdit *bool `json:"minorEdit"`
	AuthorId *string `json:"authorId"`
}

type VersionedEntity struct {
	Title *string `json:"title"`
	Id *string `json:"id"`
	Body *BodyBulk `json:"body"`
}

type ContentClassificationLevelDeleteRequest struct {
	Status string `json:"status"`
}

type ContentClassificationLevelUpdateRequest struct {
	Id string `json:"id"`
	Status string `json:"status"`
}

type PageCreateRequest struct {
	SpaceId string `json:"spaceId"`
	Status *string `json:"status"`
	Title *string `json:"title"`
	ParentId *string `json:"parentId"`
	Body interface{} `json:"body"`
	Subtype *string `json:"subtype"`
}

type PageTitleUpdateRequest struct {
	Status string `json:"status"`
	Title string `json:"title"`
}

type PageUpdateRequest struct {
	Id string `json:"id"`
	Status string `json:"status"`
	Title string `json:"title"`
	SpaceId interface{} `json:"spaceId"`
	ParentId interface{} `json:"parentId"`
	OwnerId interface{} `json:"ownerId"`
	Body interface{} `json:"body"`
	Version map[string]interface{} `json:"version"`
}

type RedactionRequest struct {
	CreatedAt time.Time `json:"createdAt"`
	CleanHistory *bool `json:"cleanHistory"`
	VersionNumber int `json:"versionNumber"`
	Body map[string]interface{} `json:"body"`
	Title map[string]interface{} `json:"title"`
}

type SpaceCreateRequest struct {
	Name string `json:"name"`
	Key *string `json:"key"`
	Alias *string `json:"alias"`
	Description map[string]interface{} `json:"description"`
	CopySpaceAccessConfiguration int `json:"copySpaceAccessConfiguration"`
	CreatePrivateSpace *bool `json:"createPrivateSpace"`
	TemplateKey *string `json:"templateKey"`
}

type SpaceDefaultClassificationLevelUpdateRequest struct {
	Id string `json:"id"`
}
