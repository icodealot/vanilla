# vanilla/mux
[![Build Status](https://github.com/icodealot/vanilla/actions/workflows/build.yml/badge.svg)](https://github.com/icodealot/vanilla/actions)
[![GoDoc](https://godoc.org/github.com/icodealot/vanilla/mux?status.svg)](https://godoc.org/github.com/icodealot/vanilla/mux)

An intentionally minimal and easy to use set of HTTP tools for Go projects including:

* `vanilla/mux`: HTTP routing by method for Go projects

Note: this module is minimal and still in early stages of development. The API should stabilize over time at which point we will begin using semantic versions for the modules (consider this v0.0.0).

---

## vanilla/mux
Package `vanilla/mux` implements a basic HTTP request multiplexer here to help declare routes based on HTTP method. The main features of `mux` are:

* Implements `http.Handler` and is compatible with standard library `http.ListenAndServe(...)`
* Requests can be routed by specific HTTP methods of GET, PUT, POST, DELETE, or OPTIONS using `HttpHandler.GET(...)` etc.
* Requests can be matched to catchall HTTP methods using `HttpHandler.HandleFunc()`
* Similar to the default `http.ServeMux` this handler will match exact routes first and then locate possible routes based on the next closest match. (Longest to shortest)

### Install mux

Assuming a previously configured and working Go toolchain:

```sh
go get -u github.com/icodealot/vanilla/mux
```

### Examples

There are several `Example...` go tests in the `examples` folder that demonstrate vanilla/mux. Clone this repository and then run the examples using `go test` as in:

```
go test ./examples -run ^ExampleHelloWorld$ -v
```

Here is a basic "Hello, World" handler:

```go
package main

import (
	"io"
	"log"
	"net/http"

	"github.com/icodealot/vanilla/mux"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, World!")
}

func main() {
	r := mux.NewHandler()
	r.HandleFunc("/", sayHello)
	log.Fatal(http.ListenAndServe(":8080", r))
}
```

If you follow along with the [Writing Applications Tutorial](https://go.dev/doc/articles/wiki/) you can integrate `vanilla/mux` as follows:

```go
package main

import (
	// code abbreviated for clarity...
    // ...

	"github.com/icodealot/vanilla/mux"
)

func main() {
	r := mux.NewHandler()
	r.GET("/view/", makeHandler(viewHandler))
	r.GET("/edit/", makeHandler(editHandler))
	r.POST("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8080", r))
}
// ...
```
