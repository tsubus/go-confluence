package confluence_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	confluence "github.com/umats/go-confluence"
)

type testQueryParams struct {
	Name  *string  `json:"name"`
	Tags  []string `json:"tags"`
	Limit *int     `json:"limit"`
}

type recordedRequest struct {
	method      string
	path        string
	rawQuery    string
	authUser    string
	authPass    string
	contentType string
	body        []byte
	handlerErr  error
}

type serverCase struct {
	method      string
	path        string
	rawQuery    string
	statusCode  int
	response    string
	contentType string
}

func newV2TestServer(t *testing.T, expected serverCase, record *recordedRequest) *httptest.Server {
	t.Helper()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		record.method = r.Method
		record.path = r.URL.Path
		record.rawQuery = r.URL.RawQuery
		record.contentType = r.Header.Get("Content-Type")
		record.authUser, record.authPass, _ = r.BasicAuth()

		body, err := ioReadAll(r)
		if err != nil {
			record.handlerErr = err
		}
		record.body = body

		if expected.contentType != "" {
			w.Header().Set("Content-Type", expected.contentType)
		}
		w.WriteHeader(expected.statusCode)
		_, writeErr := w.Write([]byte(expected.response))
		if writeErr != nil && record.handlerErr == nil {
			record.handlerErr = writeErr
		}
	})

	return httptest.NewServer(handler)
}

