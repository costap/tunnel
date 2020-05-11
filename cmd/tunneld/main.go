package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/costap/tunnel/internal/app/tunneld"
)

var s *tunneld.Server

func main() {

	c := tunneld.ConfigInit()

	fmt.Printf("Starting with config %v", c)

	s = tunneld.NewServer(c)

	go s.Run()

	http.HandleFunc("/", index)
	http.HandleFunc("/health", health)
	http.HandleFunc("/stop", stop)
	http.HandleFunc("/restart", restart)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", c.AdminPort), nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not supported.", http.StatusMethodNotAllowed)
	}
	fmt.Fprint(w, "OK")
}

func health(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not supported.", http.StatusMethodNotAllowed)
	}
	if s.IsStarted() {
		fmt.Fprint(w, fmt.Sprintf("{ \"started\": true, \"connected\": %v }", s.IsConnected()))
	} else {
		http.Error(w, "{ \"started\": false }", http.StatusServiceUnavailable)
	}
}

func stop(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not supported.", http.StatusMethodNotAllowed)
	}
	s.Stop()
	os.Exit(0)
}

func restart(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not supported.", http.StatusMethodNotAllowed)
	}
	s.Stop()
	time.Sleep(100 * time.Millisecond)
	go s.Run()
}
