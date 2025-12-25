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
		r.ParseForm()
		name := r.FormValue("name")
		email := r.FormValue("email")
		message := r.FormValue("message")
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

	http.HandleFunc("/form", handleForm)

	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}



