package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	. "koushoku/config"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type Handler func(c *Context)
type Handlers []Handler

var server *gin.Engine
var secretHandler func()

func Init() {
	if strings.EqualFold(Config.Mode, "production") {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	server = gin.Default()
	LoadTemplates()

	server.ForwardedByClientIP = true
	server.RedirectTrailingSlash = true
	server.RemoveExtraSlash = true

	server.Use(gzip.Gzip(gzip.DefaultCompression))
	if secretHandler != nil {
		secretHandler()
	}
}

func Start(port int) {
	if gin.Mode() != gin.DebugMode {
		log.Println("Listening and serving HTTP on :", port)
	}

	srv := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: server}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}

func (h Handler) wrap() gin.HandlerFunc {
	return func(c *gin.Context) {
		var context *Context
		if v, exists := c.Get("context"); exists {
			context = v.(*Context)
		} else {
			context = &Context{Context: c}
			c.Set("context", context)
		}
		h(context)
	}
}

func (h Handlers) wrap() []gin.HandlerFunc {
	ginHandlers := make([]gin.HandlerFunc, len(h))
	for i := range h {
		ginHandlers[i] = h[i].wrap()
	}
	return ginHandlers
}

func Group(relativePath string, handlers ...Handler) *gin.RouterGroup {
	return server.Group(relativePath, Handlers(handlers).wrap()...)
}

func Handle(method string, relativePath string, handlers ...Handler) {
	server.Handle(method, relativePath, Handlers(handlers).wrap()...)
}

func NoRoute(handlers ...Handler) {
	server.NoRoute(Handlers(handlers).wrap()...)
}

func GET(relativePath string, handlers ...Handler) {
	server.GET(relativePath, Handlers(handlers).wrap()...)
}

func POST(relativePath string, handlers ...Handler) {
	server.POST(relativePath, Handlers(handlers).wrap()...)
}

func PATCH(relativePath string, handlers ...Handler) {
	server.PATCH(relativePath, Handlers(handlers).wrap()...)
}

func PUT(relativePath string, handlers ...Handler) {
	server.PUT(relativePath, Handlers(handlers).wrap()...)
}

func DELETE(relativePath string, handlers ...Handler) {
	server.DELETE(relativePath, Handlers(handlers).wrap()...)
}

func HEAD(relativePath string, handlers ...Handler) {
	server.HEAD(relativePath, Handlers(handlers).wrap()...)
}
