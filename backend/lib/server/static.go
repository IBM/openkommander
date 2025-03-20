package server

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed static/index.html static/assets/*.css static/assets/*.js
var staticFiles embed.FS

func RegisterStaticRoutes(router *gin.Engine) {
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}

	serveFile := func(c *gin.Context, filePath string) {
		content, err := fs.ReadFile(staticFS, filePath)
		if err != nil {
			log.Printf("Error reading file %s: %v", filePath, err)
			c.String(http.StatusNotFound, "File not found")
			return
		}

		ext := filepath.Ext(filePath)
		switch ext {
		case ".html":
			c.Header("Content-Type", "text/html")
		case ".css":
			c.Header("Content-Type", "text/css")
		case ".js":
			c.Header("Content-Type", "application/javascript")
		}
		
		c.String(http.StatusOK, string(content))
	}

	router.GET("/", func(c *gin.Context) {
		serveFile(c, "index.html")
	})
	
	router.GET("/index.html", func(c *gin.Context) {
		serveFile(c, "index.html")
	})
	
	router.GET("/assets/*filepath", func(c *gin.Context) {
		path := c.Param("filepath")
		if path != "" && path[0] == '/' {
			path = path[1:] // remove leading slash
		}
		serveFile(c, "assets/"+path)
	})
	
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		
		if strings.HasPrefix(path, "/api/") {
			c.Next()
			return
		}
		
		trimmedPath := strings.TrimPrefix(path, "/")
		if _, err := fs.Stat(staticFS, trimmedPath); err != nil {
			// this is important for SPA routing
			serveFile(c, "index.html")
			return
		}
		
		serveFile(c, trimmedPath)
	})
}
