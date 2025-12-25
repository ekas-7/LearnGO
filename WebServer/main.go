package main

import (
	"fmt"
	"log"
	"net/http"
)

func handleHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Handle form submission
		if err := r.ParseForm(); err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
			log.Printf("error parsing form: %v", err)
			return
		}

		name := r.FormValue("name")
		email := r.FormValue("email")
		message := r.FormValue("message")

		// Server-side logging
		clientIP := r.RemoteAddr
		userAgent := r.Header.Get("User-Agent")
		log.Printf("Form submission: ip=%s user_agent=%q name=%q email=%q message=%q", clientIP, userAgent, name, email, message)

		fmt.Fprintf(w, "Received message from %s <%s>: %s\n", name, email, message)
	} else {
		// Serve the form HTML
		http.ServeFile(w, r, "./static/form.html")
	}
}

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)

	http.HandleFunc("/hello", handleHello)

	http.HandleFunc("/submit", handleForm)

	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
