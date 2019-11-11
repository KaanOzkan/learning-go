package main

import (
	"encoding/json"
	"github.com/domainr/whois"
	"log"
	"net/http"
)

type Response struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}

func whoisQuery(data string) (string, error) {
	response, err := whois.Fetch(data)
	if err != nil {
		return "", err
	}
	return string(response.Body), nil
}

func jsonResponse(w http.ResponseWriter, x interface{}) {
	bytes, err := json.Marshal(x) // generate json
	if err != nil {
		panic(err)
	}

	// send response to client
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func handleRequests(fileServer http.Handler) {
	// "/"
	http.HandleFunc(
		"/",
		func(w http.ResponseWriter, r *http.Request) {
			fileServer.ServeHTTP(w, r)
		},
	)

	// "/whois"
	http.HandleFunc("/whois", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		data := r.PostFormValue("data")
		result, err := whoisQuery(data)

		if err != nil {
			jsonResponse(w, Response{Error: err.Error()})
			return
		}
		jsonResponse(w, Response{Result: result})
	})

}

func main() {
	fmt.Println("HERE")
	// Serve static files
	fileServer := http.FileServer(http.Dir("static/"))

	// Handle requests
	handleRequests(fileServer)

	// Start server
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
