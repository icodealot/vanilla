// Run this example with:
// go test ./examples -run ^ExampleRestAPI -v
// ...then open a web browser to http://localhost:8080/time
//
// Or, use cURL:
// curl -i http://localhost:8080/time
//
// To POST data on Windows using cURL:
// curl -i --header "content-type: application/json" -d "{\"type\":\"CST\",\"value\":\"Sun, 07 Aug 2022 18:18:03 CST\"}" http://localhost:8080/time

package examples

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/icodealot/vanilla/mux"
)

type TimeData struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

func ExampleRestAPI() {
	r := mux.NewHandler()

	// Map HTTP GET requests to a handler
	r.GET("/time", getTime)

	// Map HTTP POST requests to a handler
	r.POST("/time", postTime)

	// all other requests will return a 404 by default

	// Start the HTTP server listening on port 8080
	http.ListenAndServe("127.0.0.1:8080", r)

	// Output:
	// Pressing ctrl+C to kill the server will "FAIL" the test. That is OK.
}

// getTime handles GET requests and sends the caller a
// type (UTC) and a string representation of the time.
func getTime(w http.ResponseWriter, r *http.Request) {
	t := &TimeData{
		Type:  "UTC",
		Value: time.Now().UTC().Format(http.TimeFormat),
	}
	response, err := json.Marshal(t)
	if err != nil {
		// do something with the error
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("content-type", "application/json")
		w.Write(response)
	}
}

// postTime handles POST requests and sends the caller a
// 200 status code, or error message.
func postTime(w http.ResponseWriter, r *http.Request) {
	var timeData TimeData
	defer r.Body.Close()
	t, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(t, &timeData); err != nil {
		http.Error(w, "unable to read time provided", http.StatusBadRequest)
		return
	}
	log.Printf("Time received: %+v\n", timeData)
	w.WriteHeader(http.StatusOK)
}
