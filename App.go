package main

import (
	"gpgHelper/api"
	"net/http"
)

func main() {
	http.HandleFunc("/api/Encode", api.EncodeHandler)

	http.ListenAndServe(":8080", nil)
}
