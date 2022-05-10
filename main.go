package main

import (
	"net/http"
	"netHttpTest/api"
)

func main() {
	srv := api.NewServer()
	http.ListenAndServe(":9090", srv)
}
