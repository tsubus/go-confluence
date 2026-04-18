//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	confluence "github.com/umats/go-confluence"
)

func main() {
	baseURL := os.Getenv("CONFLUENCE_URL")
	username := os.Getenv("CONFLUENCE_USERNAME")
	password := os.Getenv("CONFLUENCE_PASSWORD")

	if baseURL == "" || username == "" || password == "" {
		log.Fatal("Missing CONFLUENCE_URL, CONFLUENCE_USERNAME or CONFLUENCE_PASSWORD")
	}

	client, err := confluence.NewClient(baseURL, confluence.WithBasicAuth(username, password))
	if err != nil {
		log.Fatalf("new client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("=== Spaces ===")
	spaces, err := client.Space().List(ctx, nil)
	if err != nil {
		log.Printf("list spaces: %v", err)
	} else {
		for _, s := range spaces.Results {
			fmt.Printf("ID: %s, Key: %s, Name: %s\n", ptrStr(s.Id), ptrStr(s.Key), ptrStr(s.Name))
		}
	}

	fmt.Println("\n=== Pages ===")
	limit := 5
	pages, err := client.Page().List(ctx, &confluence.PageListParams{Limit: &limit})
	if err != nil {
		log.Printf("list pages: %v", err)
	} else {
		for _, p := range pages.Results {
			fmt.Printf("ID: %s, Title: %s\n", ptrStr(p.Id), ptrStr(p.Title))
		}
	}

	fmt.Println("\n=== Attachments ===")
	attachments, err := client.Attachment().List(ctx, &confluence.AttachmentListParams{Limit: &limit})
	if err != nil {
		log.Printf("list attachments: %v", err)
	} else {
		for _, a := range attachments.Results {
			fmt.Printf("ID: %s, Title: %s\n", ptrStr(a.Id), ptrStr(a.Title))
		}
	}

	fmt.Println("\n=== Labels ===")
	labels, err := client.Label().List(ctx, &confluence.LabelListParams{Limit: &limit})
	if err != nil {
		log.Printf("list labels: %v", err)
	} else {
		for _, l := range labels.Results {
			fmt.Printf("ID: %s, Name: %s\n", ptrStr(l.Id), ptrStr(l.Name))
		}
	}

}

func ptrStr(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return *s
}
