package app

import (
	"context"
	"net/http"
	"testing"

	"rendering/internal/domain"
)

func TestIndex(t *testing.T) {
	t.Parallel()
	page, err := NewPages().Index(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if page.Name != "index" || page.Status != http.StatusOK {
		t.Fatalf("page %+v", page)
	}
	if page.View.Brand != domain.Brand {
		t.Fatalf("brand %q", page.View.Brand)
	}
}

func TestStatusMapping(t *testing.T) {
	t.Parallel()
	pages := NewPages()
	codes := []int{
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusMethodNotAllowed,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
	}
	for _, code := range codes {
		page, err := pages.Status(context.Background(), code)
		if err != nil {
			t.Fatal(err)
		}
		if page.Name != statusNames[code] || page.Status != code {
			t.Fatalf("code %d page %+v", code, page)
		}
	}
	page, err := pages.Status(context.Background(), http.StatusTeapot)
	if err != nil {
		t.Fatal(err)
	}
	if page.Name != "500" || page.Status != http.StatusInternalServerError {
		t.Fatalf("fallback %+v", page)
	}
}