func ioReadAll(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

func TestPageService_Methods(t *testing.T) {
	tests := []struct {
		name     string
		execute  func(client *confluence.Client) error
		want     serverCase
		wantBody bool
	}{
		{
			name: "get page",
			execute: func(client *confluence.Client) error {
				_, err := client.Page().Get(context.Background(), 77, nil)
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/pages/77",
				statusCode:  http.StatusOK,
				response:    `{"id":"77"}`,
				contentType: "application/json",
			},
		},
		{
			name: "list pages",
			execute: func(client *confluence.Client) error {
				limit := 10
				_, err := client.Page().List(context.Background(), &confluence.PageListParams{Limit: &limit})
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/pages",
				rawQuery:    "limit=10",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"id":"42"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "create page",
			execute: func(client *confluence.Client) error {
				rootLevel := true
				req := confluence.PageCreateRequest{SpaceId: "SPACE"}
				_, err := client.Page().Create(context.Background(), &confluence.PageCreateParams{RootLevel: &rootLevel}, req)
				return err
			},
			want: serverCase{
				method:      http.MethodPost,
				path:        "/wiki/api/v2/pages",
				rawQuery:    "root-level=true",
				statusCode:  http.StatusOK,
				response:    `{"id":"99"}`,
				contentType: "application/json",
			},
			wantBody: true,
		},
		{
			name: "update page",
			execute: func(client *confluence.Client) error {
				req := confluence.PageUpdateRequest{Id: "77"}
				_, err := client.Page().Update(context.Background(), 77, req)
				return err
			},
			want: serverCase{
				method:      http.MethodPut,
				path:        "/wiki/api/v2/pages/77",
				statusCode:  http.StatusOK,
				response:    `{"id":"77"}`,
				contentType: "application/json",
			},
			wantBody: true,
		},
		{
			name: "delete page",
			execute: func(client *confluence.Client) error {
				return client.Page().Delete(context.Background(), 77)
			},
			want: serverCase{
				method:     http.MethodDelete,
				path:       "/wiki/api/v2/pages/77",
				statusCode: http.StatusNoContent,
				response:   "",
			},
		},
		{
			name: "update page title",
			execute: func(client *confluence.Client) error {
				req := confluence.PageTitleUpdateRequest{Status: "current"}
				_, err := client.Page().UpdateTitle(context.Background(), 77, req)
				return err
			},
			want: serverCase{
				method:      http.MethodPut,
				path:        "/wiki/api/v2/pages/77/title",
				statusCode:  http.StatusOK,
				response:    `{"id":"77"}`,
				contentType: "application/json",
			},
			wantBody: true,
		},
		{
			name: "get page attachments",
			execute: func(client *confluence.Client) error {
				limit := 5
				_, err := client.Page().GetAttachments(context.Background(), 77, &confluence.PageAttachmentsParams{Limit: &limit})
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/pages/77/attachments",
				rawQuery:    "limit=5",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"id":"att1"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get page labels",
			execute: func(client *confluence.Client) error {
				limit := 5
				_, err := client.Page().GetLabels(context.Background(), 77, &confluence.PageLabelsParams{Limit: &limit})
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/pages/77/labels",
				rawQuery:    "limit=5",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"id":"l1"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get page operations",
			execute: func(client *confluence.Client) error {
				_, err := client.Page().GetOperations(context.Background(), 77)
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/pages/77/operations",
				statusCode:  http.StatusOK,
				response:    `{"operations":[]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get page versions",
			execute: func(client *confluence.Client) error {
				limit := 5
				_, err := client.Page().GetVersions(context.Background(), 77, &confluence.PageVersionsParams{Limit: &limit})
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/pages/77/versions",
				rawQuery:    "limit=5",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"number":1}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get page version details",
			execute: func(client *confluence.Client) error {
				_, err := client.Page().GetVersionDetails(context.Background(), 77, 3)
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/pages/77/versions/3",
				statusCode:  http.StatusOK,
				response:    `{"number":3}`,
				contentType: "application/json",
			},
		},
		{
			name: "get page ancestors",
			execute: func(client *confluence.Client) error {
				limit := 5
				_, err := client.Page().GetAncestors(context.Background(), 77, &confluence.PageAncestorsParams{Limit: &limit})
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/pages/77/ancestors",
				rawQuery:    "limit=5",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"id":"1"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get page descendants",
			execute: func(client *confluence.Client) error {
				depth := 2
				_, err := client.Page().GetDescendants(context.Background(), 77, &confluence.PageDescendantsParams{Depth: &depth})
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/pages/77/descendants",
				rawQuery:    "depth=2",
				statusCode:  http.StatusOK,
				response:    `{"page":[],"attachment":[]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get page children",
			execute: func(client *confluence.Client) error {
				limit := 5
				_, err := client.Page().GetChildren(context.Background(), 77, &confluence.PageChildrenParams{Limit: &limit})
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/pages/77/children",
				rawQuery:    "limit=5",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"id":"1"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get page direct children",
			execute: func(client *confluence.Client) error {
				limit := 5
				_, err := client.Page().GetDirectChildren(context.Background(), 77, &confluence.PageChildrenParams{Limit: &limit})
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/pages/77/direct-children",
				rawQuery:    "limit=5",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"id":"1"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get page footer comments",
			execute: func(client *confluence.Client) error {
				limit := 5
				_, err := client.Page().GetFooterComments(context.Background(), 77, &confluence.PageCommentsParams{Limit: &limit})
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/pages/77/footer-comments",
				rawQuery:    "limit=5",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"id":"c1"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get page inline comments",
			execute: func(client *confluence.Client) error {
				limit := 5
				_, err := client.Page().GetInlineComments(context.Background(), 77, &confluence.PageCommentsParams{Limit: &limit})
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/pages/77/inline-comments",
				rawQuery:    "limit=5",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"id":"c2"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get page content properties",
			execute: func(client *confluence.Client) error {
				limit := 5
				_, err := client.Page().GetContentProperties(
					context.Background(),
					77,
					&confluence.ContentPropertiesParams{Limit: &limit},
				)
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/pages/77/properties",
				rawQuery:    "limit=5",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"key":"k1"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get page content property by id",
			execute: func(client *confluence.Client) error {
				_, err := client.Page().GetContentPropertyByID(context.Background(), 77, 99)
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/pages/77/properties/99",
				statusCode:  http.StatusOK,
				response:    `{"key":"k1"}`,
				contentType: "application/json",
			},
		},
		{
			name: "create page content property",
			execute: func(client *confluence.Client) error {
				key := "k1"
				req := confluence.ContentPropertyCreateRequest{Key: &key}
				_, err := client.Page().CreateContentProperty(context.Background(), 77, req)
				return err
			},
			want: serverCase{
				method:      http.MethodPost,
				path:        "/wiki/api/v2/pages/77/properties",
				statusCode:  http.StatusOK,
				response:    `{"key":"k1"}`,
				contentType: "application/json",
			},
			wantBody: true,
		},
		{
			name: "update page content property",
			execute: func(client *confluence.Client) error {
				key := "k1"
				req := confluence.ContentPropertyUpdateRequest{Key: &key}
				_, err := client.Page().UpdateContentProperty(context.Background(), 77, 99, req)
				return err
			},
			want: serverCase{
				method:      http.MethodPut,
				path:        "/wiki/api/v2/pages/77/properties/99",
				statusCode:  http.StatusOK,
				response:    `{"key":"k1"}`,
				contentType: "application/json",
			},
			wantBody: true,
		},
		{
			name: "delete page content property",
			execute: func(client *confluence.Client) error {
				return client.Page().DeleteContentProperty(context.Background(), 77, 99)
			},
			want: serverCase{
				method:     http.MethodDelete,
				path:       "/wiki/api/v2/pages/77/properties/99",
				statusCode: http.StatusNoContent,
				response:   "",
			},
		},
		{
			name: "get page custom content",
			execute: func(client *confluence.Client) error {
				_, err := client.Page().GetCustomContent(context.Background(), 77, &confluence.PageCustomContentParams{Type: "foo"})
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/pages/77/custom-content",
				rawQuery:    "type=foo",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"id":"cc1"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get page like count",
			execute: func(client *confluence.Client) error {
				_, err := client.Page().GetLikeCount(context.Background(), 77)
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/pages/77/likes/count",
				statusCode:  http.StatusOK,
				response:    `{"count":5}`,
				contentType: "application/json",
			},
		},
		{
			name: "get page like users",
			execute: func(client *confluence.Client) error {
				limit := 5
				_, err := client.Page().GetLikeUsers(context.Background(), 77, &confluence.PageLikeUsersParams{Limit: &limit})
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/pages/77/likes/users",
				rawQuery:    "limit=5",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"id":"u1"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get page classification level",
			execute: func(client *confluence.Client) error {
				_, err := client.Page().GetClassificationLevel(context.Background(), 77)
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/pages/77/classification-level",
				statusCode:  http.StatusOK,
				response:    `{"id":"cl1"}`,
				contentType: "application/json",
			},
		},
		{
			name: "update page classification level",
			execute: func(client *confluence.Client) error {
				req := confluence.ContentClassificationLevelUpdateRequest{Id: "cl1"}
				_, err := client.Page().UpdateClassificationLevel(context.Background(), 77, req)
				return err
			},
			want: serverCase{
				method:      http.MethodPut,
				path:        "/wiki/api/v2/pages/77/classification-level",
				statusCode:  http.StatusOK,
				response:    `{"id":"cl1"}`,
				contentType: "application/json",
			},
			wantBody: true,
		},
		{
			name: "reset page classification level",
			execute: func(client *confluence.Client) error {
				req := confluence.ContentClassificationLevelDeleteRequest{Status: "current"}
				_, err := client.Page().ResetClassificationLevel(context.Background(), 77, req)
				return err
			},
			want: serverCase{
				method:      http.MethodPost,
				path:        "/wiki/api/v2/pages/77/classification-level/reset",
				statusCode:  http.StatusOK,
				response:    `{"id":"cl1"}`,
				contentType: "application/json",
			},
			wantBody: true,
		},
		{
			name: "redact page",
			execute: func(client *confluence.Client) error {
				req := confluence.RedactionRequest{CreatedAt: time.Now()}
				_, err := client.Page().Redact(context.Background(), 77, req)
				return err
			},
			want: serverCase{
				method:      http.MethodPost,
				path:        "/wiki/api/v2/pages/77/redact",
				statusCode:  http.StatusOK,
				response:    `{"id":"r1"}`,
				contentType: "application/json",
			},
			wantBody: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var record recordedRequest
			expected := tt.want
			server := newV2TestServer(t, expected, &record)
			defer server.Close()

			client, err := confluence.NewClient(server.URL, confluence.WithBasicAuth("user", "pass"))
			require.NoError(t, err)

			err = tt.execute(client)
			require.NoError(t, err)
			require.NoError(t, record.handlerErr)
			require.Equal(t, expected.method, record.method)
			require.Equal(t, expected.path, record.path)
			require.Equal(t, expected.rawQuery, record.rawQuery)
			if tt.wantBody {
				require.NotEmpty(t, record.body)
				if expected.contentType != "" {
					require.Equal(t, expected.contentType, record.contentType)
				}
			}
		})
	}
}

func TestSpaceService_Methods(t *testing.T) {
	tests := []struct {
		name    string
		execute func(client *confluence.Client) error
		want    serverCase
	}{
		{
			name: "get space",
			execute: func(client *confluence.Client) error {
				_, err := client.Space().Get(context.Background(), 42, nil)
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/spaces/42",
				statusCode:  http.StatusOK,
				response:    `{"id":"42"}`,
				contentType: "application/json",
			},
		},
		{
			name: "list spaces",
			execute: func(client *confluence.Client) error {
				limit := 5
				_, err := client.Space().List(context.Background(), &confluence.SpaceListParams{Limit: &limit})
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/spaces",
				rawQuery:    "limit=5",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"id":"1"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "create space",
			execute: func(client *confluence.Client) error {
				req := confluence.SpaceCreateRequest{Name: "Engineering"}
				_, err := client.Space().Create(context.Background(), req)
				return err
			},
			want: serverCase{
				method:      http.MethodPost,
				path:        "/wiki/api/v2/spaces",
				statusCode:  http.StatusOK,
				response:    `{"id":"abc"}`,
				contentType: "application/json",
			},
		},
		{
			name: "get space pages",
			execute: func(client *confluence.Client) error {
				limit := 5
				_, err := client.Space().GetPages(context.Background(), 42, &confluence.SpacePagesParams{Limit: &limit})
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/spaces/42/pages",
				rawQuery:    "limit=5",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"id":"p1"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get space blog posts",
			execute: func(client *confluence.Client) error {
				limit := 5
				_, err := client.Space().GetBlogPosts(context.Background(), 42, &confluence.SpaceBlogPostsParams{Limit: &limit})
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/spaces/42/blogposts",
				rawQuery:    "limit=5",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"id":"bp1"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get space labels",
			execute: func(client *confluence.Client) error {
				limit := 3
				_, err := client.Space().GetLabels(context.Background(), 42, &confluence.SpaceLabelsParams{Limit: &limit})
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/spaces/42/labels",
				rawQuery:    "limit=3",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"id":"l1","name":"test-label"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get space content labels",
			execute: func(client *confluence.Client) error {
				limit := 3
				_, err := client.Space().GetContentLabels(context.Background(), 42, &confluence.SpaceLabelsParams{Limit: &limit})
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/spaces/42/content/labels",
				rawQuery:    "limit=3",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"id":"l2","name":"content-label"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get space custom content",
			execute: func(client *confluence.Client) error {
				_, err := client.Space().GetCustomContent(
					context.Background(),
					42,
					&confluence.SpaceCustomContentParams{Type: "foo"},
				)
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/spaces/42/custom-content",
				rawQuery:    "type=foo",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"id":"cc1"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get space operations",
			execute: func(client *confluence.Client) error {
				_, err := client.Space().GetOperations(context.Background(), 42)
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/spaces/42/operations",
				statusCode:  http.StatusOK,
				response:    `{"operations":[]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get space permissions",
			execute: func(client *confluence.Client) error {
				limit := 5
				_, err := client.Space().GetPermissions(context.Background(), 42, &confluence.SpacePermissionsParams{Limit: &limit})
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/spaces/42/permissions",
				rawQuery:    "limit=5",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"id":"perm1"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get space properties",
			execute: func(client *confluence.Client) error {
				limit := 5
				_, err := client.Space().GetProperties(context.Background(), 42, &confluence.SpacePropertiesParams{Limit: &limit})
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/spaces/42/properties",
				rawQuery:    "limit=5",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"key":"k1"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get space property by id",
			execute: func(client *confluence.Client) error {
				_, err := client.Space().GetPropertyByID(context.Background(), 42, 99)
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/spaces/42/properties/99",
				statusCode:  http.StatusOK,
				response:    `{"key":"k1"}`,
				contentType: "application/json",
			},
		},
		{
			name: "create space property",
			execute: func(client *confluence.Client) error {
				key := "k1"
				req := confluence.SpacePropertyCreateRequest{Key: &key}
				_, err := client.Space().CreateProperty(context.Background(), 42, req)
				return err
			},
			want: serverCase{
				method:      http.MethodPost,
				path:        "/wiki/api/v2/spaces/42/properties",
				statusCode:  http.StatusOK,
				response:    `{"key":"k1"}`,
				contentType: "application/json",
			},
		},
		{
			name: "update space property",
			execute: func(client *confluence.Client) error {
				key := "k1"
				req := confluence.SpacePropertyUpdateRequest{Key: &key}
				_, err := client.Space().UpdateProperty(context.Background(), 42, 99, req)
				return err
			},
			want: serverCase{
				method:      http.MethodPut,
				path:        "/wiki/api/v2/spaces/42/properties/99",
				statusCode:  http.StatusOK,
				response:    `{"key":"k1"}`,
				contentType: "application/json",
			},
		},
		{
			name: "delete space property",
			execute: func(client *confluence.Client) error {
				return client.Space().DeleteProperty(context.Background(), 42, 99)
			},
			want: serverCase{
				method:     http.MethodDelete,
				path:       "/wiki/api/v2/spaces/42/properties/99",
				statusCode: http.StatusNoContent,
				response:   "",
			},
		},
		{
			name: "get space default classification level",
			execute: func(client *confluence.Client) error {
				_, err := client.Space().GetDefaultClassificationLevel(context.Background(), 42)
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/spaces/42/classification-level/default",
				statusCode:  http.StatusOK,
				response:    `{"id":"cl1"}`,
				contentType: "application/json",
			},
		},
		{
			name: "update space default classification level",
			execute: func(client *confluence.Client) error {
				req := confluence.SpaceDefaultClassificationLevelUpdateRequest{}
				_, err := client.Space().UpdateDefaultClassificationLevel(context.Background(), 42, req)
				return err
			},
			want: serverCase{
				method:      http.MethodPut,
				path:        "/wiki/api/v2/spaces/42/classification-level/default",
				statusCode:  http.StatusOK,
				response:    `{"id":"cl1"}`,
				contentType: "application/json",
			},
		},
		{
			name: "delete space default classification level",
			execute: func(client *confluence.Client) error {
				req := confluence.SpaceDefaultClassificationLevelUpdateRequest{}
				_, err := client.Space().DeleteDefaultClassificationLevel(context.Background(), 12, req)
				return err
			},
			want: serverCase{
				method:      http.MethodDelete,
				path:        "/wiki/api/v2/spaces/12/classification-level/default",
				statusCode:  http.StatusOK,
				response:    `{"id":"x"}`,
				contentType: "application/json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var record recordedRequest
			expected := tt.want
			server := newV2TestServer(t, expected, &record)
			defer server.Close()

			client, err := confluence.NewClient(server.URL, confluence.WithBasicAuth("user", "pass"))
			require.NoError(t, err)

			err = tt.execute(client)
			require.NoError(t, err)
			require.NoError(t, record.handlerErr)
			require.Equal(t, expected.method, record.method)
			require.Equal(t, expected.path, record.path)
			require.Equal(t, expected.rawQuery, record.rawQuery)
		})
	}
}

func TestAttachmentService_Methods(t *testing.T) {
	tests := []struct {
		name    string
		execute func(client *confluence.Client) error
		want    serverCase
	}{
		{
			name: "list attachments",
			execute: func(client *confluence.Client) error {
				limit := 5
				_, err := client.Attachment().List(context.Background(), &confluence.AttachmentListParams{Limit: &limit})
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/attachments",
				rawQuery:    "limit=5",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"id":"att1"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get attachment",
			execute: func(client *confluence.Client) error {
				_, err := client.Attachment().Get(context.Background(), "att123", nil)
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/attachments/att123",
				statusCode:  http.StatusOK,
				response:    `{"id":"att123"}`,
				contentType: "application/json",
			},
		},
		{
			name: "delete attachment no purge",
			execute: func(client *confluence.Client) error {
				return client.Attachment().Delete(context.Background(), 999, false)
			},
			want: serverCase{
				method:     http.MethodDelete,
				path:       "/wiki/api/v2/attachments/999",
				statusCode: http.StatusNoContent,
				response:   "",
			},
		},
		{
			name: "delete attachment purge",
			execute: func(client *confluence.Client) error {
				return client.Attachment().Delete(context.Background(), 999, true)
			},
			want: serverCase{
				method:     http.MethodDelete,
				path:       "/wiki/api/v2/attachments/999",
				rawQuery:   "purge=true",
				statusCode: http.StatusNoContent,
				response:   "",
			},
		},
		{
			name: "get attachment labels",
			execute: func(client *confluence.Client) error {
				limit := 5
				_, err := client.Attachment().GetLabels(
					context.Background(),
					"123",
					&confluence.AttachmentLabelsParams{Limit: &limit},
				)
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/attachments/123/labels",
				rawQuery:    "limit=5",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"id":"l1"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get attachment operations",
			execute: func(client *confluence.Client) error {
				_, err := client.Attachment().GetOperations(context.Background(), "att123")
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/attachments/att123/operations",
				statusCode:  http.StatusOK,
				response:    `{"operations":[]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get attachment versions",
			execute: func(client *confluence.Client) error {
				limit := 5
				_, err := client.Attachment().GetVersions(
					context.Background(),
					"att123",
					&confluence.AttachmentVersionsParams{Limit: &limit},
				)
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/attachments/att123/versions",
				rawQuery:    "limit=5",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"number":1}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get attachment version details",
			execute: func(client *confluence.Client) error {
				_, err := client.Attachment().GetVersionDetails(context.Background(), "att123", 3)
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/attachments/att123/versions/3",
				statusCode:  http.StatusOK,
				response:    `{"number":3}`,
				contentType: "application/json",
			},
		},
		{
			name: "get attachment comments",
			execute: func(client *confluence.Client) error {
				limit := 5
				_, err := client.Attachment().GetComments(
					context.Background(),
					"att123",
					&confluence.AttachmentCommentsParams{Limit: &limit},
				)
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/attachments/att123/footer-comments",
				rawQuery:    "limit=5",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"id":"c1"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get attachment content properties",
			execute: func(client *confluence.Client) error {
				limit := 5
				_, err := client.Attachment().GetContentProperties(
					context.Background(),
					"att123",
					&confluence.ContentPropertiesParams{Limit: &limit},
				)
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/attachments/att123/properties",
				rawQuery:    "limit=5",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"key":"k1"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get attachment content property by id",
			execute: func(client *confluence.Client) error {
				_, err := client.Attachment().GetContentPropertyByID(context.Background(), "att123", 99)
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/attachments/att123/properties/99",
				statusCode:  http.StatusOK,
				response:    `{"key":"k1"}`,
				contentType: "application/json",
			},
		},
		{
			name: "create attachment content property",
			execute: func(client *confluence.Client) error {
				key := "k1"
				req := confluence.ContentPropertyCreateRequest{Key: &key}
				_, err := client.Attachment().CreateContentProperty(context.Background(), "att123", req)
				return err
			},
			want: serverCase{
				method:      http.MethodPost,
				path:        "/wiki/api/v2/attachments/att123/properties",
				statusCode:  http.StatusOK,
				response:    `{"key":"k1"}`,
				contentType: "application/json",
			},
		},
		{
			name: "update attachment content property",
			execute: func(client *confluence.Client) error {
				key := "k1"
				req := confluence.ContentPropertyUpdateRequest{Key: &key}
				_, err := client.Attachment().UpdateContentProperty(context.Background(), "att123", 99, req)
				return err
			},
			want: serverCase{
				method:      http.MethodPut,
				path:        "/wiki/api/v2/attachments/att123/properties/99",
				statusCode:  http.StatusOK,
				response:    `{"key":"k1"}`,
				contentType: "application/json",
			},
		},
		{
			name: "delete attachment content property",
			execute: func(client *confluence.Client) error {
				return client.Attachment().DeleteContentProperty(context.Background(), "att123", 99)
			},
			want: serverCase{
				method:     http.MethodDelete,
				path:       "/wiki/api/v2/attachments/att123/properties/99",
				statusCode: http.StatusNoContent,
				response:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var record recordedRequest
			expected := tt.want
			server := newV2TestServer(t, expected, &record)
			defer server.Close()

			client, err := confluence.NewClient(server.URL, confluence.WithBasicAuth("user", "pass"))
			require.NoError(t, err)

			err = tt.execute(client)
			require.NoError(t, err)
			require.NoError(t, record.handlerErr)
			require.Equal(t, expected.method, record.method)
			require.Equal(t, expected.path, record.path)
			require.Equal(t, expected.rawQuery, record.rawQuery)
		})
	}
}

func TestLabelService_Methods(t *testing.T) {
	tests := []struct {
		name    string
		execute func(client *confluence.Client) error
		want    serverCase
	}{
		{
			name: "list labels",
			execute: func(client *confluence.Client) error {
				limit := 2
				_, err := client.Label().List(context.Background(), &confluence.LabelListParams{Limit: &limit})
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/labels",
				rawQuery:    "limit=2",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"id":"l1"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get label attachments",
			execute: func(client *confluence.Client) error {
				limit := 5
				_, err := client.Label().GetAttachments(context.Background(), 7, &confluence.LabelAttachmentsParams{Limit: &limit})
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/labels/7/attachments",
				rawQuery:    "limit=5",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"id":"att1"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get label pages",
			execute: func(client *confluence.Client) error {
				limit := 5
				_, err := client.Label().GetPages(context.Background(), 7, &confluence.LabelPagesParams{Limit: &limit})
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/labels/7/pages",
				rawQuery:    "limit=5",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"id":"p1"}]}`,
				contentType: "application/json",
			},
		},
		{
			name: "get label blog posts",
			execute: func(client *confluence.Client) error {
				limit := 5
				_, err := client.Label().GetBlogPosts(context.Background(), 7, &confluence.LabelBlogPostsParams{Limit: &limit})
				return err
			},
			want: serverCase{
				method:      http.MethodGet,
				path:        "/wiki/api/v2/labels/7/blogposts",
				rawQuery:    "limit=5",
				statusCode:  http.StatusOK,
				response:    `{"results":[{"id":"bp1"}]}`,
				contentType: "application/json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var record recordedRequest
			expected := tt.want
			server := newV2TestServer(t, expected, &record)
			defer server.Close()

			client, err := confluence.NewClient(server.URL, confluence.WithBasicAuth("user", "pass"))
			require.NoError(t, err)

			err = tt.execute(client)
			require.NoError(t, err)
			require.NoError(t, record.handlerErr)
			require.Equal(t, expected.method, record.method)
			require.Equal(t, expected.path, record.path)
			require.Equal(t, expected.rawQuery, record.rawQuery)
		})
	}
}

func TestV2DecodeResponse_EmptyBody(t *testing.T) {
	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("")),
	}

	client, err := confluence.NewClient("http://example.com")
	require.NoError(t, err)

	apiErr, err := confluence.DecodeResponseForTest(client, response, &map[string]any{})
	require.NoError(t, err)
	require.Empty(t, apiErr)
}

func TestBuildQuery_OmitsNilFields(t *testing.T) {
	params := &testQueryParams{
		Name:  nil,
		Tags:  nil,
		Limit: nil,
	}
	values, err := confluence.BuildQueryForTest(params)
	require.NoError(t, err)
	require.Equal(t, url.Values{}, values)
}

func TestV2ClientDoJSON_EncodeError(t *testing.T) {
	client, err := confluence.NewClient("http://example.com")
	require.NoError(t, err)

	bad := make(chan int)
	err = confluence.DoJSONForTest(context.Background(), client, http.MethodPost, "pages", nil, bad, &map[string]any{})
	require.Error(t, err)
	require.ErrorContains(t, err, "encode request")
}

func TestV2ClientDoJSON_QueryEncodingError(t *testing.T) {
	params := map[string]any{"bad": json.RawMessage("{")}
	_, err := confluence.BuildQueryForTest(params)
	require.Error(t, err)
	require.ErrorContains(t, err, "struct, map, or url.Values")
}
