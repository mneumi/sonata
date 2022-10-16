package sonata

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/mneumi/sonata/render"
)

const defaultMaxMemory = 32 << 20 // 32M

type Context struct {
	W             http.ResponseWriter
	R             *http.Request
	engine        *Engine
	queryCache    url.Values
	postFormCache url.Values
}

func (c *Context) initQueryCache() {
	if c.R != nil {
		c.queryCache = c.R.URL.Query()
	} else {
		c.queryCache = url.Values{}
	}
}

func (c *Context) initPostFormCache() {
	if c.R != nil {
		if err := c.R.ParseMultipartForm(defaultMaxMemory); err != nil {
			if !errors.Is(err, http.ErrNotMultipart) {
				log.Println(err)
			}
		}
		c.postFormCache = c.R.PostForm
	} else {
		c.postFormCache = url.Values{}
	}
}

func (c *Context) GetQuery(key string) string {
	c.initQueryCache()
	return c.queryCache.Get(key)
}

func (c *Context) GetQueryArray(key string) ([]string, bool) {
	c.initQueryCache()
	values, ok := c.queryCache[key]
	return values, ok
}

func (c *Context) DefaultQuery(key string, defaultValue string) string {
	value := c.GetQuery(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (c *Context) GetPostForm(key string) string {
	c.initPostFormCache()
	return c.postFormCache.Get(key)
}

func (c *Context) GetPostFormArray(key string) ([]string, bool) {
	c.initPostFormCache()
	values, ok := c.postFormCache[key]
	return values, ok
}

func (c *Context) DefaultPostForm(key string, defaultValue string) string {
	value := c.GetPostForm(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (c *Context) HTML(status int, html string) error {
	return c.Render(status, &render.HTML{
		Data: html,
	})
}

func (c *Context) HTMLTemplate(name string, data any, filenames ...string) error {
	c.W.Header().Set("Content-Type", "text/html; charset=utf8")
	t := template.New(name)
	t, err := t.ParseFiles(filenames...)
	if err != nil {
		return err
	}
	err = t.Execute(c.W, data)
	return err
}

func (c *Context) HTMLTemplateGlob(name string, data any, pattern string) error {
	c.W.Header().Set("Content-Type", "text/html; charset=utf8")
	t := template.New(name)
	t, err := t.ParseGlob(pattern)
	if err != nil {
		return err
	}
	err = t.Execute(c.W, data)
	return err
}

func (c *Context) Template(name string, data any) error {
	return c.Render(http.StatusOK, &render.HTML{
		Name:       name,
		Data:       data,
		IsTemplate: true,
		Template:   c.engine.htmlRender.Template,
	})
}

func (c *Context) JSON(status int, data any) error {
	return c.Render(status, &render.JSON{
		Data: data,
	})
}

func (c *Context) XML(status int, data any) error {
	return c.Render(status, &render.XML{
		Data: data,
	})
}

func (c *Context) Redirect(status int, location string) error {
	return c.Render(status, &render.Redirect{
		Status:   status,
		Request:  c.R,
		Location: location,
	})
}

func (c *Context) String(status int, format string, values ...any) error {
	return c.Render(status, &render.String{
		Format: format,
		Data:   values,
	})
}

func (c *Context) Render(status int, r render.Render) error {
	r.WriteContentType(c.W)
	r.WriteHeader(status, c.W)
	return r.Render(c.W)
}
