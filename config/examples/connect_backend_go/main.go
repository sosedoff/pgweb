package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type BackendRequest struct {
	Resource string            `json:"resource"`
	Token    string            `json:"token"`
	Headers  map[string]string `json:"headers"`
}

type BackendResponse struct {
	DatabaseURL string `json:"database_url"`
}

func main() {
	resources := map[string]string{
		"id1": "postgres://localhost:5432/db1?sslmode=disable",
		"id2": "postgres://localhost:5432/db2?sslmode=disable",
		"id3": "postgres://localhost:5432/db3?sslmode=disable",
	}

	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		backendReq := BackendRequest{}

		if err := json.NewDecoder(req.Body).Decode(&backendReq); err != nil {
			rw.WriteHeader(400)
			fmt.Fprintf(rw, "error while parsing request: %v", err)
			return
		}

		res, ok := resources[backendReq.Resource]
		if !ok {
			rw.WriteHeader(404)
			return
		}

		resp := BackendResponse{
			DatabaseURL: res,
		}

		json.NewEncoder(rw).Encode(resp)
	})

	if err := http.ListenAndServe(":4567", nil); err != nil {
		log.Fatal(err)
	}
}
