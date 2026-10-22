package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"rendering/internal/adapters/templates"
	"rendering/internal/app"
	"rendering/internal/domain"

	"github.com/gin-gonic/gin"
)

type fakePages struct{}

func (fakePages) Index(context.Context) (domain.Page, error) {
	return domain.Page{Name: "index", Status: http.StatusOK, View: domain.IndexView()}, nil
}

func (fakePages) Status(_ context.Context, code int) (domain.Page, error) {
	return domain.Page{Name: "500", Status: code}, nil
}

type fakeRenderer struct{}

func (fakeRenderer) Render(context.Context, string, any) ([]byte, error) {
	return []byte("<html>ok</html>"), nil
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func testRouter() *gin.Engine {
	router := gin.New()
	Register(router, NewHandler(fakePages{}, fakeRenderer{}))
	return router
}

func TestHealthz(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	testRouter().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "ok" {
		t.Fatalf("body %s", rec.Body.String())
	}
}

func TestPageSecurityHeaders(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	testRouter().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	headers := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"Referrer-Policy":        "strict-origin-when-cross-origin",
		"X-Frame-Options":        "DENY",
	}
	for key, want := range headers {
		if got := rec.Header().Get(key); got != want {
			t.Fatalf("%s = %q", key, got)
		}
	}
	csp := rec.Header().Get("Content-Security-Policy")
	if csp == "" || !strings.Contains(csp, "script-src 'self'") || !strings.Contains(csp, "style-src 'self'") {
		t.Fatalf("csp %q", csp)
	}
}

func TestIndexUsesDefaultTitle(t *testing.T) {
	t.Parallel()
	rec := serve(t, http.MethodGet, "/")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "Rendering") {
		t.Fatalf("status %d body %s", rec.Code, rec.Body.String())
	}
}

func TestUnknownPathUses404Template(t *testing.T) {
	t.Parallel()
	rec := serve(t, http.MethodGet, "/missing")
	body := rec.Body.String()
	if rec.Code != http.StatusNotFound || !strings.Contains(body, "404") || !strings.Contains(body, "Not Found") {
		t.Fatalf("status %d body %s", rec.Code, body)
	}
	if strings.Contains(body, "Page not found") {
		t.Fatalf("inline 404 body %s", body)
	}
}

func TestWrongMethodUses405Template(t *testing.T) {
	t.Parallel()
	rec := serve(t, http.MethodPost, "/")
	if rec.Code != http.StatusMethodNotAllowed || !strings.Contains(rec.Body.String(), "405") {
		t.Fatalf("status %d body %s", rec.Code, rec.Body.String())
	}
}

func TestRendererFailureWrites500(t *testing.T) {
	t.Parallel()
	router := gin.New()
	Register(router, NewHandler(app.NewPages(), failRenderer{}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusInternalServerError || !strings.Contains(rec.Body.String(), "template error") {
		t.Fatalf("status %d body %s", rec.Code, rec.Body.String())
	}
}

type failRenderer struct{}

func (failRenderer) Render(context.Context, string, any) ([]byte, error) {
	return nil, errRender
}

type renderError struct{}

func (renderError) Error() string { return "render failed" }

var errRender = renderError{}

func serve(t *testing.T, method, path string) *httptest.ResponseRecorder {
	t.Helper()
	eng, err := templates.Load("../../../templates")
	if err != nil {
		t.Fatal(err)
	}
	router := gin.New()
	Register(router, NewHandler(app.NewPages(), eng))
	req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}
