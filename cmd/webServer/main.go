package main

import (
	"net/http"
	"path/filepath"

	. "koushoku/config"

	"koushoku/cache"
	"koushoku/database"
	"koushoku/server"
	"koushoku/services"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Init()
	cache.Init()

	if err := services.AnalyzeStats(); err != nil {
		return
	}
	server.Init()

	assets := server.Group("/")
	assets.Use(func(c *gin.Context) {
		c.Header("Cache-Control", "public, max-age=300")
	})

	assets.Static("/js", filepath.Join(Config.Directories.Root, "assets/js"))
	assets.Static("/css", filepath.Join(Config.Directories.Root, "assets/css"))
	assets.Static("/fonts", filepath.Join(Config.Directories.Root, "assets/fonts"))

	assets.StaticFile("/cover.jpg", filepath.Join(Config.Directories.Root, "cover.jpg"))
	assets.StaticFile("/robots.txt", filepath.Join(Config.Directories.Root, "robots.txt"))
	assets.StaticFile("/updates.txt", filepath.Join(Config.Directories.Root, "updates.txt"))

	assets.StaticFile("/favicon.ico", filepath.Join(Config.Directories.Root, "favicon.ico"))
	assets.StaticFile("/favicon-16x16.png", filepath.Join(Config.Directories.Root, "favicon-16x16.png"))
	assets.StaticFile("/favicon-32x32.png", filepath.Join(Config.Directories.Root, "favicon-32x32.png"))
	assets.StaticFile("/apple-touch-icon.png", filepath.Join(Config.Directories.Root, "apple-touch-icon.png"))
	assets.StaticFile("/android-chrome-192x192.png", filepath.Join(Config.Directories.Root, "android-chrome-192x192.png"))
	assets.StaticFile("/android-chrome-512x512.png", filepath.Join(Config.Directories.Root, "android-chrome-512x512.png"))

	server.GET("/", index)
	server.GET("/about", server.WithName("About"), about)
	server.GET("/search", search)
	server.GET("/stats", server.WithName("Stats"), stats)
	server.GET("/sitemap.xml", sitemap)

	server.GET("/archive/:id", archive)
	server.GET("/archive/:id/:slug", archive)
	server.GET("/archive/:id/:slug/:pageNum", read)
	server.GET("/artists", artists)
	server.GET("/artists/:slug", artist)
	server.GET("/circles", circles)
	server.GET("/circles/:slug", circle)
	server.GET("/magazines", magazines)
	server.GET("/magazines/:slug", magazine)
	server.GET("/parodies", parodies)
	server.GET("/parodies/:slug", parody)
	server.GET("/tags", tags)
	server.GET("/tags/:slug", tag)

	server.GET("/submit", server.WithName("Submit"), submit)
	server.POST("/submit", server.WithName("Submit"), server.WithRateLimit("Submit?", "10-D"), submitPost)
	server.GET("/submissions", submisisions)

	server.POST("/api/purge-cache", purgeCache)
	server.POST("/api/reload-templates", reloadTemplates)

	server.NoRoute(func(c *server.Context) {
		c.HTML(http.StatusNotFound, "error.html")
	})

	server.Start(Config.Server.WebPort)
}
