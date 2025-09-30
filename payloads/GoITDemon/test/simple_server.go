package main

import (
	"fmt"
	"log"
	"net/http"
	"io"
)

func main() {
	fmt.Println("[+] Starting simple test server for GoITDemon on :40056")
	
	// Handle all routes
	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/api/login", handleRequest)
	http.HandleFunc("/api/session", handleRequest)
	http.HandleFunc("/api/check", handleRequest)
	
	fmt.Println("[+] Server listening on http://127.0.0.1:40056")
	fmt.Println("[+] Waiting for GoITDemon connections...")
	
	log.Fatal(http.ListenAndServe(":40056", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	
	fmt.Printf("[*] %s %s from %s\n", r.Method, r.URL.Path, r.RemoteAddr)
	fmt.Printf("    Headers: %v\n", r.Header)
	fmt.Printf("    Body length: %d bytes\n", len(body))
	
	if len(body) > 0 {
		fmt.Printf("    Body preview: %x...\n", body[:min(len(body), 32)])
	}
	
	// Send a simple response
	w.Header().Set("Content-Type", "application/octet-stream")
	
	// For initial checkin, send back a simple acknowledgment
	if len(body) > 0 {
		// This is likely metadata or a task request
		response := []byte("OK")
		w.Write(response)
		fmt.Printf("    Response: %d bytes sent\n", len(response))
	} else {
		// Empty response for basic requests
		w.WriteHeader(http.StatusOK)
	}
	
	fmt.Println()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}