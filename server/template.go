package server

import (
	html "html/template"
	text "text/template"

	"bytes"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"koushoku/cache"
	. "koushoku/config"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type RenderOptions struct {
	Cache  bool
	Data   map[string]any
	Name   string
	Status int
}

const (
	htmlContentType = "text/html; charset=utf-8"
	xmlContentType  = "application/xml; charset=utf-8"
)

var (
	htmlTemplates       *html.Template
	xmlTemplates        *text.Template
	mu                  sync.Mutex
	ErrTemplateNotFound = errors.New("Template not found")
)

func LoadTemplates() {
	mu.Lock()
	defer mu.Unlock()

	var htmlFiles, xmlFiles []string
	err := filepath.Walk(filepath.Join(Config.Directories.Templates),
		func(path string, stat fs.FileInfo, err error) error {
			if err != nil || stat.IsDir() {
				return err
			}

			if strings.HasSuffix(path, ".html") {
				htmlFiles = append(htmlFiles, path)
			} else {
				xmlFiles = append(xmlFiles, path)
			}

			return nil
		})
	if err != nil {
		log.Fatalln(err)
	}

	if len(htmlFiles) > 0 {
		htmlTemplates, err = html.New("").Funcs(helper).ParseFiles(htmlFiles...)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if len(xmlFiles) > 0 {
		xmlTemplates, err = text.New("").Funcs(helper).ParseFiles(xmlFiles...)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func parseTemplate(name string, data any) ([]byte, error) {
	if strings.HasSuffix(name, ".html") {
		return parseHtmlTemplate(name, data)
	}
	return parseXmlTemplate(name, data)
}

func parseHtmlTemplate(name string, data any) ([]byte, error) {
	if gin.Mode() == gin.DebugMode {
		LoadTemplates()
	}

	t := htmlTemplates.Lookup(name)
	if t == nil {
		return nil, ErrTemplateNotFound

	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		log.Println(err)
		return nil, err
	}
	return buf.Bytes(), nil
}

func parseXmlTemplate(name string, data any) ([]byte, error) {
	if gin.Mode() == gin.DebugMode {
		LoadTemplates()
	}

	t := xmlTemplates.Lookup(name)
	if t == nil {
		return nil, ErrTemplateNotFound

	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		log.Println(err)
		return nil, err
	}
	return buf.Bytes(), nil
}

func getTemplate(name, key string) ([]byte, bool) {
	var v any
	var err error

	if len(key) > 0 {
		v, err = cache.Templates.Get(fmt.Sprintf("%s:%s", name, key))
	} else {
		v, err = cache.Templates.Get(name)
	}

	if err != nil {
		return nil, false
	}
	return v.([]byte), true
}

func setTemplate(name, key string, data any) ([]byte, error) {
	buf, err := parseTemplate(name, data)
	if err != nil {
		return nil, err
	}

	if len(key) > 0 {
		cache.Templates.Set(fmt.Sprintf("%s:%s", name, key), buf, 0)
	} else {
		cache.Templates.Set(name, buf, 0)
	}
	return buf, nil
}

func renderTemplate(c *Context, opts *RenderOptions) {
	var buf []byte
	if opts.Cache {
		var ok bool
		if buf, ok = getTemplate(opts.Name, c.cacheKey()); !ok {
			var err error
			buf, err = setTemplate(opts.Name, c.cacheKey(), opts.Data)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return
			}
		}
	} else {
		buf, _ = parseTemplate(opts.Name, opts.Data)
	}

	contentType := htmlContentType
	if strings.HasSuffix(opts.Name, ".xml") {
		contentType = xmlContentType
	}
	c.Data(opts.Status, contentType, buf)
}
