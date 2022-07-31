package mux

import (
	"net/http"
	"sort"
	"strings"
	"sync"
)

type RouteEntry struct {
	handler http.HandlerFunc
	route   string
	method  string
}

// HttpHandler As per https://go.dev/blog/maps, map order is
// not guaranteed to be the same. (This is a language decision)
// A sorted slice of string keys in the handler is used to
// match routes from the most to the least specific pattern.
type HttpHandler struct {
	// Keys will be method:route such as: "GET:/items"
	// Could also use a struct key such as {string, string}
	// but keeping it string-based for now.
	routes map[string]RouteEntry

	mutex sync.RWMutex // sync on Write operations (see addRouteEntry)

	sk []string // sorted slice of keys by length of route
}

func NewHandler() *HttpHandler {
	return &HttpHandler{
		routes: make(map[string]RouteEntry),
	}
}

func (mux *HttpHandler) GET(r string, h http.HandlerFunc) {
	mux.addRouteEntry(r, h, "GET")
}

func (mux *HttpHandler) PUT(r string, h http.HandlerFunc) {
	mux.addRouteEntry(r, h, "PUT")
}

func (mux *HttpHandler) POST(r string, h http.HandlerFunc) {
	mux.addRouteEntry(r, h, "POST")
}

func (mux *HttpHandler) DELETE(r string, h http.HandlerFunc) {
	mux.addRouteEntry(r, h, "DELETE")
}

func (mux *HttpHandler) OPTIONS(r string, h http.HandlerFunc) {
	mux.addRouteEntry(r, h, "OPTIONS")
}

func (mux *HttpHandler) HandleFunc(r string, h http.HandlerFunc) {
	mux.addRouteEntry(r, h, "*")
}

func (mux *HttpHandler) addRouteEntry(r string, h http.HandlerFunc, m string) {
	if len(r) == 0 || h == nil || len(m) == 0 {
		panic("mux: nil handler or empty route or method")
	}
	mux.mutex.Lock()
	defer mux.mutex.Unlock()
	key := m + ":" + r
	mux.routes[key] = RouteEntry{
		handler: h,
		route:   r,
		method:  m,
	}
	mux.sk = append(mux.sk, key)
	n := len(mux.sk)
	if n == 1 || len(key) <= len(mux.sk[n-2]) { // already sorted so just return
		return
	}
	sort.Slice(mux.sk, func(i, j int) bool { // readability over speed given a one time hit
		return len(mux.sk[i]) > len(mux.sk[j]) // sort the longest path to shortest
	})
}

// Using a length sorted slice for those familiar with default mux
// and match an exact route first if applicable
func (mux *HttpHandler) findRoute(path string) http.Handler {
	// First try to grab the exact route
	if re, ok := mux.routes[path]; ok {
		return re.handler
	}
	// Second find routes with a matching prefix
	for _, key := range mux.sk {
		if key[0] == '*' { // catchall so replace method in path with *
			_, urlPattern, found := strings.Cut(path, ":")
			if found {
				path = "*:" + urlPattern
			}
		}
		if strings.HasPrefix(path, key) {
			return mux.routes[key].handler
		}
	}
	return nil
}

func (mux *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := mux.findRoute(r.Method + ":" + r.URL.Path)
	if handler == nil {
		http.NotFound(w, r)
		return
	}
	handler.ServeHTTP(w, r)
}
