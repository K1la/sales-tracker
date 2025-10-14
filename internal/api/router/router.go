package router

import (
	"github.com/K1la/sales-tracker/internal/api/handlers/analytics"
	"github.com/K1la/sales-tracker/internal/api/handlers/items"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/wb-go/wbf/ginext"
)

func New(ih *items.Handler, ah *analytics.Handler) *ginext.Engine {
	e := ginext.New("")
	e.Use(ginext.Recovery(), ginext.Logger())

	// API routes
	api := e.Group("/api")
	{
		itemsGroup := api.Group("/items")
		{
			itemsGroup.POST("", ih.Create)
			itemsGroup.GET("", ih.GetAll)
			itemsGroup.GET("/:id", ih.GetByID)
			itemsGroup.PUT("/:id", ih.Update)
			itemsGroup.DELETE("/:id", ih.Delete)
		}

		analyticsGroup := api.Group("/analytics")
		{
			analyticsGroup.GET("", ah.GetSummary)
		}
	}

	// Frontend: serve files from ./web without conflicting wildcard
	e.NoRoute(func(c *ginext.Context) {
		if c.Request.URL.Path == "/" {
			http.ServeFile(c.Writer, c.Request, "./web/index.html")
			return
		}
		// Serve only files under /web/ directly from disk
		if strings.HasPrefix(c.Request.URL.Path, "/web/") {
			safe := filepath.Clean("." + c.Request.URL.Path)
			http.ServeFile(c.Writer, c.Request, safe)
			return
		}
		c.Status(http.StatusNotFound)
	})

	return e
}
