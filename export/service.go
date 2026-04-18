package export

import (
	"context"
	"io"

	"github.com/umats/go-confluence/internal/transport"
)

// Service provides PDF export operations.
type Service struct {
	client *transport.Client
}

// NewService creates a new export service.
func NewService(client *transport.Client) *Service {
	return &Service{client: client}
}

// Page exports a page and returns the PDF bytes.
func (s *Service) Page(ctx context.Context, pageID string) ([]byte, error) {
	return newExporter(s.client).Page(ctx, pageID)
}

// PageTo exports a page and streams the PDF to the writer.
func (s *Service) PageTo(ctx context.Context, pageID string, writer io.Writer) error {
	return newExporter(s.client).PageTo(ctx, pageID, writer)
}
