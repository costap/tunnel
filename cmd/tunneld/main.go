package main

import (
	"fmt"
	"github.com/costap/tunnel/internal/app/tunneld"
	"log"
	"net/http"
)

var s *tunneld.Server

func main() {

	c := tunneld.ConfigInit()

	fmt.Printf("Starting with config %v", c)

	s = tunneld.NewServer(c)

	go s.Run()

	http.HandleFunc("/", index)
	http.HandleFunc("/health", health)
	http.HandleFunc( "/stop", stop)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func index(w http.ResponseWriter, r *http.Request){
	if r.Method != "GET" {
		http.Error(w, "Method not supported.", http.StatusMethodNotAllowed)
	}
	fmt.Fprint(w, "OK")
}

func health(w http.ResponseWriter, r *http.Request){
	if r.Method != "GET" {
		http.Error(w, "Method not supported.", http.StatusMethodNotAllowed)
	}
	if s.IsStarted() {
		fmt.Fprint(w, "{ \"started\": true }")
	} else {
		http.Error(w, "{ \"started\": false }", http.StatusServiceUnavailable)
	}
}

func stop(w http.ResponseWriter, r *http.Request){
	if r.Method != "POST" {
		http.Error(w, "Method not supported.", http.StatusMethodNotAllowed)
	}
	s.Stop()
	fmt.Fprint(w, "STOPPING")
}