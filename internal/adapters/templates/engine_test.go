package templates

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"rendering/internal/app"
)

func TestLoadMissingDir(t *testing.T) {
	t.Parallel()
	_, err := Load(t.TempDir())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRenderKnownAndUnknown(t *testing.T) {
	t.Parallel()
	eng := loadTemplates(t)
	ctx := context.Background()
	page, err := app.NewPages().Index(ctx)
	if err != nil {
		t.Fatal(err)
	}
	body, err := eng.Render(ctx, page.Name, page.View)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "Rendering") {
		t.Fatalf("index body %s", body)
	}
	if _, err := eng.Render(ctx, "missing", nil); err == nil {
		t.Fatal("expected unknown define to fail")
	}
}

func TestStatusDefinesContainCode(t *testing.T) {
	t.Parallel()
	eng := loadTemplates(t)
	ctx := context.Background()
	pages := app.NewPages()
	for _, code := range []int{400, 401, 403, 404, 405, 500, 502, 503} {
		page, err := pages.Status(ctx, code)
		if err != nil {
			t.Fatal(err)
		}
		body, err := eng.Render(ctx, page.Name, page.View)
		if err != nil {
			t.Fatalf("%d: %v", code, err)
		}
		if !strings.Contains(string(body), strconv.Itoa(code)) {
			t.Fatalf("%d body %s", code, body)
		}
	}
}

func loadTemplates(t *testing.T) *Engine {
	t.Helper()
	eng, err := Load("../../../templates")
	if err != nil {
		t.Fatal(err)
	}
	return eng
}
