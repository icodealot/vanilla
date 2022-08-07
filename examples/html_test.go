// Run this example with:
// go test ./examples -run ^ExampleUsingHtmlTemplate -v
// ... then open a web browser to http://localhost:8080/time

package examples

import (
	"html/template"
	"net/http"
	"time"

	"github.com/icodealot/vanilla/mux"
)

type Page struct {
	Title string
	Body  string
}

func ExampleUsingHtmlTemplate() {
	r := mux.NewHandler()

	// Map HTTP GET requests to a handler
	r.GET("/time", tellTime)

	// Start the HTTP server listening on port 8080
	http.ListenAndServe("127.0.0.1:8080", r)

	// Output:
	// Pressing ctrl+C to kill the server will "FAIL" the test. That is OK.
}

func tellTime(w http.ResponseWriter, r *http.Request) {
	page := &Page{
		Title: "What time is it?",
		Body:  "It's UTC time: " + time.Now().UTC().Format(http.TimeFormat),
	}
	t, _ := template.ParseFiles("example.html")
	t.Execute(w, page)
}
