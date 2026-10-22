package http

import (
	stdhttp "net/http"

	"rendering/internal/domain"
	"rendering/internal/ports"

	"github.com/gin-gonic/gin"
)

// Handler translates HTTP calls into the pages port and writes rendered bytes.
type Handler struct {
	pages    ports.Pages
	renderer ports.Renderer
}

// NewHandler returns an HTTP adapter.
func NewHandler(pages ports.Pages, renderer ports.Renderer) *Handler {
	return &Handler{pages: pages, renderer: renderer}
}

const csp = "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self'; font-src 'self'; connect-src 'self'; object-src 'none'; base-uri 'self'; form-action 'self'; frame-ancestors 'none'"

// SecurityHeaders sets the response headers required on every request.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Writer.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Content-Security-Policy", csp)
		c.Next()
	}
}

// Register mounts page routes on router. It does not replace the current server.
func Register(router *gin.Engine, h *Handler) {
	router.HandleMethodNotAllowed = true
	router.Use(SecurityHeaders())
	router.Use(gin.CustomRecovery(func(c *gin.Context, _ any) {
		h.writeStatus(c, stdhttp.StatusInternalServerError)
	}))
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(stdhttp.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/", h.index)
	router.GET("/index.html", h.index)
	router.NoRoute(func(c *gin.Context) {
		h.writeStatus(c, stdhttp.StatusNotFound)
	})
	router.NoMethod(func(c *gin.Context) {
		h.writeStatus(c, stdhttp.StatusMethodNotAllowed)
	})
}

func (h *Handler) index(c *gin.Context) {
	page, err := h.pages.Index(c.Request.Context())
	if err != nil {
		h.writeStatus(c, stdhttp.StatusInternalServerError)
		return
	}
	h.write(c, page)
}

func (h *Handler) writeStatus(c *gin.Context, code int) {
	page, err := h.pages.Status(c.Request.Context(), code)
	if err != nil {
		c.String(stdhttp.StatusInternalServerError, "template error")
		return
	}
	h.write(c, page)
}

func (h *Handler) write(c *gin.Context, page domain.Page) {
	body, err := h.renderer.Render(c.Request.Context(), page.Name, page.View)
	if err != nil {
		if page.Status == stdhttp.StatusInternalServerError {
			c.String(stdhttp.StatusInternalServerError, "template error")
			return
		}
		h.writeStatus(c, stdhttp.StatusInternalServerError)
		return
	}
	c.Data(page.Status, "text/html; charset=utf-8", body)
}
