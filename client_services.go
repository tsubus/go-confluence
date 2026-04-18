package confluence

import (
	"github.com/umats/go-confluence/attachment"
	"github.com/umats/go-confluence/export"
	"github.com/umats/go-confluence/internal/transport"
	"github.com/umats/go-confluence/label"
	"github.com/umats/go-confluence/page"
	"github.com/umats/go-confluence/space"
)

// Page returns the PageService for REST v2 endpoints.
func (c *Client) Page() PageService {
	return page.NewService(newTransportClient(c))
}

// Space returns the SpaceService for REST v2 endpoints.
func (c *Client) Space() SpaceService {
	return space.NewService(newTransportClient(c))
}

// Attachment returns the AttachmentService for REST v2 endpoints.
func (c *Client) Attachment() AttachmentService {
	return attachment.NewService(newTransportClient(c))
}

// Label returns the LabelService for REST v2 endpoints.
func (c *Client) Label() LabelService {
	return label.NewService(newTransportClient(c))
}

// Export returns the ExportService for PDF exports.
func (c *Client) Export() ExportService {
	return export.NewService(newTransportClient(c))
}

var _ PageService = (*page.Service)(nil)
var _ SpaceService = (*space.Service)(nil)
var _ AttachmentService = (*attachment.Service)(nil)
var _ LabelService = (*label.Service)(nil)
var _ ExportService = (*export.Service)(nil)

func newTransportClient(client *Client) *transport.Client {
	return &transport.Client{
		BaseURL:                  client.baseURL,
		HTTPClient:               client.httpClient,
		Username:                 client.username,
		Password:                 client.password,
		PollInterval:             client.pollInterval,
		PollTimeout:              client.pollTimeout,
		AllowedRedirectHosts:     client.allowedRedirectHosts,
		AllowCrossHostContentURL: client.allowCrossHostContentURL,
	}
}
