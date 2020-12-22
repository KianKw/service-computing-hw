package resthttp

import (
	"context"
	"net/http"
	"strings"
	"sync"
)

type Handle func(http.ResponseWriter, *http.Request, Params)

type Param struct {
	Key   string
	Value string
}

type Params []Param

func (ps Params) ByName(name string) string {
	for _, p := range ps {
		if p.Key == name {
			return p.Value
		}
	}
	return ""
}

type paramsKey struct{}

var ParamsKey = paramsKey{}

func ParamsFromContext(ctx context.Context) Params {
	p, _ := ctx.Value(ParamsKey).(Params)
	return p
}

var MatchedRoutePathParam = "$matchedRoutePath"

func (ps Params) MatchedRoutePath() string {
	return ps.ByName(MatchedRoutePathParam)
}

type Api struct {
	router *Router
}

func NewApi() *Api {
	return &Api{
		router: New(),
	}
}

type Middleware struct {
	method string
	path   string
	handle Handle
}

func (api *Api) SetRouter(middlewares ...*Middleware) {
	for _, middleware := range middlewares {
		switch middleware.method {
		case "GET":
			api.router.GET(middleware.path, middleware.handle)
		case "POST":
			api.router.POST(middleware.path, middleware.handle)
		case "PUT":
			api.router.PUT(middleware.path, middleware.handle)
		case "DELETE":
			api.router.DELETE(middleware.path, middleware.handle)
		}
	}
}

func GET(p string, h Handle) *Middleware {
	return &Middleware{
		method: "GET",
		path:   p,
		handle: h,
	}
}

func POST(p string, h Handle) *Middleware {
	return &Middleware{
		method: "POST",
		path:   p,
		handle: h,
	}
}

func PUT(p string, h Handle) *Middleware {
	return &Middleware{
		method: "PUT",
		path:   p,
		handle: h,
	}
}

func DELETE(p string, h Handle) *Middleware {
	return &Middleware{
		method: "DELETE",
		path:   p,
		handle: h,
	}
}
func (api *Api) MakeHandler() http.Handler {
	return api.router
}

type Router struct {
	tires                  map[string]*node
	pPools                 sync.Pool
	maxParams              uint16
	isStoreThePath         bool
	RedirectTrailingSlash  bool
	isRedirectTheUnchangeP bool
	isAllowMethod          bool
	isAllWork              string
	NotFound               http.Handler
	MethodNotAllowed       http.Handler
	errorHandle            func(http.ResponseWriter, *http.Request, interface{})
}

var _ http.Handler = New()

func New() *Router {
	return &Router{
		RedirectTrailingSlash:  true,
		isRedirectTheUnchangeP: true,
		isAllowMethod:          true,
	}
}

func (router *Router) getParams() *Params {
	ps, _ := router.pPools.Get().(*Params)
	*ps = (*ps)[0:0] // reset slice
	return ps
}

func (router *Router) putParams(ps *Params) {
	if ps != nil {
		router.pPools.Put(ps)
	}
}

func (router *Router) saveMatchedRoutePath(path string, handle Handle) Handle {
	return func(w http.ResponseWriter, request *http.Request, ps Params) {
		if ps == nil {
			psp := router.getParams()
			ps = (*psp)[0:1]
			ps[0] = Param{Key: MatchedRoutePathParam, Value: path}
			handle(w, request, ps)
			router.putParams(psp)
		} else {
			ps = append(ps, Param{Key: MatchedRoutePathParam, Value: path})
			handle(w, request, ps)
		}
	}
}

// GET is a shortcut for router.Handle(http.MethodGet, path, handle)
func (router *Router) GET(path string, handle Handle) {
	router.Handle(http.MethodGet, path, handle)
}

// HEAD is a shortcut for router.Handle(http.MethodHead, path, handle)
func (router *Router) HEAD(path string, handle Handle) {
	router.Handle(http.MethodHead, path, handle)
}

// OPTIONS is a shortcut for router.Handle(http.MethodOptions, path, handle)
func (router *Router) OPTIONS(path string, handle Handle) {
	router.Handle(http.MethodOptions, path, handle)
}

// POST is a shortcut for router.Handle(http.MethodPost, path, handle)
func (router *Router) POST(path string, handle Handle) {
	router.Handle(http.MethodPost, path, handle)
}

// PUT is a shortcut for router.Handle(http.MethodPut, path, handle)
func (router *Router) PUT(path string, handle Handle) {
	router.Handle(http.MethodPut, path, handle)
}

// PATCH is a shortcut for router.Handle(http.MethodPatch, path, handle)
func (router *Router) PATCH(path string, handle Handle) {
	router.Handle(http.MethodPatch, path, handle)
}

// DELETE is a shortcut for router.Handle(http.MethodDelete, path, handle)
func (router *Router) DELETE(path string, handle Handle) {
	router.Handle(http.MethodDelete, path, handle)
}

func cleanPath(p string) string {
	if p == "" {
		return "/"
	}

	var b strings.Builder
	b.Grow(len(p))

	if p[0] != '/' {
		b.Write([]byte{'/'})
	}

	index := strings.Index(p, "//")
	if index == -1 {
		_, err := b.WriteString(p)
		if err != nil {
			panic(err)
		}
		return b.String()
	}

	b.WriteString(p[:index+1])

	slash := true
	for i := index + 2; i < len(p); i++ {
		if p[i] == '/' {
			if slash {
				continue
			}
			slash = true
		} else {
			slash = false
		}
		b.Write([]byte{p[i]})
	}
	return b.String()
}

