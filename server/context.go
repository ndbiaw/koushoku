package server

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	. "koushoku/config"

	"github.com/gin-gonic/gin"
)

type Context struct {
	*gin.Context

	sync.RWMutex
	MapData map[string]any
}

func (c *Context) GetURL() string {
	u, _ := url.Parse(Config.Meta.BaseURL)
	u.Path = c.Request.URL.Path
	u.RawQuery = c.Request.URL.RawQuery
	return u.String()
}

func (c *Context) preHTML(code *int) {
	if err, ok := c.GetData("error"); ok {
		err := strings.ToLower(err.(error).Error())
		if strings.Contains(err, "does not exist") || strings.Contains(err, "not found") {
			*code = http.StatusNotFound
		}
	}

	c.SetData("status", *code)
	c.SetData("statusText", http.StatusText(*code))

	if v, ok := c.MapData["name"]; !ok || len(v.(string)) == 0 {
		c.SetData("name", http.StatusText(*code))
	}

	c.SetData("title", Config.Meta.Title)
	c.SetData("description", Config.Meta.Description)
	c.SetData("baseURL", Config.Meta.BaseURL)
	c.SetData("dataBaseURL", Config.Meta.DataBaseURL)
	c.SetData("language", Config.Meta.Language)
	c.SetData("url", c.GetURL())
	c.SetData("query", c.Request.URL.Query())
}

func (c *Context) HTML(code int, name string) {
	c.preHTML(&code)
	renderTemplate(c, &RenderOptions{
		Data:   c.MapData,
		Name:   name,
		Status: code,
	})
}

func (c *Context) Cache(code int, name string) {
	if gin.Mode() == gin.DebugMode {
		c.HTML(code, name)
	} else {
		c.preHTML(&code)
		renderTemplate(c, &RenderOptions{
			Cache:  true,
			Data:   c.MapData,
			Name:   name,
			Status: code,
		})
	}
}

func (c *Context) cacheKey() string {
	return c.GetURL()
}

func (c *Context) IsCached(name string) bool {
	_, ok := getTemplate(name, c.cacheKey())
	return ok
}

func (c *Context) TryCache(name string) bool {
	if c.IsCached(name) {
		c.Cache(http.StatusOK, name)
		return true
	}
	return false
}

func (c *Context) ErrorJSON(code int, message string, err error) {
	c.JSON(code, gin.H{
		"error": gin.H{
			"message": message,
			"cause":   err.Error(),
		},
	})
}

func (c *Context) GetData(key string) (any, bool) {
	c.RLock()
	defer c.RUnlock()

	v, exists := c.MapData[key]
	return v, exists
}

func (c *Context) SetData(key string, value any) {
	c.Lock()
	defer c.Unlock()

	if c.MapData == nil {
		c.MapData = make(map[string]any)
	}
	c.MapData[key] = value
}

func (c *Context) SetCookie(name, value string, expires *time.Time) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Secure:   c.Request.TLS != nil || strings.HasPrefix(Config.Meta.BaseURL, "https"),
		HttpOnly: true,
	}

	if expires == nil {
		cookie.MaxAge = -1
	} else {
		cookie.Expires = *expires
	}
	http.SetCookie(c.Writer, cookie)
}

func (c *Context) ParamInt(name string) (int, error) {
	return strconv.Atoi(c.Param(name))
}

func (c *Context) ParamInt64(name string) (int64, error) {
	return strconv.ParseInt(c.Param(name), 10, 64)
}

type readLimiter struct {
	io.ReadSeeker
	r io.Reader
}

func (r readLimiter) Read(p []byte) (int, error) {
	return r.r.Read(p)
}

const (
	rate     = 1 << 20
	capacity = 1 << 20
)

func (c *Context) serveContent(stream bool, stat os.FileInfo, content io.ReadSeeker) {
	if stream {
		c.Writer.Header().Set("Content-Type", "application/octet-stream")
		c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", stat.Name()))
	}

	//bucket := ratelimit.NewBucketWithRate(rate, capacity)
	//limiter := readLimiter{content, ratelimit.Reader(content, bucket)}
	//http.ServeContent(c.Writer, c.Request, stat.Name(), stat.ModTime(), limiter)
	http.ServeContent(c.Writer, c.Request, stat.Name(), stat.ModTime(), content)
}

func (c *Context) serveFile(stream bool, filepath string) {
	stat, err := os.Stat(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			c.Status(http.StatusNotFound)
		} else {
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	f, err := os.Open(filepath)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	c.serveContent(stream, stat, f)
}

func (c *Context) ServeFile(filepath string) {
	c.serveFile(false, filepath)
}

func (c *Context) StreamFile(filepath string) {
	c.serveFile(true, filepath)
}

func (c *Context) ServeData(stat os.FileInfo, data io.ReadSeeker) {
	c.serveContent(false, stat, data)
}

func (c *Context) StreamData(stat os.FileInfo, data io.ReadSeeker) {
	c.serveContent(true, stat, data)
}
