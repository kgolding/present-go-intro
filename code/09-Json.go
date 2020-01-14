package main

import (
	"encoding/json"
	"net/http"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age`
	note string `json:"note"`
}

func main() {
	http.HandleFunc("/", JsonServer)
	http.ListenAndServe(":8080", nil)
}

func JsonServer(w http.ResponseWriter, r *http.Request) {
	a := Person{"Arthur Dent", 42, "HHGTG"}
	buf, err := json.Marshal(a)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
}
