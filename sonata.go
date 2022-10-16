package sonata

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"text/template"

	"github.com/mneumi/sonata/render"
)

const AnyMethod = "AnyMethod"

type HandleFunc func(ctx *Context)

type MiddlewareFunc func(HandleFunc) HandleFunc

type routerGroup struct {
	name                     string
	handleFuncMap            map[string]map[string]HandleFunc
	middlewaresHandleFuncMap map[string]map[string][]MiddlewareFunc
	middlewares              []MiddlewareFunc
}

func (rg *routerGroup) Use(middlewareFuncs ...MiddlewareFunc) {
	rg.middlewares = append(rg.middlewares, middlewareFuncs...)
}

func (rg *routerGroup) methodHandle(ctx *Context, name string, method string, handle HandleFunc) {
	middlewareFuncs := rg.middlewaresHandleFuncMap[name][method]
	for index := len(middlewareFuncs) - 1; index >= 0; index-- {
		handle = middlewareFuncs[index](handle)
	}

	for index := len(rg.middlewares) - 1; index >= 0; index-- {
		handle = rg.middlewares[index](handle)
	}

	handle(ctx)
}

func (rg *routerGroup) handle(name string, method string, handleFunc HandleFunc, middlewareFunc ...MiddlewareFunc) {
	_, ok := rg.handleFuncMap[name]
	if !ok {
		rg.handleFuncMap[name] = make(map[string]HandleFunc)
		rg.middlewaresHandleFuncMap[name] = make(map[string][]MiddlewareFunc)
	}
	_, ok = rg.handleFuncMap[name][method]
	if ok {
		panic("err: register same route")
	}
	rg.handleFuncMap[name][method] = handleFunc
	rg.middlewaresHandleFuncMap[name][method] = append(rg.middlewaresHandleFuncMap[name][method], middlewareFunc...)
}

func (rg *routerGroup) Any(name string, handleFunc HandleFunc, middlewareFunc ...MiddlewareFunc) {
	rg.handle(name, AnyMethod, handleFunc, middlewareFunc...)
}

func (rg *routerGroup) Get(name string, handleFunc HandleFunc, middlewareFunc ...MiddlewareFunc) {
	rg.handle(name, http.MethodGet, handleFunc, middlewareFunc...)
}

func (rg *routerGroup) Post(name string, handleFunc HandleFunc, middlewareFunc ...MiddlewareFunc) {
	rg.handle(name, http.MethodPost, handleFunc, middlewareFunc...)
}

func (rg *routerGroup) Put(name string, handleFunc HandleFunc, middlewareFunc ...MiddlewareFunc) {
	rg.handle(name, http.MethodPut, handleFunc, middlewareFunc...)
}

func (rg *routerGroup) Delete(name string, handleFunc HandleFunc, middlewareFunc ...MiddlewareFunc) {
	rg.handle(name, http.MethodDelete, handleFunc, middlewareFunc...)
}

func (rg *routerGroup) Patch(name string, handleFunc HandleFunc, middlewareFunc ...MiddlewareFunc) {
	rg.handle(name, http.MethodPatch, handleFunc, middlewareFunc...)
}

type router struct {
	routerGroups []*routerGroup
}

func (rg *router) Group(name string) *routerGroup {
	routerGroup := &routerGroup{
		name:                     name,
		handleFuncMap:            make(map[string]map[string]HandleFunc),
		middlewaresHandleFuncMap: make(map[string]map[string][]MiddlewareFunc),
		middlewares:              []MiddlewareFunc{},
	}
	rg.routerGroups = append(rg.routerGroups, routerGroup)
	return routerGroup
}

type Engine struct {
	router
	funcMap    template.FuncMap
	htmlRender render.HTMLRender
	pool       sync.Pool
}

func New() *Engine {
	e := &Engine{
		router: router{},
	}
	e.pool.New = func() any {
		return e.allocateContext()
	}
	return e
}

func (e *Engine) allocateContext() any {
	return &Context{
		engine: e,
	}
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

func (e *Engine) LoadTemplate(pattern string) {
	t := template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
	e.SetHTMLTemplate(t)
}

func (e *Engine) SetHTMLTemplate(t *template.Template) {
	e.htmlRender = render.HTMLRender{
		Template: t,
	}
}

func (e *Engine) httpRequestHandle(ctx *Context, w http.ResponseWriter, r *http.Request) {
	method := r.Method

	for _, group := range e.routerGroups {
		for name, handleFuncMap := range group.handleFuncMap {
			url := "/" + group.name + name
			if r.RequestURI == url {
				if handle, ok := handleFuncMap[AnyMethod]; ok {
					group.methodHandle(ctx, name, AnyMethod, handle)
					return
				}

				if handle, ok := handleFuncMap[method]; ok {
					group.methodHandle(ctx, name, method, handle)
					return
				}

				w.WriteHeader(http.StatusMethodNotAllowed)
				fmt.Fprintf(w, "%s %s not allow", r.RequestURI, method)

				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s %s not found", r.RequestURI, method)
	}
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := e.pool.Get().(*Context)
	ctx.W = w
	ctx.R = r
	e.httpRequestHandle(ctx, w, r)
	e.pool.Put(ctx)
}

func (e *Engine) Run() {
	if err := http.ListenAndServe(":8111", e); err != nil {
		log.Fatal(err)
	}
}
