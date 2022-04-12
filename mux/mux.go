package mux

import (
	"net/http"
	"strings"
)

type RouteEntry struct {
	handler http.HandlerFunc
	route   string
	method  string
}

type HttpHandler struct {
	//Keys will be method:route such as...
	// GET:/items
	routes map[string]*RouteEntry
}

func NewHandler() *HttpHandler {
	return &HttpHandler{
		routes: make(map[string]*RouteEntry),
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

func (mux *HttpHandler) addRouteEntry(r string, h http.HandlerFunc, m string) {
	if len(r) == 0 || h == nil || len(m) == 0 {
		panic("mux: nil handler or empty route or method")
	}
	mux.routes[m+":"+r] = &RouteEntry{
		handler: h,
		route:   r,
		method:  m,
	}
}

func (mux *HttpHandler) findRoute(path string) *RouteEntry {
	// First try to grab the exact route
	if re, ok := mux.routes[path]; ok {
		return re
	}
	// Second find routes with a matching prefix
	for k, re := range mux.routes {
		// consider using a length sorted slice similar to Default mux
		if strings.HasPrefix(path, k) {
			return re
		}
	}
	return nil
}

func (mux *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := mux.findRoute(r.Method + ":" + r.URL.Path)
	if route == nil {
		http.NotFound(w, r)
		return
	}
	route.handler(w, r)
}
