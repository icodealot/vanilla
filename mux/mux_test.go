package mux

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddRoutes(t *testing.T) {
	r := NewHandler()
	r.HandleFunc("/", http.NotFound)
	r.GET("/", http.NotFound)
	r.GET("/a", http.NotFound)
	r.PUT("/bbb", http.NotFound)

	if len(r.sk) != 4 || len(r.routes) != 4 {
		t.Error("Expected 4 routes to be stored in handler")
	}
}

func TestSortedRoutes(t *testing.T) {
	r := NewHandler()
	r.HandleFunc("/", http.NotFound)
	r.GET("/", http.NotFound)
	r.GET("/a", http.NotFound)
	r.GET("/bbb", http.NotFound)

	n := len(r.sk)
	if r.sk[n-1] != "*:/" {
		t.Error("Shortest url pattern should be last")
	}

	if r.sk[0] != "GET:/bbb" {
		t.Error("Longest url pattern should be first")
	}
}

func TestValidRoute(t *testing.T) {
	r := NewHandler()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("Expected HTTP 200 but received %d", resp.StatusCode)
	}
}

func TestInvalidRoute(t *testing.T) {
	r := NewHandler()
	r.HandleFunc("/a", func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest("GET", "/b", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != 404 {
		t.Errorf("Expected HTTP 404 but received %d", resp.StatusCode)
	}
}

func TestExactMatch(t *testing.T) {
	r := NewHandler()
	r.HandleFunc("/a/aa", func(w http.ResponseWriter, r *http.Request) {})
	r.HandleFunc("/a", http.NotFound)

	req := httptest.NewRequest("GET", "/a/aa", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("Expected HTTP 200 but received %d", resp.StatusCode)
	}
}

func TestRedirect(t *testing.T) {
	r := NewHandler()
	r.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {})
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/home", http.StatusFound)
	})

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != 302 {
		t.Errorf("Expected HTTP 302 but received %d", resp.StatusCode)
	}
}

func TestCatchAll(t *testing.T) {
	r := NewHandler()
	r.PUT("/", http.NotFound)
	r.POST("/", http.NotFound)
	r.DELETE("/", http.NotFound)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "success")
	})

	key := r.sk[len(r.sk)-1]
	if key[0] != '*' {
		t.Error("Expected catchall route to have * method prefix")
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "success" {
		t.Errorf("Expected HTTP success but received %s", string(body))
	}

	req = httptest.NewRequest("PUT", "/", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp = w.Result()
	if resp.StatusCode != 404 {
		t.Errorf("Expected HTTP 404 but received %d", resp.StatusCode)
	}
}