func (router *Router) Handle(m, path string, handle Handle) {
	varsNum := uint16(0)

	if router.isStoreThePath {
		varsNum++
		handle = router.saveMatchedRoutePath(path, handle)
	}

	if router.tires == nil {
		router.tires = make(map[string]*node)
	}

	root := router.tires[m]
	if root == nil {
		root = new(node)
		router.tires[m] = root

		router.isAllWork = router.allowed("*", "")
	}

	root.addRoute(path, handle)

	// Update maxParams
	if paramsNum := ParamNum(path); paramsNum+varsNum > router.maxParams {
		router.maxParams = paramsNum + varsNum
	}

	// Lazy-init pPools alloc func
	if router.pPools.New == nil && router.maxParams > 0 {
		router.pPools.New = func() interface{} {
			ps := make(Params, 0, router.maxParams)
			return &ps
		}
	}
}

// Handler is an adapter which allows the usage of an http.Handler as a
// request handle.
// The Params are available in the request context under ParamsKey.
func (router *Router) Handler(m, path string, handler http.Handler) {
	router.Handle(m, path,
		func(w http.ResponseWriter, request *http.Request, p Params) {
			if len(p) > 0 {
				ctx := request.Context()
				ctx = context.WithValue(ctx, ParamsKey, p)
				request = request.WithContext(ctx)
			}
			handler.ServeHTTP(w, request)
		},
	)
}

// HandlerFunc is an adapter which allows the usage of an http.HandlerFunc as a
// request handle.
func (router *Router) HandlerFunc(m, path string, handler http.HandlerFunc) {
	router.Handler(m, path, handler)
}

func (router *Router) ServeFiles(path string, root http.FileSystem) {
	if len(path) < 10 || path[len(path)-10:] != "/*filepath" {
		panic("path must end with /*filepath in path '" + path + "'")
	}

	fileServer := http.FileServer(root)

	router.GET(path, func(w http.ResponseWriter, request *http.Request, ps Params) {
		request.URL.Path = ps.ByName("filepath")
		fileServer.ServeHTTP(w, request)
	})
}

func (router *Router) myRecover(w http.ResponseWriter, request *http.Request) {
	if rcv := recover(); rcv != nil {
		router.errorHandle(w, request, rcv)
	}
}

func (router *Router) Lookup(m, path string) (Handle, Params, bool) {
	if root := router.tires[m]; root != nil {
		handle, ps, tsr := root.getValue(path, router.getParams)
		if handle == nil {
			router.putParams(ps)
			return nil, nil, tsr
		}
		if ps == nil {
			return handle, nil, tsr
		}
		return handle, *ps, tsr
	}
	return nil, nil, false
}

func (router *Router) allowed(path, reqMethod string) (allow string) {
	allowed := make([]string, 0, 9)

	if path == "*" { // server-wide
		// empty method is used for internal calls to refresh the cache
		if reqMethod == "" {
			for m := range router.tires {
				if m == http.MethodOptions {
					continue
				}
				// Add request method to list of allowed methods
				allowed = append(allowed, m)
			}
		} else {
			return router.isAllWork
		}
	} else { // specific path
		for m := range router.tires {
			// Skip the requested m - we already tried this one
			if m == reqMethod || m == http.MethodOptions {
				continue
			}

			handle, _, _ := router.tires[m].getValue(path, nil)
			if handle != nil {
				// Add request m to list of allowed methods
				allowed = append(allowed, m)
			}
		}
	}

	if len(allowed) > 0 {
		// Add request m to list of allowed methods
		allowed = append(allowed, http.MethodOptions)

		for i, l := 1, len(allowed); i < l; i++ {
			for j := i; j > 0 && allowed[j] < allowed[j-1]; j-- {
				allowed[j], allowed[j-1] = allowed[j-1], allowed[j]
			}
		}

		// return as comma separated list
		return strings.Join(allowed, ", ")
	}

	return allow
}

// ServeHTTP makes the router implement the http.Handler interface.
func (router *Router) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	if router.errorHandle != nil {
		defer router.myRecover(w, request)
	}

	reqPath := request.URL.Path

	if root := router.tires[request.Method]; root != nil {
		if handle, ps, tsr := root.getValue(reqPath, router.getParams); handle != nil {
			if ps != nil {
				handle(w, request, *ps)
				router.putParams(ps)
			} else {
				handle(w, request, nil)
			}
			return
		} else if request.Method != http.MethodConnect && reqPath != "/" {
			// Moved Permanently, request with GET method
			code := http.StatusMovedPermanently
			if request.Method != http.MethodGet {
				// Permanent Redirect, request with same method
				code = http.StatusPermanentRedirect
			}

			if tsr && router.RedirectTrailingSlash {
				if len(reqPath) > 1 && reqPath[len(reqPath)-1] == '/' {
					request.URL.Path = reqPath[:len(reqPath)-1]
				} else {
					request.URL.Path = reqPath + "/"
				}
				http.Redirect(w, request, request.URL.String(), code)
				return
			}

			// Try to fix the request reqPath
			if router.isRedirectTheUnchangeP {
				fixedPath, found := root.findCaseInsensitivePath(
					cleanPath(reqPath),
					router.RedirectTrailingSlash,
				)
				if found {
					request.URL.Path = fixedPath
					http.Redirect(w, request, request.URL.String(), code)
					return
				}
			}
		}
	}

	if router.isAllowMethod { // Handle 405
		if allow := router.allowed(reqPath, request.Method); allow != "" {
			w.Header().Set("Allow", allow)
			if router.MethodNotAllowed != nil {
				router.MethodNotAllowed.ServeHTTP(w, request)
			} else {
				http.Error(w,
					http.StatusText(http.StatusMethodNotAllowed),
					http.StatusMethodNotAllowed,
				)
			}
			return
		}
	}

	// Handle 404
	if router.NotFound != nil {
		router.NotFound.ServeHTTP(w, request)
	} else {
		http.NotFound(w, request)
	}
}
