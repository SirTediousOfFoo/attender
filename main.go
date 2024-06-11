package main


import (

	"fmt"
	"net/http"
)

func main() {
	// Create a custom file server
	fileServer := http.FileServer(http.Dir("templates/"))
	// Serve up the index page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})
	// Set the correct MIME type for CSS files
	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		fileServer.ServeHTTP(w, r)
	})

	fmt.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
