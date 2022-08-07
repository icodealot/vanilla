// Run this example with:
// go test ./examples -run ^ExampleHelloWorld$ -v
// ... then open a web browser to http://localhost:8080/hello

package examples

import (
	"io"
	"net/http"

	"github.com/icodealot/vanilla/mux"
)

func ExampleHelloWorld() {
	r := mux.NewHandler()

	// Map HTTP GET requests to the sayHello() handler
	r.GET("/hello", sayHello)

	// Map everything else to the index() handler by default
	r.HandleFunc("/", index)

	// Start the HTTP server listening on port 8080
	http.ListenAndServe("127.0.0.1:8080", r)

	// Output:
	// Pressing ctrl+C to kill the server will "FAIL" the test. That is OK.
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, World!")
}

func index(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Goodbye")
}
